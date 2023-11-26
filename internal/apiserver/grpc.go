// Copyright 2020 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/skeleton1231/gotal/pkg/log"
	"google.golang.org/grpc"
)

type grpcAPIServer struct {
	*grpc.Server
	address string
}

func (s *grpcAPIServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		logrus.Fatalf("failed to listen: %s", err.Error())
	}

	go func() {
		if err := s.Serve(listen); err != nil {
			logrus.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	log.Infof("start grpc server at %s", s.address)
}

func (s *grpcAPIServer) Close() {
	s.GracefulStop()
	log.Infof("GRPC server on %s stopped", s.address)
}
