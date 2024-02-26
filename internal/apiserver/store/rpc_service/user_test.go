package rpc_service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/skeleton1231/gotal/internal/apiserver/store/mocks"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	pb "github.com/skeleton1231/gotal/internal/proto/user"
	"github.com/stretchr/testify/assert"
)

func TestUserGrpcServiceImpl_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserServiceClient := mocks.NewMockUserServiceClient(mockCtrl)

	// 设置期望
	mockUserServiceClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(&pb.CreateResponse{}, nil). // 假设Create调用成功
		Times(1)

	// 创建服务实例
	ds := &datastore{client: mockUserServiceClient}
	userService := newUser(ds)

	// 调用服务方法
	user := &model.User{
		ObjectMeta: model.ObjectMeta{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		},
		Name:  "Test User",
		Email: "test@example.com",
	}
	err := userService.Create(context.Background(), user, model.CreateOptions{})

	// 验证
	assert.NoError(t, err)
}
