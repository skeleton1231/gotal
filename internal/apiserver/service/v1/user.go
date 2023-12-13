package service

import (
	"context"
	"regexp"
	"sync"

	"github.com/skeleton1231/gotal/internal/apiserver/store"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/pkg/log"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *model.User, opts model.CreateOptions) error
	Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error
	Delete(ctx context.Context, userId uint64, opts model.DeleteOptions) error
	// DeleteCollection(ctx context.Context, userIds []uint64, opts model.DeleteOptions) error
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
// func (u *userService) DeleteCollection(ctx context.Context, userIds []uint64, opts model.DeleteOptions) error {
// 	if err := u.store.Users().DeleteCollection(ctx, userIds, opts); err != nil {
// 		return errors.WithCode(code.ErrDatabase, err.Error())
// 	}
// 	return nil
// }

// Get implements UserSrv.
func (u *userService) Get(ctx context.Context, userId uint64, opts model.GetOptions) (*model.User, error) {
	user, err := u.store.Users().Get(ctx, userId, opts)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List implements UserSrv.
// List returns user list in the storage. This function has a good performance.
func (u *userService) List(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	users, err := u.store.Users().List(ctx, opts)
	if err != nil {
		log.Record(ctx).Errorf("list users from storage failed: %s", err.Error())

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	var m sync.Map

	for _, user := range users.Items {
		wg.Add(1)

		go func(user *model.User) {
			defer wg.Done()

			m.Store(
				user.ID,
				&model.User{
					ObjectMeta: model.ObjectMeta{
						ID:        user.ID,
						Extend:    user.Extend,
						CreatedAt: user.CreatedAt,
						UpdatedAt: user.UpdatedAt,
						Status:    user.Status,
					},
					Name:         user.Name,
					Email:        user.Email,
					StripeID:     user.StripeID,
					DiscordID:    user.DiscordID,
					TotalCredits: user.TotalCredits,
				})
		}(user)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	infos := make([]*model.User, 0, len(users.Items))
	for _, user := range users.Items {
		info, _ := m.Load(user.ID)
		infos = append(infos, info.(*model.User))
	}

	log.Record(ctx).Debugf("get %d users from backend storage.", len(infos))

	return &model.UserList{ListMeta: users.ListMeta, Items: infos}, nil
}

// Update implements UserSrv.
func (u *userService) Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error {

	if err := u.store.Users().Update(ctx, user, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}
