// user_grpc_service.go
package rpc_service

import (
	"context"

	"github.com/skeleton1231/gotal/internal/apiserver/store"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	pbO "github.com/skeleton1231/gotal/internal/proto/options"
	pb "github.com/skeleton1231/gotal/internal/proto/user"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// userGrpcServiceImpl 实现 UserStore 接口
type userGrpcServiceImpl struct {
	client pb.UserServiceClient
}

// NewUserGrpcService 创建一个新的实例
func newUser(ds *datastore) store.UserStore {
	return &userGrpcServiceImpl{ds.client}
}

func (s *userGrpcServiceImpl) Create(ctx context.Context, user *model.User, opts model.CreateOptions) error {

	// 创建CreateOptions的实例
	createOpts := &pbO.CreateOptions{}

	pbUser := model.UserToProto(user) // 转换为Protobuf格式
	_, err := s.client.Create(ctx, &pb.CreateRequest{
		User:    pbUser, // 使用转换后的用户信息
		Options: createOpts,
	})
	return err
}

func (s *userGrpcServiceImpl) Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error {
	udpateOpts := &pbO.UpdateOptions{}
	pbUser := model.UserToProto(user) // 转换为Protobuf格式
	_, err := s.client.Update(ctx, &pb.UpdateRequest{
		User:    pbUser, // 使用转换后的用户信息
		Options: udpateOpts,
	})
	return err
}

func (s *userGrpcServiceImpl) Delete(ctx context.Context, userId uint64, opts model.DeleteOptions) error {
	_, err := s.client.Delete(ctx, &pb.DeleteRequest{})
	return err
}

func (s *userGrpcServiceImpl) Get(ctx context.Context, userId uint64, opts model.GetOptions) (*model.User, error) {
	var user *model.User
	getOpts := &pbO.GetOptions{}
	userPb, err := s.client.Get(ctx, &pb.GetRequest{
		UserId:  userId,
		Options: getOpts,
		//
	})
	pbUser := userPb.GetUser()
	user, _ = model.ProtoToUser(pbUser)
	return user, err
}

func (s *userGrpcServiceImpl) GetByUsername(ctx context.Context, username string, opts model.GetOptions) (*model.User, error) {
	var user *model.User

	userPb, err := s.client.GetByUsername(ctx, &pb.GetByUsernameRequest{
		Username: username,
	})
	pb := userPb.GetUser()
	user.ID = pb.Meta.GetId()
	user.Name = pb.GetName()
	user.DiscordID = pb.GetDiscordId()
	user.Email = pb.GetEmail()
	return user, err
}

func (s *userGrpcServiceImpl) List(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	// 创建 pb.ListOptions 实例
	pbOpts := &pbO.ListOptions{
		LabelSelector: wrapperspb.String(opts.LabelSelector),
		FieldSelector: wrapperspb.String(opts.FieldSelector),
	}

	// 对于 Limit 和 Offset，只有在它们不为 nil 时才设置
	if opts.Limit != nil {
		pbOpts.Limit = wrapperspb.Int64(*opts.Limit)
	}
	if opts.Offset != nil {
		pbOpts.Offset = wrapperspb.Int64(*opts.Offset)
	}

	// 构造 ListRequest，包括转换后的 ListOptions
	pbReq := &pb.ListRequest{
		Options: pbOpts,
	}

	// 调用 List 方法
	pbList, err := s.client.List(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	// 将 protobuf 返回的 UserList 转换为 model.UserList
	userList := &model.UserList{
		Items: []*model.User{},
	}
	for _, pbUser := range pbList.GetUsers().Items {
		user, err := model.ProtoToUser(pbUser)
		if err != nil {
			// 处理转换错误
			continue
		}
		userList.Items = append(userList.Items, user)
	}

	return userList, nil
}

func (s *userGrpcServiceImpl) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return s.client.ChangePassword(ctx, req)
}
