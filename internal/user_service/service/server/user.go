package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	pb "github.com/skeleton1231/gotal/internal/proto/user"
	"github.com/skeleton1231/gotal/internal/user_service/store"
	"github.com/skeleton1231/gotal/pkg/log"
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
func (s *UserServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	obj := req.GetUser()
	// 先初始化 user 变量
	user := &model.User{
		Email:           obj.Email,
		Name:            obj.Name,
		EmailVerifiedAt: obj.GetEmailVerifiedAt().AsTime(),
		TrialEndsAt:     obj.GetTrialEndsAt().AsTime(),
		ObjectMeta: model.ObjectMeta{
			Status: int(obj.GetMeta().Status),
		},
	}

	log.Infof("EmailVerifiedAt: %+v\n", obj.GetEmailVerifiedAt().AsTime())
	log.Infof("TrialEndsAt: %+v\n", obj.GetTrialEndsAt().AsTime())

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
// Update 方法更新一个用户
func (s *UserServiceServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	// 从请求中提取用户ID
	userID := req.GetUser().GetMeta().GetId()

	// 从数据库或存储中检索现有的用户信息
	existingUser, err := s.store.Users().Get(ctx, userID, model.GetOptions{})
	if err != nil {
		// 处理错误，例如用户不存在的情况
		return nil, err
	}

	// 将请求中的新数据赋值到现有的用户对象上
	// 假设 ProtoToUser 是一个将pb.User转换为model.User的函数，并返回一个*model.User
	// 这个函数需要实现字段的合适映射和赋值
	updatedUser, err := model.ProtoToUser(req.GetUser())
	if err != nil {
		return nil, err
	}

	// 将更新后的数据赋值到现有的用户对象上，这里需要根据实际情况调整字段赋值
	if updatedUser.Name != "" {
		existingUser.Name = updatedUser.Name
	}
	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}
	// 更新其他需要更新的字段...

	// 使用更新后的用户信息进行更新操作
	err = s.store.Users().Update(ctx, existingUser, model.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	// 创建响应对象，这里假设更新操作成功后不需要返回更新后的用户信息
	// 如果需要，可以从存储中重新获取用户信息，并转换为Protobuf格式返回
	return &pb.UpdateResponse{
		User: model.UserToProto(existingUser),
	}, nil
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
		User: model.UserToProto(user),
	}

	// 返回 GetResponse 对象作为响应
	return response, nil
}

// 实现 List 方法
func (s *UserServiceServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	// 在这里编写 List 方法的具体实现
	return &pb.ListResponse{}, nil
}

func (s *UserServiceServer) GetByUsername(ctx context.Context, req *pb.GetByUsernameRequest) (*pb.GetByUsernameResponse, error) {
	user, err := s.store.Users().GetByUsername(ctx, req.GetUsername(), model.GetOptions{})
	if err != nil {
		// 处理错误，比如返回gRPC的错误码
		return nil, err
	}

	// 假设有一个函数UserToProto转换model.User到protobuf的User
	pbUser := model.UserToProto(user)

	return &pb.GetByUsernameResponse{
		User: pbUser,
	}, nil
}

// 实现 ChangePassword 方法
func (s *UserServiceServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// 在这里编写 ChangePassword 方法的具体实现
	return &pb.ChangePasswordResponse{}, nil
}
