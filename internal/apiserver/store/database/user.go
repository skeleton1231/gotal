package database

import (
	"context"

	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/fields"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *model.User, opts model.CreateOptions) error {
	return u.db.Create(&user).Error
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error {
	return u.db.Save(user).Error
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string, opts model.DeleteOptions) error {

	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	err := u.db.Where("name = ?", username).Delete(&model.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts model.DeleteOptions) error {

	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.Where("name in (?)", usernames).Delete(&model.User{}).Error
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts model.GetOptions) (*model.User, error) {
	user := &model.User{}
	err := u.db.Where("name = ? and status = 1", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}

// List return all users.
func (u *users) List(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	ret := &model.UserList{}
	ol := model.Unpointer(opts.Offset, opts.Limit)

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	d := u.db.Where("name like ? and status = 1", "%"+username+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

// ListOptional show a more graceful query method.
func (u *users) ListOptional(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	ret := &model.UserList{}
	ol := model.Unpointer(opts.Offset, opts.Limit)

	where := model.User{}

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, found := selector.RequiresExactMatch("name")
	if found {
		where.Name = username
	}

	d := u.db.Where(where).
		// Not(whereNot).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}
