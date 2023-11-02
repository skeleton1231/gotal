// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package storage defines redis storage.

package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"errors"

	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config defines options for redis cluster.
type Config struct {
	Host                  string
	Port                  int
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	EnableCluster         bool
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

// ErrRedisIsDown is returned when we can't communicate with redis.
var ErrRedisIsDown = errors.New("storage: Redis is either down or ws not configured")

var (
	singlePool      atomic.Value
	singleCachePool atomic.Value
	redisUp         atomic.Value
)

var disableRedis atomic.Value

// DisableRedis very handy when testsing it allows to dynamically enable/disable talking with redisW.
func DisableRedis(ok bool) {
	if ok {
		redisUp.Store(false)
		disableRedis.Store(true)

		return
	}
	redisUp.Store(true)
	disableRedis.Store(false)
}

func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}

	return true
}

// Connected returns true if we are connected to redis.
func Connected() bool {
	if v := redisUp.Load(); v != nil {
		return v.(bool)
	}

	return false
}

func singleton(cache bool) redis.UniversalClient {
	if cache {
		v := singleCachePool.Load()
		if v != nil {
			return v.(redis.UniversalClient)
		}

		return nil
	}
	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}

	return nil
}

func singletonV2(cache bool) redis.UniversalClient {
	if cache {
		v := singleCachePool.Load()
		if v != nil {
			client, ok := v.(redis.UniversalClient)
			if !ok {
				logrus.Error("Stored value in singleCachePool is not of type redis.UniversalClient")
				return nil
			}
			return client
		}

		return nil
	}

	v := singlePool.Load()
	if v != nil {
		client, ok := v.(redis.UniversalClient)
		if !ok {
			logrus.Error("Stored value in singlePool is not of type redis.UniversalClient")
			return nil
		}
		return client
	}

	return nil
}

type RedisClusterV2 struct {
	KeyPrefix string
	HashKeys  bool
	IsCache   bool
}

func ConnectToRedisV2(ctx context.Context, config *Config) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	c := []RedisClusterV2{
		{},
		{IsCache: true},
	}

	for {
		if shouldConnect() {
			ok := establishConnection(ctx, c, config)
			redisUp.Store(ok)
		}

		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			continue
		}
	}
}

func establishConnection(ctx context.Context, clusters []RedisClusterV2, config *Config) bool {
	for _, v := range clusters {
		if !connectSingletonV2(ctx, v.IsCache, config) {
			return false
		}
		if !clusterConnectionIsOpenV2(ctx, v) {
			return false
		}
	}
	return true
}

func connectSingletonV2(ctx context.Context, isCache bool, config *Config) bool {
	// NOTE: Assuming the connect logic remains the same, but if it differs, adjust accordingly.
	if singletonV2(isCache) == nil {
		logrus.Debug("Connecting to redis cluster")
		if isCache {
			singleCachePool.Store(NewRedisClusterPoolV2(isCache, config))
			return true
		}
		singlePool.Store(NewRedisClusterPoolV2(isCache, config))
		return true
	}
	return true
}

func clusterConnectionIsOpenV2(ctx context.Context, cluster RedisClusterV2) bool {
	c := singletonV2(cluster.IsCache)
	if c == nil {
		logrus.Warn("Redis client is nil")
		return false
	}

	// Generating UUID
	testKey := "redis-test-" + uuid.NewV4().String()

	if err := c.Set(ctx, testKey, "test", time.Second).Err(); err != nil {
		logrus.Warnf("Error trying to set test key: %s", err.Error())
		return false
	}
	if _, err := c.Get(ctx, testKey).Result(); err != nil {
		logrus.Warnf("Error trying to get test key: %s", err.Error())
		return false
	}

	return true
}

func getRedisAddrs(config *Config) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

// RedisOpts is the overridden type of redis.UniversalOptions. simple() and cluster() functions are not public in redis
// library.
// Therefore, they are redefined here to use in the creation of a new redis cluster logic.
// We don't want to use redis.NewUniversalClient() logic.
type RedisOptsV2 redis.UniversalOptions

func (opts *RedisOptsV2) simple() *redis.Options {
	return &redis.Options{
		Addr:         opts.Addrs[0],
		Password:     opts.Password,
		DB:           opts.DB,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolSize:     opts.PoolSize,
		TLSConfig:    opts.TLSConfig,
	}
}

func (opts *RedisOptsV2) cluster() *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:        opts.Addrs,
		Password:     opts.Password,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolSize:     opts.PoolSize,
		TLSConfig:    opts.TLSConfig,
	}
}

func (opts *RedisOptsV2) failover() *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:   opts.MasterName,
		Password:     opts.Password,
		DB:           opts.DB,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolSize:     opts.PoolSize,
		TLSConfig:    opts.TLSConfig,
	}
}

// NewRedisClusterPool create a redis cluster pool.
func NewRedisClusterPoolV2(isCache bool, config *Config) redis.UniversalClient {
	// redisSingletonMu is locked and we know the singleton is nil
	logrus.Debug("Creating new Redis connection pool")

	// poolSize applies per cluster node and not for the whole cluster.
	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	timeout := 5 * time.Second

	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config

	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := &RedisOptsV2{
		Addrs:        getRedisAddrs(config),
		MasterName:   config.MasterName,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		//IdleTimeout:  240 * timeout,
		PoolSize:  poolSize,
		TLSConfig: tlsConfig,
	}

	if opts.MasterName != "" {
		logrus.Info("--> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(opts.failover())
	} else if config.EnableCluster {
		logrus.Info("--> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(opts.cluster())
	} else {
		logrus.Info("--> [REDIS] Creating single-node client")
		client = redis.NewClient(opts.simple())
	}

	return client
}

// Connect will establish a connection this is always true because we are dynamically using redis.
func (r *RedisClusterV2) Connect() bool {
	return true
}

func (r *RedisClusterV2) singleton() redis.UniversalClient {
	return singletonV2(r.IsCache)
}

func (r *RedisClusterV2) hashKey(in string) string {
	if !r.HashKeys {
		// Not hashing? Return the raw key
		return in
	}

	return HashStr(in)
}

func (r *RedisClusterV2) fixKey(keyName string) string {
	return r.KeyPrefix + r.hashKey(keyName)
}

func (r *RedisClusterV2) cleanKey(keyName string) string {
	return strings.Replace(keyName, r.KeyPrefix, "", 1)
}

func (r *RedisClusterV2) up() error {
	if !Connected() {
		return ErrRedisIsDown
	}

	return nil
}

// GetKey will retrieve a key from the database.
func (r *RedisClusterV2) GetKey(ctx context.Context, keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}

	cluster := r.singleton()

	value, err := cluster.Get(ctx, r.fixKey(keyName)).Result()
	if err != nil {
		logrus.Debugf("Error trying to get value: %s", err.Error())

		return "", ErrKeyNotFound
	}

	return value, nil
}

func (r *RedisClusterV2) GetMultiKey(ctx context.Context, keys []string) ([]string, error) {
	if err := r.up(); err != nil {
		return nil, err
	}
	cluster := r.singleton()

	// Directly create the fixed keyNames slice
	keyNames := make([]string, len(keys))
	for index, key := range keys {
		keyNames[index] = r.fixKey(key)
	}

	result := make([]string, 0)

	switch v := cluster.(type) {
	case *redis.ClusterClient:
		// ... (rest of the logic remains the same)
		getCmds := make([]*redis.StringCmd, len(keyNames))
		pipe := v.Pipeline()
		for i, key := range keyNames {
			getCmds[i] = pipe.Get(ctx, key)
		}
		_, err := pipe.Exec(ctx)
		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.Errorf("Error trying to get value: %s", err.Error())
			return nil, err
		}
		for _, cmd := range getCmds {
			result = append(result, cmd.Val())
		}

	case *redis.Client:
		values, err := v.MGet(ctx, keyNames...).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.Errorf("Error trying to get value: %s", err.Error())
			return nil, err
		}
		for _, val := range values {
			if val == nil {
				result = append(result, "")
			} else {
				result = append(result, fmt.Sprint(val))
			}
		}
	}

	notEmpty := false
	for _, val := range result {
		if val != "" {
			notEmpty = true
			break
		}
	}

	if notEmpty {
		return result, nil
	}
	return nil, ErrKeyNotFound
}

// GetKeyTTL return ttl of the given key.
func (r *RedisClusterV2) GetKeyTTL(ctx context.Context, keyName string) (ttl int64, err error) {
	if err = r.up(); err != nil {
		return 0, err
	}
	duration, err := r.singleton().TTL(ctx, r.fixKey(keyName)).Result()

	return int64(duration.Seconds()), err
}

// GetRawKey return the value of the given key.
func (r *RedisClusterV2) GetRawKey(ctx context.Context, keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}
	value, err := r.singleton().Get(ctx, keyName).Result()
	if err != nil {
		logrus.Debugf("Error trying to get value: %s", err.Error())

		return "", ErrKeyNotFound
	}

	return value, nil
}

// GetExp return the expiry of the given key.
func (r *RedisClusterV2) GetExp(ctx context.Context, keyName string) (int64, error) {
	logrus.Debugf("Getting exp for key: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		return 0, err
	}

	value, err := r.singleton().TTL(ctx, r.fixKey(keyName)).Result()
	if err != nil {
		logrus.Errorf("Error trying to get TTL: ", err.Error())

		return 0, ErrKeyNotFound
	}

	return int64(value.Seconds()), nil
}

// SetExp set expiry of the given key.
func (r *RedisClusterV2) SetExp(ctx context.Context, keyName string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Expire(ctx, r.fixKey(keyName), timeout).Err()
	if err != nil {
		logrus.Errorf("Could not EXPIRE key: %s", err.Error())
	}

	return err
}

// SetKey will create (or update) a key value in the store.
func (r *RedisClusterV2) SetKey(ctx context.Context, keyName, session string, timeout time.Duration) error {
	logrus.Debugf("[STORE] SET Raw key is: %s", keyName)
	logrus.Debugf("[STORE] Setting key: %s", r.fixKey(keyName))

	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Set(ctx, r.fixKey(keyName), session, timeout).Err()
	if err != nil {
		logrus.Errorf("Error trying to set value: %s", err.Error())

		return err
	}

	return nil
}

// SetRawKey set the value of the given key.
func (r *RedisClusterV2) SetRawKey(ctx context.Context, keyName, session string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Set(ctx, keyName, session, timeout).Err()
	if err != nil {
		logrus.Errorf("Error trying to set value: %s", err.Error())

		return err
	}

	return nil
}

// Decrement will decrement a key in redis.
func (r *RedisClusterV2) Decrement(ctx context.Context, keyName string) {
	keyName = r.fixKey(keyName)
	logrus.Debugf("Decrementing key: %s", keyName)
	if err := r.up(); err != nil {
		return
	}
	err := r.singleton().Decr(ctx, keyName).Err()
	if err != nil {
		logrus.Errorf("Error trying to decrement value: %s", err.Error())
	}
}

// IncrememntWithExpire will increment a key in redis.
func (r *RedisClusterV2) IncrememntWithExpire(ctx context.Context, keyName string, expire int64) int64 {
	logrus.Debugf("Incrementing raw key: %s", keyName)
	if err := r.up(); err != nil {
		return 0
	}
	// This function uses a raw key, so we shouldn't call fixKey
	fixedKey := keyName
	val, err := r.singleton().Incr(ctx, fixedKey).Result()

	if err != nil {
		logrus.Errorf("Error trying to increment value: %s", err.Error())
	} else {
		logrus.Debugf("Incremented key: %s, val is: %d", fixedKey, val)
	}

	if val == 1 && expire > 0 {
		logrus.Debug("--> Setting Expire")
		r.singleton().Expire(ctx, fixedKey, time.Duration(expire)*time.Second)
	}

	return val
}

// GetKeys will return all keys according to the filter (filter is a prefix - e.g. tyk.keys.*).
func (r *RedisClusterV2) GetKeys(ctx context.Context, filter string) []string {
	if err := r.up(); err != nil {
		return nil
	}
	client := r.singleton()

	filterHash := ""
	if filter != "" {
		filterHash = r.hashKey(filter)
	}
	searchStr := r.KeyPrefix + filterHash + "*"
	logrus.Debugf("[STORE] Getting list by: %s", searchStr)

	fnFetchKeys := func(client *redis.Client) ([]string, error) {
		values := make([]string, 0)

		iter := client.Scan(ctx, 0, searchStr, 0).Iterator()
		for iter.Next(ctx) {
			values = append(values, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return nil, err
		}

		return values, nil
	}

	var err error
	var values []string
	sessions := make([]string, 0)

	switch v := client.(type) {
	case *redis.ClusterClient:
		ch := make(chan []string)

		go func() {
			err = v.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
				values, err = fnFetchKeys(client)
				if err != nil {
					return err
				}

				ch <- values

				return nil
			})
			close(ch)
		}()

		for res := range ch {
			sessions = append(sessions, res...)
		}
	case *redis.Client:
		sessions, err = fnFetchKeys(v)
	}

	if err != nil {
		logrus.Errorf("Error while fetching keys: %s", err)

		return nil
	}

	for i, v := range sessions {
		sessions[i] = r.cleanKey(v)
	}

	return sessions
}

// GetKeysAndValuesWithFilter will return all keys and their values with a filter.
func (r *RedisClusterV2) GetKeysAndValuesWithFilter(ctx context.Context, filter string) map[string]string {
	if err := r.up(); err != nil {
		return nil
	}
	keys := r.GetKeys(ctx, filter)
	if keys == nil {
		logrus.Error("Error trying to get filtered client keys")

		return nil
	}

	if len(keys) == 0 {
		return nil
	}

	for i, v := range keys {
		keys[i] = r.KeyPrefix + v
	}

	client := r.singleton()
	values := make([]string, 0)

	switch v := client.(type) {
	case *redis.ClusterClient:
		{
			getCmds := make([]*redis.StringCmd, 0)
			pipe := v.Pipeline()
			for _, key := range keys {
				getCmds = append(getCmds, pipe.Get(ctx, key))
			}
			_, err := pipe.Exec(ctx)
			if err != nil && !errors.Is(err, redis.Nil) {
				logrus.Errorf("Error trying to get client keys: %s", err.Error())

				return nil
			}

			for _, cmd := range getCmds {
				values = append(values, cmd.Val())
			}
		}
	case *redis.Client:
		{
			result, err := v.MGet(ctx, keys...).Result()
			if err != nil {
				logrus.Errorf("Error trying to get client keys: %s", err.Error())

				return nil
			}

			for _, val := range result {
				strVal := fmt.Sprint(val)
				if strVal == "<nil>" {
					strVal = ""
				}
				values = append(values, strVal)
			}
		}
	}

	m := make(map[string]string)
	for i, v := range keys {
		m[r.cleanKey(v)] = values[i]
	}

	return m
}

// GetKeysAndValues will return all keys and their values - not to be used lightly.
func (r *RedisClusterV2) GetKeysAndValues(ctx context.Context) map[string]string {
	return r.GetKeysAndValuesWithFilter(ctx, "")
}

// DeleteKey will remove a key from the database.
func (r *RedisClusterV2) DeleteKey(ctx context.Context, keyName string) bool {
	if err := r.up(); err != nil {
		// logrus.Debug(err)
		return false
	}
	logrus.Debugf("DEL Key was: %s", keyName)
	logrus.Debugf("DEL Key became: %s", r.fixKey(keyName))
	n, err := r.singleton().Del(ctx, r.fixKey(keyName)).Result()
	if err != nil {
		logrus.Errorf("Error trying to delete key: %s", err.Error())
	}

	return n > 0
}

// DeleteAllKeys will remove all keys from the database.
func (r *RedisClusterV2) DeleteAllKeys(ctx context.Context) bool {
	if err := r.up(); err != nil {
		return false
	}
	n, err := r.singleton().FlushAll(ctx).Result()
	if err != nil {
		logrus.Errorf("Error trying to delete keys: %s", err.Error())
	}

	if n == "OK" {
		return true
	}

	return false
}

// DeleteRawKey will remove a key from the database without prefixing, assumes user knows what they are doing.
func (r *RedisClusterV2) DeleteRawKey(ctx context.Context, keyName string) bool {
	if err := r.up(); err != nil {
		return false
	}
	n, err := r.singleton().Del(ctx, keyName).Result()
	if err != nil {
		logrus.Errorf("Error trying to delete key: %s", err.Error())
	}

	return n > 0
}

// DeleteScanMatch will remove a group of keys in bulk.
func (r *RedisClusterV2) DeleteScanMatch(ctx context.Context, pattern string) bool {
	if err := r.up(); err != nil {
		return false
	}
	client := r.singleton()
	logrus.Debugf("Deleting: %s", pattern)

	fnScan := func(client *redis.Client) ([]string, error) {
		values := make([]string, 0)

		iter := client.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			values = append(values, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return nil, err
		}

		return values, nil
	}

	var err error
	var keys []string
	var values []string

	switch v := client.(type) {
	case *redis.ClusterClient:
		ch := make(chan []string)
		go func() {
			err = v.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
				values, err = fnScan(client)
				if err != nil {
					return err
				}

				ch <- values

				return nil
			})
			close(ch)
		}()

		for vals := range ch {
			keys = append(keys, vals...)
		}
	case *redis.Client:
		keys, err = fnScan(v)
	}

	if err != nil {
		logrus.Errorf("SCAN command field with err: %s", err.Error())

		return false
	}

	if len(keys) > 0 {
		for _, name := range keys {
			logrus.Infof("Deleting: %s", name)
			err := client.Del(ctx, name).Err()
			if err != nil {
				logrus.Errorf("Error trying to delete key: %s - %s", name, err.Error())
			}
		}
		logrus.Infof("Deleted: %d records", len(keys))
	} else {
		logrus.Debug("RedisClusterV2 called DEL - Nothing to delete")
	}

	return true
}

// DeleteKeys will remove a group of keys in bulk.
func (r *RedisClusterV2) DeleteKeys(ctx context.Context, keys []string) bool {
	if err := r.up(); err != nil {
		return false
	}
	if len(keys) > 0 {
		for i, v := range keys {
			keys[i] = r.fixKey(v)
		}

		logrus.Debugf("Deleting: %v", keys)
		client := r.singleton()
		switch v := client.(type) {
		case *redis.ClusterClient:
			{
				pipe := v.Pipeline()
				for _, k := range keys {
					pipe.Del(ctx, k)
				}

				if _, err := pipe.Exec(ctx); err != nil {
					logrus.Errorf("Error trying to delete keys: %s", err.Error())
				}
			}
		case *redis.Client:
			{
				_, err := v.Del(ctx, keys...).Result()
				if err != nil {
					logrus.Errorf("Error trying to delete keys: %s", err.Error())
				}
			}
		}
	} else {
		logrus.Debug("RedisClusterV2 called DEL - Nothing to delete")
	}

	return true
}

// StartPubSubHandler will listen for a signal and run the callback for
// every subscription and message event.
func (r *RedisClusterV2) StartPubSubHandler(ctx context.Context, channel string, callback func(interface{})) error {
	if err := r.up(); err != nil {
		return err
	}
	client := r.singleton()
	if client == nil {
		return errors.New("redis connection failed")
	}

	pubsub := client.Subscribe(ctx, channel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
		logrus.Errorf("Error while receiving pubsub message: %s", err.Error())

		return err
	}

	for msg := range pubsub.Channel() {
		callback(msg)
	}

	return nil
}

// Publish publish a message to the specify channel.
func (r *RedisClusterV2) Publish(ctx context.Context, channel, message string) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Publish(ctx, channel, message).Err()
	if err != nil {
		logrus.Errorf("Error trying to set value: %s", err.Error())

		return err
	}

	return nil
}

// GetAndDeleteSet get and delete a key.
func (r *RedisClusterV2) GetAndDeleteSet(ctx context.Context, keyName string) []interface{} {
	logrus.Debugf("Getting raw key set: %s", keyName)
	if err := r.up(); err != nil {
		return nil
	}
	logrus.Debugf("keyName is: %s", keyName)
	fixedKey := r.fixKey(keyName)
	logrus.Debugf("Fixed keyname is: %s", fixedKey)

	client := r.singleton()

	var lrange *redis.StringSliceCmd
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(ctx, fixedKey, 0, -1)
		pipe.Del(ctx, fixedKey)

		return nil
	})
	if err != nil {
		logrus.Errorf("Multi command failed: %s", err.Error())

		return nil
	}

	vals := lrange.Val()
	logrus.Debugf("Analytics returned: %d", len(vals))
	if len(vals) == 0 {
		return nil
	}

	logrus.Debugf("Unpacked vals: %d", len(vals))
	result := make([]interface{}, len(vals))
	for i, v := range vals {
		result[i] = v
	}

	return result
}

// AppendToSet append a value to the key set.
func (r *RedisClusterV2) AppendToSet(ctx context.Context, keyName, value string) {
	fixedKey := r.fixKey(keyName)

	logrus.WithField("keyName", keyName).Debug("Pushing to raw key list")
	logrus.WithField("fixedKey", fixedKey).Debug("Appending to fixed key list")

	if err := r.up(); err != nil {
		return
	}
	if err := r.singleton().RPush(ctx, fixedKey, value).Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"keyName":  keyName,
			"fixedKey": fixedKey,
			"error":    err.Error(),
		}).Error("Error trying to append to set keys")
	}
}

// Exists check if keyName exists.
func (r *RedisClusterV2) Exists(ctx context.Context, keyName string) (bool, error) {
	fixedKey := r.fixKey(keyName)

	logrus.WithField("keyName", fixedKey).Debug("Checking if exists")

	exists, err := r.singleton().Exists(ctx, fixedKey).Result()
	if err != nil {
		logrus.WithField("keyName", fixedKey).Errorf("Error trying to check if key exists: %s", err.Error())
		return false, err
	}
	if exists == 1 {
		return true, nil
	}

	return false, nil
}

// RemoveFromList delete a value from a list identified with the keyName.
func (r *RedisClusterV2) RemoveFromList(ctx context.Context, keyName, value string) error {
	fixedKey := r.fixKey(keyName)

	logrus.WithFields(logrus.Fields{
		"keyName":  keyName,
		"fixedKey": fixedKey,
		"value":    value,
	}).Debug("Removing value from list")

	if err := r.singleton().LRem(ctx, fixedKey, 0, value).Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"keyName":  keyName,
			"fixedKey": fixedKey,
			"value":    value,
			"error":    err.Error(),
		}).Error("LREM command failed")

		return err
	}

	return nil
}

// GetListRange gets range of elements of list identified by keyName.
func (r *RedisClusterV2) GetListRange(ctx context.Context, keyName string, from, to int64) ([]string, error) {
	fixedKey := r.fixKey(keyName)

	elements, err := r.singleton().LRange(ctx, fixedKey, from, to).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"keyName":  keyName,
			"fixedKey": fixedKey,
			"from":     from,
			"to":       to,
			"error":    err.Error(),
		}).Error("LRANGE command failed")

		return nil, err
	}

	return elements, nil
}

// AppendToSetPipelined append values to redis pipeline.
func (r *RedisClusterV2) AppendToSetPipelined(ctx context.Context, key string, values [][]byte) {
	if len(values) == 0 {
		return
	}

	fixedKey := r.fixKey(key)
	if err := r.up(); err != nil {
		logrus.Debug(err.Error())

		return
	}
	client := r.singleton()

	pipe := client.Pipeline()
	for _, val := range values {
		pipe.RPush(ctx, fixedKey, val)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		logrus.Errorf("Error trying to append to set keys: %s", err.Error())
	}

	// if we need to set an expiration time
	if storageExpTime := int64(viper.GetDuration("analytics.storage-expiration-time")); storageExpTime != int64(-1) {
		// If there is no expiry on the analytics set, we should set it.
		exp, _ := r.GetExp(ctx, key)
		if exp == -1 {
			_ = r.SetExp(ctx, key, time.Duration(storageExpTime)*time.Second)
		}
	}
}

// GetSet return key set value.
func (r *RedisClusterV2) GetSet(ctx context.Context, keyName string) (map[string]string, error) {
	logrus.Debugf("Getting from key set: %s", keyName)
	logrus.Debugf("Getting from fixed key set: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		return nil, err
	}
	val, err := r.singleton().SMembers(ctx, r.fixKey(keyName)).Result()
	if err != nil {
		logrus.Errorf("Error trying to get key set: %s", err.Error())

		return nil, err
	}

	result := make(map[string]string)
	for i, value := range val {
		result[strconv.Itoa(i)] = value
	}

	return result, nil
}

// AddToSet add value to key set.
func (r *RedisClusterV2) AddToSet(ctx context.Context, keyName, value string) {
	logrus.Debugf("Pushing to raw key set: %s", keyName)
	logrus.Debugf("Pushing to fixed key set: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		return
	}
	err := r.singleton().SAdd(ctx, r.fixKey(keyName), value).Err()
	if err != nil {
		logrus.Errorf("Error trying to append keys: %s", err.Error())
	}
}

// RemoveFromSet remove a value from key set.
func (r *RedisClusterV2) RemoveFromSet(ctx context.Context, keyName, value string) {
	logrus.Debugf("Removing from raw key set: %s", keyName)
	logrus.Debugf("Removing from fixed key set: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		logrus.Debug(err.Error())

		return
	}
	err := r.singleton().SRem(ctx, r.fixKey(keyName), value).Err()
	if err != nil {
		logrus.Errorf("Error trying to remove keys: %s", err.Error())
	}
}

// IsMemberOfSet return whether the given value belong to key set.
func (r *RedisClusterV2) IsMemberOfSet(ctx context.Context, keyName, value string) bool {
	if err := r.up(); err != nil {
		logrus.Debug(err.Error())

		return false
	}
	val, err := r.singleton().SIsMember(ctx, r.fixKey(keyName), value).Result()
	if err != nil {
		logrus.Errorf("Error trying to check set member: %s", err.Error())

		return false
	}

	logrus.Debugf("SISMEMBER %s %s %v %v", keyName, value, val, err)

	return val
}

// SetRollingWindow will append to a sorted set in redis and extract a timed window of values.
func (r *RedisClusterV2) SetRollingWindow(
	ctx context.Context,
	keyName string,
	per int64,
	valueOverride string,
	pipeline bool,
) (int, []interface{}) {
	logrus.Debugf("Incrementing raw key: %s", keyName)
	if err := r.up(); err != nil {
		logrus.Debug(err.Error())

		return 0, nil
	}
	logrus.Debugf("keyName is: %s", keyName)
	now := time.Now()
	logrus.Debugf("Now is: %v", now)
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)
	logrus.Debugf("Then is: %v", onePeriodAgo)

	client := r.singleton()
	var zrange *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		zrange = pipe.ZRange(ctx, keyName, 0, -1)

		element := redis.Z{
			Score: float64(now.UnixNano()),
		}

		if valueOverride != "-1" {
			element.Member = valueOverride
		} else {
			element.Member = strconv.Itoa(int(now.UnixNano()))
		}

		pipe.ZAdd(ctx, keyName, element)
		pipe.Expire(ctx, keyName, time.Duration(per)*time.Second)

		return nil
	}

	var err error
	if pipeline {
		_, err = client.Pipelined(ctx, pipeFn)
	} else {
		_, err = client.TxPipelined(ctx, pipeFn)
	}

	if err != nil {
		logrus.Errorf("Multi command failed: %s", err.Error())

		return 0, nil
	}

	values := zrange.Val()

	// Check actual value
	if values == nil {
		return 0, nil
	}

	intVal := len(values)
	result := make([]interface{}, len(values))

	for i, v := range values {
		result[i] = v
	}

	logrus.Debugf("Returned: %d", intVal)

	return intVal, result
}

// GetRollingWindow return rolling window.
func (r RedisClusterV2) GetRollingWindow(ctx context.Context, keyName string, per int64, pipeline bool) (int, []interface{}) {
	if err := r.up(); err != nil {
		logrus.Debug(err.Error())

		return 0, nil
	}
	now := time.Now()
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)

	client := r.singleton()
	var zrange *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		zrange = pipe.ZRange(ctx, keyName, 0, -1)

		return nil
	}

	var err error
	if pipeline {
		_, err = client.Pipelined(ctx, pipeFn)
	} else {
		_, err = client.TxPipelined(ctx, pipeFn)
	}
	if err != nil {
		logrus.Errorf("Multi command failed: %s", err.Error())

		return 0, nil
	}

	values := zrange.Val()

	// Check actual value
	if values == nil {
		return 0, nil
	}

	intVal := len(values)
	result := make([]interface{}, intVal)
	for i, v := range values {
		result[i] = v
	}

	logrus.Debugf("Returned: %d", intVal)

	return intVal, result
}

// GetKeyPrefix returns storage key prefix.
func (r *RedisClusterV2) GetKeyPrefix() string {
	return r.KeyPrefix
}

// AddToSortedSet adds value with given score to sorted set identified by keyName.
func (r *RedisClusterV2) AddToSortedSet(ctx context.Context, keyName, value string, score float64) {
	fixedKey := r.fixKey(keyName)

	logrus.WithFields(logrus.Fields{
		"keyName":  keyName,
		"fixedKey": fixedKey,
	}).Debug("Pushing raw key to sorted set")

	if err := r.up(); err != nil {
		logrus.Debug(err.Error())
		return
	}

	member := redis.Z{Score: score, Member: value}
	if err := r.singleton().ZAdd(ctx, fixedKey, member).Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"keyName":  keyName,
			"fixedKey": fixedKey,
			"error":    err.Error(),
		}).Error("ZADD command failed")
	}
}

// GetSortedSetRange gets range of elements of sorted set identified by keyName.
func (r *RedisClusterV2) GetSortedSetRange(ctx context.Context, keyName, scoreFrom, scoreTo string) ([]string, []float64, error) {
	fixedKey := r.fixKey(keyName)
	logrus.WithFields(logrus.Fields{
		"keyName":   keyName,
		"fixedKey":  fixedKey,
		"scoreFrom": scoreFrom,
		"scoreTo":   scoreTo,
	}).Debug("Getting sorted set range")

	args := redis.ZRangeBy{Min: scoreFrom, Max: scoreTo}
	values, err := r.singleton().ZRangeByScoreWithScores(ctx, fixedKey, &args).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"keyName":   keyName,
			"fixedKey":  fixedKey,
			"scoreFrom": scoreFrom,
			"scoreTo":   scoreTo,
			"error":     err.Error(),
		}).Error("ZRANGEBYSCORE command failed")
		return nil, nil, err
	}

	if len(values) == 0 {
		return nil, nil, nil
	}

	elements := make([]string, len(values))
	scores := make([]float64, len(values))

	for i, v := range values {
		elements[i] = fmt.Sprint(v.Member)
		scores[i] = v.Score
	}

	return elements, scores, nil
}

// RemoveSortedSetRange removes range of elements from sorted set identified by keyName.
func (r *RedisClusterV2) RemoveSortedSetRange(ctx context.Context, keyName, scoreFrom, scoreTo string) error {
	fixedKey := r.fixKey(keyName)

	logrus.WithFields(logrus.Fields{
		"keyName":   keyName,
		"fixedKey":  fixedKey,
		"scoreFrom": scoreFrom,
		"scoreTo":   scoreTo,
	}).Debug("Removing sorted set range")

	if err := r.singleton().ZRemRangeByScore(ctx, fixedKey, scoreFrom, scoreTo).Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"keyName":   keyName,
			"fixedKey":  fixedKey,
			"scoreFrom": scoreFrom,
			"scoreTo":   scoreTo,
			"error":     err.Error(),
		}).Error("ZREMRANGEBYSCORE command failed")

		return err
	}

	return nil
}
