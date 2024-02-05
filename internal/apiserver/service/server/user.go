package service

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/skeleton1231/gotal/internal/apiserver/service"
	"github.com/skeleton1231/gotal/internal/apiserver/store"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
)

// userServiceServer 是 UserServiceServer 接口的实现
type UserServiceServer struct {
	// 这里可以包含服务器需要的任何状态或依赖
	store store.Factory
	pb.UnimplementedUserServiceServer
}

var (
	userServer *UserServiceServer
	once       sync.Once
)

// GetCacheInsOr return cache server instance with given factory.
func GetUserInsOr(store store.Factory) (*UserServiceServer, error) {
	if store != nil {
		once.Do(func() {
			userServer = &UserServiceServer{store: store}
		})
	}

	if userServer == nil {
		return nil, fmt.Errorf("got nil user server")
	}

	return userServer, nil
}

// 实现 Create 方法
func (s *UserServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	// 在这里编写 Create 方法的具体实现
	s.store.Users().Create(ctx, &model.User{}, model.CreateOptions{})
	return &pb.CreateResponse{}, nil
}

// 实现 Update 方法
func (s *UserServiceServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	// 在这里编写 Update 方法的具体实现
	return &pb.UpdateResponse{}, nil
}

// 实现 Delete 方法
func (s *UserServiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	// 在这里编写 Delete 方法的具体实现
	return &pb.DeleteResponse{}, nil
}

// 实现 Get 方法
func (s *UserServiceServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	// 在这里编写 Get 方法的具体实现
	return &pb.GetResponse{}, nil
}

// 实现 List 方法
func (s *UserServiceServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	// 在这里编写 List 方法的具体实现
	return &pb.ListResponse{}, nil
}

// 实现 ChangePassword 方法
func (s *UserServiceServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// 在这里编写 ChangePassword 方法的具体实现
	return &pb.ChangePasswordResponse{}, nil
}
