package rpc_service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/skeleton1231/gotal/internal/apiserver/store"
	pb "github.com/skeleton1231/gotal/internal/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type datastore struct {
	client pb.UserServiceClient
}

// Close implements store.Factory.
func (ds *datastore) Close() error {
	return nil
}

func (ds *datastore) Users() store.UserStore {
	return newUser(ds)
}

var (
	rpcServerFactory store.Factory
	once             sync.Once
)

// GetRPCServerFactory returns a gRPC client factory with TLS.
// It connects to the server at the given address using the specified CA.
func GetRPCServerFactory(address string, clientCA string) (store.Factory, error) {
	var initErr error

	once.Do(func() {
		creds, err := credentials.NewClientTLSFromFile(clientCA, "")
		if err != nil {
			logrus.Errorf("credentials.NewClientTLSFromFile err: %v", err)
			initErr = err
			return
		}

		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
		if err != nil {
			logrus.Errorf("Connect to grpc server failed, error: %s", err)
			initErr = err
			return
		}

		client := pb.NewUserServiceClient(conn)
		rpcServerFactory = &datastore{client: client}
		logrus.Infof("Connected to grpc server, address: %s", address)
	})

	if initErr != nil {
		return nil, initErr
	}

	if rpcServerFactory == nil {
		errMsg := "failed to get rpcserver store factory"
		logrus.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	return rpcServerFactory, nil
}

// GetRPCServerFactoryNoTLS creates a gRPC client factory without enforcing TLS.
// It connects to the server at the given address.
func GetRPCServerFactoryNoTLS(serverAddr string) (store.Factory, error) {
	var initErr error

	once.Do(func() {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
		if err != nil {
			logrus.Errorf("Failed to dial: %v", err)
			initErr = err
			return
		}

		// Note: Do not close the connection here. It will be closed when the factory is closed.
		// defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		rpcServerFactory = &datastore{client: client}
	})

	if initErr != nil {
		return nil, initErr // Return the initialization error if any
	}

	if rpcServerFactory == nil {
		logrus.Errorf("failed to get rpcserver store factory")
		return nil, errors.New("failed to get rpcserver store factory")
	}

	return rpcServerFactory, nil
}
