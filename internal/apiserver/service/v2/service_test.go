package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	sserve1 "github.com/skeleton1231/gotal/internal/apiserver/service/server"
	"github.com/skeleton1231/gotal/internal/apiserver/store/mock_store"
	pb "github.com/skeleton1231/gotal/internal/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	ctrl := gomock.NewController(nil) // 注意: 这里应该传递*t.T，因为是init函数，我们用nil代替
	mockFactory := mock_store.NewMockFactory(ctrl)
	mockUserStore := mock_store.NewMockUserStore(ctrl)

	// 设置mock行为
	mockFactory.EXPECT().Users().Return(mockUserStore).AnyTimes()
	mockUserStore.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// 使用mockFactory代替真实的store实例
	userService, err := sserve1.GetUserInsOr(mockFactory)
	if err != nil || userService == nil {
		panic("UserService initialization failed: " + err.Error())
	}
	pb.RegisterUserServiceServer(s, userService)

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestUserService_Create(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.Create(ctx, &pb.CreateRequest{
		User: &pb.User{
			Name:  "Gonda1231",
			Email: "Gonda1231@gmail.com",
		},
	})
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}
	if resp.GetUser().GetName() != "Gonda1231" || resp.GetUser().GetEmail() != "Gonda1231@gmail.com" {
		t.Errorf("Unexpected response: got %v", resp.GetUser())
	}
}
