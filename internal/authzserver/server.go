package authzserver

import (
	"context"
	"log"

	"github.com/skeleton1231/gotal/internal/authzserver/config"
	genericOptions "github.com/skeleton1231/gotal/internal/pkg/options"
	genericApiServer "github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/skeleton1231/gotal/pkg/cache"
	"github.com/skeleton1231/gotal/pkg/shutdown"
	posixsignal "github.com/skeleton1231/gotal/pkg/shutdown/managers"
)

// authzServer struct holds all the necessary fields for the authorization server.
type authzServer struct {
	gs               *shutdown.GracefulShutdown   // Graceful shutdown manager
	rpcServer        string                       // Address of the RPC server
	clientCA         string                       // Client Certificate Authority
	redisOptions     *genericOptions.RedisOptions // Configuration options for Redis
	genericApiServer *genericApiServer.APIServer  // Generic API server
	redisCancelFunc  context.CancelFunc           // Function to cancel Redis context
}

// PrepareRun initializes the server and returns a preparedAuthzServer ready to be run.
func (s *authzServer) PrepareRun() preparedAuthzServer {
	_ = s.initializes()
	// Router Initialization (not shown here)
	return preparedAuthzServer{s}
}

// buildCacheConfig constructs the cache configuration from the server's Redis options.
func (s *authzServer) buildCacheConfig() *cache.Config {
	// Create and return a new cache configuration using the server's Redis options.
	return &cache.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}
}

// initializes sets up necessary components for the authorization server.
func (s *authzServer) initializes() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.redisCancelFunc = cancel

	// Start connecting to Redis with the configuration provided.
	go cache.ConnectToRedisV2(ctx, s.buildCacheConfig())

	return nil
}

// preparedAuthzServer struct represents an authorization server that is ready to run.
type preparedAuthzServer struct {
	*authzServer
}

// Run starts the authorization server and its components.
func (s preparedAuthzServer) Run() error {
	stopCh := make(chan struct{})

	// Start the graceful shutdown manager.
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// Start the generic API server in a separate goroutine.
	go s.genericApiServer.Run()

	// Register a shutdown callback for clean up.
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		// Close API server and cancel Redis context on shutdown.
		s.genericApiServer.Close()
		s.redisCancelFunc()
		return nil
	}))

	// Block until an interrupt signal is received.
	<-stopCh
	return nil
}

// buildGenericConfig constructs the configuration for the generic API server.
func buildGenericConfig(cfg *config.Config) (genericConfig *genericApiServer.Config, err error) {
	// Initialize a new generic API server configuration.
	genericConfig = genericApiServer.NewConfig()

	// Apply various configurations to the generic server config.
	// Any error encountered is returned immediately.
	if err = cfg.GenericServerOptions.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.FeatureOptions.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.SecureServing.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}
	return
}

// createAuthzServer sets up and returns a new authorization server.
func createAuthzServer(cfg *config.Config) (*authzServer, error) {
	// Initialize a new graceful shutdown manager and add a POSIX signal manager.
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// Build the generic API server configuration.
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Complete the configuration and create a new generic API server.
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	// Construct the authorization server with the necessary components.
	server := &authzServer{
		// Assign the relevant fields from the configuration and initialized components.
		gs:               gs,
		redisOptions:     cfg.RedisOptions,
		rpcServer:        cfg.RPCServer,
		clientCA:         cfg.ClientCA,
		genericApiServer: genericServer,
	}

	return server, nil
}
