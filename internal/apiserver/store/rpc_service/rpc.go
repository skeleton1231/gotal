package rpc_service

import (
	"sync"

	"github.com/skeleton1231/gotal/internal/apiserver/store"
	pb "github.com/skeleton1231/gotal/internal/proto/user"
	"github.com/skeleton1231/gotal/pkg/log"
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

// GetAPIServerFactoryOrDie return cache instance and panics on any error.
func GetRPCServerFactoryOrDie(address string, clientCA string) store.Factory {
	once.Do(func() {
		var (
			err   error
			conn  *grpc.ClientConn
			creds credentials.TransportCredentials
		)

		creds, err = credentials.NewClientTLSFromFile(clientCA, "")
		if err != nil {
			log.Panicf("credentials.NewClientTLSFromFile err: %v", err)
		}

		conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Panicf("Connect to grpc server failed, error: %s", err.Error())
		}

		rpcServerFactory = &datastore{pb.NewUserServiceClient(conn)}
		log.Infof("Connected to grpc server, address: %s", address)
	})

	if rpcServerFactory == nil {
		log.Panicf("failed to get rpcserver store fatory")
	}

	return rpcServerFactory
}
