package service

import (
	"context"

	pb "github.com/skeleton1231/gotal/internal/apiserver/service"

	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
)

// UserSrv 定义了处理用户请求的函数。
type UserGrpcSrv interface {
	Create(ctx context.Context, user *model.User, opts model.CreateOptions) error
	// Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error
	// Delete(ctx context.Context, userId uint64, opts model.DeleteOptions) error
	// // DeleteCollection(ctx context.Context, userIds []uint64, opts model.DeleteOptions) error
	// Get(ctx context.Context, userId uint64, opts model.GetOptions) (*model.User, error)
	// List(ctx context.Context, opts model.ListOptions) (*model.UserList, error)
	// ChangePassword(ctx context.Context, user *model.User) error
	// 定义其他方法...
}

var _ UserGrpcSrv = (*userGrpcService)(nil)

// userService 是UserSrv接口的实现。
type userGrpcService struct {
	client pb.UserServiceClient // gRPC用户服务客户端
}

// newUserService 创建一个新的userService实例。
func newUserGrpcService(client pb.UserServiceClient) UserGrpcSrv {
	return &userGrpcService{
		client: client,
	}
}

// Create 通过gRPC客户端调用远程创建用户服务。
func (u *userGrpcService) Create(ctx context.Context, user *model.User, opts model.CreateOptions) error {
	// 将model.User转换为protobuf定义的User类型
	pbUser := &pb.User{
		// 初始化字段...
	}
	_, err := u.client.Create(ctx, &pb.CreateRequest{
		User: pbUser,
	})
	return err
}

// 实现其他UserSrv接口方法...
