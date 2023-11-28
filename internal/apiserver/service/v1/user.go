package service

import (
	"context"
	"regexp"

	"github.com/skeleton1231/gotal/internal/apiserver/store"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *model.User, opts model.CreateOptions) error
	Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error
	Delete(ctx context.Context, userId uint64, opts model.DeleteOptions) error
	DeleteCollection(ctx context.Context, userIds []uint64, opts model.DeleteOptions) error
	Get(ctx context.Context, userId uint64, opts model.GetOptions) (*model.User, error)
	List(ctx context.Context, opts model.ListOptions) (*model.UserList, error)
	ChangePassword(ctx context.Context, user *model.User) error
}

type userService struct {
	store store.Factory
}

var _ UserSrv = (*userService)(nil)

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

// ChangePassword implements UserSrv.
func (u *userService) ChangePassword(ctx context.Context, user *model.User) error {
	// Save Password changed fields.
	if err := u.store.Users().Update(ctx, user, model.UpdateOptions{}); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// Create implements UserSrv.
func (u *userService) Create(ctx context.Context, user *model.User, opts model.CreateOptions) error {
	if err := u.store.Users().Create(ctx, user, opts); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'idx_name'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// Delete implements UserSrv.
func (u *userService) Delete(ctx context.Context, userId uint64, opts model.DeleteOptions) error {
	if err := u.store.Users().Delete(ctx, userId, opts); err != nil {
		return err
	}
	return nil
}

// DeleteCollection implements UserSrv.
func (u *userService) DeleteCollection(ctx context.Context, userIds []uint64, opts model.DeleteOptions) error {
	if err := u.store.Users().DeleteCollection(ctx, userIds, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

// Get implements UserSrv.
func (u *userService) Get(ctx context.Context, userId uint64, opts model.GetOptions) (*model.User, error) {
	user, err := u.store.Users().Get(ctx, userId, opts)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List implements UserSrv.
func (u *userService) List(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	panic("unimplemented")
}

// Update implements UserSrv.
func (u *userService) Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error {

	if err := u.store.Users().Update(ctx, user, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}
