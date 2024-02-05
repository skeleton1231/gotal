// user_grpc_service.go
package service

import (
	"context"

	pb "github.com/skeleton1231/gotal/internal/apiserver/service"
)

// UserGrpcService 定义了新的用户服务接口，使用gRPC进行通信
type UserGrpcService interface {
	Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error)
	Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
	Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error)
	List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error)
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
}

// userGrpcServiceImpl 实现 UserGrpcService 接口
type userGrpcServiceImpl struct {
	client pb.UserServiceClient
}

// NewUserGrpcService 创建一个新的 UserGrpcService 实例
func NewUserGrpcService(client pb.UserServiceClient) UserGrpcService {
	return &userGrpcServiceImpl{
		client: client,
	}
}

func (s *userGrpcServiceImpl) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	return s.client.Create(ctx, req)
}

func (s *userGrpcServiceImpl) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return s.client.Update(ctx, req)
}

func (s *userGrpcServiceImpl) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return s.client.Delete(ctx, req)
}

func (s *userGrpcServiceImpl) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	return s.client.Get(ctx, req)
}

func (s *userGrpcServiceImpl) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	return s.client.List(ctx, req)
}

func (s *userGrpcServiceImpl) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return s.client.ChangePassword(ctx, req)
}
