package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/user_service/store/mock_store"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {
	// Initialize the GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Asserts that all expected calls were made

	// Create a mock user store using the controller
	mockUserStore := mock_store.NewMockUserStore(ctrl)
	mockStoreFactory := mock_store.NewMockFactory(ctrl)

	// Set up your expectations as before
	user := &model.User{
		Name: "test user",
		// Other necessary user fields
	}
	mockUserStore.EXPECT().Create(gomock.Any(), user, gomock.Any()).Return(nil)
	mockStoreFactory.EXPECT().Users().Return(mockUserStore)

	// Create userService instance and execute test as before
	userService := NewService(mockStoreFactory)
	err := userService.Users().Create(context.Background(), user, model.CreateOptions{})

	// Verify results as before
	assert.NoError(t, err)
}

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mock_store.NewMockUserStore(ctrl)
	mockStoreFactory := mock_store.NewMockFactory(ctrl)

	user := &model.User{
		ObjectMeta: model.ObjectMeta{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		},
		Name: "John Doe",
	}
	mockUserStore.EXPECT().Update(gomock.Any(), user, gomock.Any()).Return(nil)
	mockStoreFactory.EXPECT().Users().Return(mockUserStore)

	userService := NewService(mockStoreFactory)
	err := userService.Users().Update(context.Background(), user, model.UpdateOptions{})

	assert.NoError(t, err)
}

func TestUserService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mock_store.NewMockUserStore(ctrl)
	mockStoreFactory := mock_store.NewMockFactory(ctrl)

	userId := uint64(1)
	mockUserStore.EXPECT().Delete(gomock.Any(), userId, gomock.Any()).Return(nil)
	mockStoreFactory.EXPECT().Users().Return(mockUserStore)

	userService := NewService(mockStoreFactory)
	err := userService.Users().Delete(context.Background(), userId, model.DeleteOptions{})

	assert.NoError(t, err)
}

func TestUserService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mock_store.NewMockUserStore(ctrl)
	mockStoreFactory := mock_store.NewMockFactory(ctrl)

	userId := uint64(1)
	expectedUser := &model.User{ObjectMeta: model.ObjectMeta{
		ID:        userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    1,
	}, Name: "John Doe"}

	mockUserStore.EXPECT().Get(gomock.Any(), userId, gomock.Any()).Return(expectedUser, nil)
	mockStoreFactory.EXPECT().Users().Return(mockUserStore)

	userService := NewService(mockStoreFactory)
	user, err := userService.Users().Get(context.Background(), userId, model.GetOptions{})

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mock_store.NewMockUserStore(ctrl)
	mockStoreFactory := mock_store.NewMockFactory(ctrl)

	expectedUsers := &model.UserList{
		Items: []*model.User{{ObjectMeta: model.ObjectMeta{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		}, Name: "John Doe"}},
	}
	mockUserStore.EXPECT().List(gomock.Any(), gomock.Any()).Return(expectedUsers, nil)
	mockStoreFactory.EXPECT().Users().Return(mockUserStore)

	userService := NewService(mockStoreFactory)
	users, err := userService.Users().List(context.Background(), model.ListOptions{})

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}
