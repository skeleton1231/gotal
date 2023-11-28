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
	Delete(ctx context.Context, username string, opts model.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts model.DeleteOptions) error
	Get(ctx context.Context, username string, opts model.GetOptions) (*model.User, error)
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
	panic("unimplemented")
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
func (*userService) Delete(ctx context.Context, username string, opts model.DeleteOptions) error {
	panic("unimplemented")
}

// DeleteCollection implements UserSrv.
func (*userService) DeleteCollection(ctx context.Context, usernames []string, opts model.DeleteOptions) error {
	panic("unimplemented")
}

// Get implements UserSrv.
func (*userService) Get(ctx context.Context, username string, opts model.GetOptions) (*model.User, error) {
	panic("unimplemented")
}

// List implements UserSrv.
func (*userService) List(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	panic("unimplemented")
}

// Update implements UserSrv.
func (*userService) Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error {
	panic("unimplemented")
}
