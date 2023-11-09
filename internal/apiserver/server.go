// Copyright 2020 Talhuang <talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/skeleton1231/gotal/internal/apiserver/config"
	"github.com/skeleton1231/gotal/internal/pkg/options"
	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/skeleton1231/gotal/pkg/cache"
	"github.com/skeleton1231/gotal/pkg/shutdown"
	posix "github.com/skeleton1231/gotal/pkg/shutdown/managers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type apiServer struct {
	gs            *shutdown.GracefulShutdown
	redisOptions  *options.RedisOptions
	httpAPIServer *server.APIServer
	gRPCAPIServer *grpcAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

// ExtraConfig defines extra configuration for the apiserver.
type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   options.GeneratableKeyCert
	mysqlOptions *options.MySQLOptions
}

func NewAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posix.NewPosixSignalManager())

	// Assgin apiServer config to APIServer
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}
	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:            gs,
		redisOptions:  cfg.RedisOptions,
		httpAPIServer: genericServer,
		gRPCAPIServer: extraServer,
	}

	return server, nil
}

func (s *apiServer) PrepareRun() preparedAPIServer {

	initRouter(s.httpAPIServer.Engine)

	s.initRedisStore()

	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {

		s.gRPCAPIServer.Close()
		s.httpAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	go s.gRPCAPIServer.Run()

	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		logrus.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.httpAPIServer.Run()
}

type completedExtraConfig struct {
	*ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

// New create a grpcAPIServer instance.
func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		logrus.Fatalf("Failed to generate credentials %s", err.Error())
	}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *server.Config, lastErr error) {
	genericConfig = server.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

// nolint: unparam
func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
		ServerCert:   cfg.SecureServing.ServerCert,
		mysqlOptions: cfg.MySQLOptions,
	}, nil
}

func (s *apiServer) initRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())

	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		cancel()

		return nil
	}))

	config := &cache.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}

	// try to connect to redis
	go cache.ConnectToRedisV2(ctx, config)
}
