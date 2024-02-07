package service

import (
	pb "github.com/skeleton1231/gotal/internal/proto/user"
)

// Service 定义了服务层的接口，它抽象了服务的功能。
type RpcService interface {
	Users() UserGrpcSrv // 返回处理用户相关操作的UserSrv实例。
}

// service 结构实现了Service接口。
// 它持有gRPC客户端的引用，用于访问远程服务。
type rpcService struct {
	userClient pb.UserServiceClient // gRPC用户服务客户端
}

// NewService 是创建service实例的构造函数。
// 它接受gRPC客户端作为参数，并返回一个Service实例。
func NewRpcService(userClient pb.UserServiceClient) RpcService {
	return &rpcService{
		userClient: userClient, // 使用提供的gRPC客户端初始化userClient字段。
	}
}

// Users 方法返回UserSrv的一个实例，该实例通过gRPC客户端与远程用户服务进行通信。
func (s *rpcService) Users() UserGrpcSrv {
	// 此处将gRPC客户端传递给UserSrv的实现
	return newUserGrpcService(s.userClient)
}
