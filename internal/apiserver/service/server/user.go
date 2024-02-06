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

// GetUserInsOr return cache server instance with given factory.
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
// 实现 Create 方法
func (s *UserServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	obj := req.GetUser()
	// 先初始化 user 变量
	user := &model.User{
		Email: obj.Email,
		Name:  obj.Name,
	}
	// 然后调用 store 方法来创建用户
	err := s.store.Users().Create(ctx, user, model.CreateOptions{})
	if err != nil {
		// 处理创建过程中可能发生的错误
		return nil, err
	}
	// 如果创建成功，返回创建的用户信息
	return &pb.CreateResponse{User: obj}, nil
}

// 实现 Update 方法
func (s *UserServiceServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	// 从请求中获取用户信息
	pbUser := req.GetUser()

	// 创建一个空的 model.User 对象
	user := &model.User{}

	// 检查是否提供了有效的用户信息
	if pbUser != nil {
		// 如果提供了姓名字段且有值，则更新姓名
		if pbUser.Name != "" {
			user.Name = pbUser.Name
		}
		// 如果提供了邮箱字段且有值，则更新邮箱
		if pbUser.Email != "" {
			user.Email = pbUser.Email
		}
		// 其他字段的检查和更新操作
	}
	s.store.Users().Update(ctx, user, model.UpdateOptions{})
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
	// 从请求中获取用户的标识（假设是用户的 ID）
	userID := req.GetUserId()

	// 使用 userID 从数据库或存储中检索用户的信息
	user, err := s.store.Users().Get(ctx, userID, model.GetOptions{})
	if err != nil {
		// 处理错误，例如用户不存在的情况
		return nil, err
	}

	// 创建 GetResponse 对象，并将检索到的用户信息赋值给它
	response := &pb.GetResponse{
		User: &pb.User{
			// 在这里赋值用户信息，例如姓名、邮箱等
			Name:  user.Name,
			Email: user.Email,
			// 其他字段的赋值
		},
	}

	// 返回 GetResponse 对象作为响应
	return response, nil
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
