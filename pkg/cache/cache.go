// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package storage defines redis storage.
package cache

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"

	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/twmb/murmur3"
)

// ErrKeyNotFound is a standard error for when a key is not found in the storage engine.
var ErrKeyNotFound = errors.New("key not found")

// Handler is a standard interface to a cache backend
type Handler interface {
	GetKey(ctx context.Context, key string) (string, error)
	GetMultiKey(ctx context.Context, keys []string) ([]string, error)
	GetRawKey(ctx context.Context, key string) (string, error)
	SetKey(ctx context.Context, key, value string, ttl int64) error
	SetRawKey(ctx context.Context, key, value string, ttl int64) error
	SetExp(ctx context.Context, key string, ttl int64) error
	GetExp(ctx context.Context, key string) (int64, error)
	GetKeys(ctx context.Context, pattern string) []string
	DeleteKey(ctx context.Context, key string) bool
	DeleteAllKeys(ctx context.Context) bool
	DeleteRawKey(ctx context.Context, key string) bool
	Connect() bool
	GetKeysAndValues(ctx context.Context) map[string]string
	GetKeysAndValuesWithFilter(ctx context.Context, pattern string) map[string]string
	DeleteKeys(ctx context.Context, keys []string) bool
	Decrement(ctx context.Context, key string)
	IncrememntWithExpire(ctx context.Context, key string, ttl int64) int64
	SetRollingWindow(ctx context.Context, key string, per int64, val string, pipeline bool) (int, []interface{})
	GetRollingWindow(ctx context.Context, key string, per int64, pipeline bool) (int, []interface{})
	GetSet(ctx context.Context, key string) (map[string]string, error)
	AddToSet(ctx context.Context, key, value string)
	GetAndDeleteSet(ctx context.Context, key string) []interface{}
	RemoveFromSet(ctx context.Context, key, value string)
	DeleteScanMatch(ctx context.Context, pattern string) bool
	GetKeyPrefix() string
	AddToSortedSet(ctx context.Context, key, value string, score float64)
	GetSortedSetRange(ctx context.Context, key, min, max string) ([]string, []float64, error)
	RemoveSortedSetRange(ctx context.Context, key, min, max string) error
	GetListRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	RemoveFromList(ctx context.Context, key, value string) error
	AppendToSet(ctx context.Context, key, value string)
	Exists(ctx context.Context, key string) (bool, error)
}

const defaultHashAlgorithm = "murmur64"

// GenerateToken generate token, if hashing algorithm is empty, use legacy key generation.
func GenerateToken(orgID, keyID, hashAlgorithm string) (string, error) {
	if keyID == "" {
		keyID = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	}

	if hashAlgorithm != "" {
		if _, err := hashFunction(hashAlgorithm); err != nil {
			hashAlgorithm = defaultHashAlgorithm
		}

		jsonToken := fmt.Sprintf(`{"org":"%s","id":"%s","h":"%s"}`, orgID, keyID, hashAlgorithm)
		encodedToken := base64.StdEncoding.EncodeToString([]byte(jsonToken))

		return encodedToken, nil
	}

	// Legacy keys
	return orgID + keyID, nil
}

// B64JSONPrefix stand for `{"` in base64.
const B64JSONPrefix = "ey"

// TokenHashAlgo ...
func TokenHashAlgo(token string) string {
	// Legacy tokens not b64 and not JSON records
	if strings.HasPrefix(token, B64JSONPrefix) {
		if jsonToken, err := base64.StdEncoding.DecodeString(token); err == nil {
			hashAlgo, _ := jsonparser.GetString(jsonToken, "h")

			return hashAlgo
		}
	}

	return ""
}

// TokenOrg ...
func TokenOrg(token string) string {
	if strings.HasPrefix(token, B64JSONPrefix) {
		if jsonToken, err := base64.StdEncoding.DecodeString(token); err == nil {
			// Checking error in case if it is a legacy tooken which just by accided has the same b64JSON prefix
			if org, err := jsonparser.GetString(jsonToken, "org"); err == nil {
				return org
			}
		}
	}

	// 24 is mongo bson id length
	if len(token) > 24 {
		return token[:24]
	}

	return ""
}

// Defines algorithm constant.
var (
	HashSha256    = "sha256"
	HashMurmur32  = "murmur32"
	HashMurmur64  = "murmur64"
	HashMurmur128 = "murmur128"
)

func hashFunction(algorithm string) (hash.Hash, error) {
	switch algorithm {
	case HashSha256:
		return sha256.New(), nil
	case HashMurmur64:
		return murmur3.New64(), nil
	case HashMurmur128:
		return murmur3.New128(), nil
	case "", HashMurmur32:
		return murmur3.New32(), nil
	default:
		return murmur3.New32(), fmt.Errorf("unknown key hash function: %s. Falling back to murmur32", algorithm)
	}
}

// HashStr return hash the give string and return.
func HashStr(in string) string {
	h, _ := hashFunction(TokenHashAlgo(in))
	_, _ = h.Write([]byte(in))

	return hex.EncodeToString(h.Sum(nil))
}

// HashKey return hash the give string and return.
func HashKey(in string) string {
	return HashStr(in)
}
