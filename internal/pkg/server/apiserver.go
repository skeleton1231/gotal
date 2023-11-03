// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skeleton1231/gotal/internal/pkg/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
)

// APIServer wraps the gin.Engine with specific configurations and capabilities.
type APIServer struct {
	middlewares         []string
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo
	ShutdownTimeout     time.Duration
	*gin.Engine
	healthz, enableMetrics, enableProfiling bool
	insecureServer, secureServer            *http.Server
}

// initAPIServer initializes the API server with necessary settings and middlewares.
func initAPIServer(s *APIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// InstallAPIs installs specific endpoints to the server based on its configuration.
func (s *APIServer) InstallAPIs() {
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})
	}

	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	if s.enableProfiling {
		pprof.Register(s.Engine)
	}
}

// Setup customizes gin settings, mainly for debugging purposes.
func (s *APIServer) Setup() {
	// Override the DebugPrintRouteFunc for customized route logging
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		// Now this log statement is within the scope where the variables are defined
		logrus.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// Rest of your setup code...
}

// InstallMiddlewares sets up any global middlewares for the server.
func (s *APIServer) InstallMiddlewares() {
	// Potential spot to install any necessary middlewares

	// install custom middlewares
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			logrus.Warnf("can not find middleware: %s", m)

			continue
		}

		logrus.Infof("install middleware: %s", m)
		s.Use(mw)
	}
}

// Run starts the API server. It sets up and runs both the insecure and secure servers.
func (s *APIServer) Run() error {
	// Setup for insecure server
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}

	// Setup for secure server
	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
	}

	var eg errgroup.Group

	// Start insecure server
	eg.Go(func() error {
		logrus.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(err.Error())
			return err
		}
		logrus.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	// Start secure server
	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}

		logrus.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())

		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(err.Error())

			return err
		}

		logrus.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// Check server health (if enabled) before continuing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		logrus.Fatal(err.Error())
	}

	return nil
}

// Close gracefully shuts down both the insecure and secure servers.
func (s *APIServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()
	_ = s.secureServer.Shutdown(ctx)
	_ = s.insecureServer.Shutdown(ctx)
}

// ping checks the health of the server by sending a request to the /healthz endpoint.
func (s *APIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		resp, err := http.Get(url) // simplified from creating a new request
		if err == nil && resp.StatusCode == http.StatusOK {
			logrus.Info("The router has been deployed successfully.")
			resp.Body.Close()
			return nil
		}

		logrus.Info("Waiting for the router, if fail will retry in 1 second.")

		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			logrus.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}
