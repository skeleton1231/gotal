package database

import (
	"context"
	"fmt"

	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/pkg/log"
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
func (u *users) Delete(ctx context.Context, userId uint64, opts model.DeleteOptions) error {

	err := u.db.Where("id = ?", userId).Delete(&model.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(ctx context.Context, userIds []uint64, opts model.DeleteOptions) error {

	return u.db.Where("id in (?)", userIds).Delete(&model.User{}).Error
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, userId uint64, opts model.GetOptions) (*model.User, error) {
	user := &model.User{}
	err := u.db.Where("id = ? and status = 1", userId).First(&user).Error
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

	// Apply field selectors to the query
	query, err := userApplyFieldSelectors(u.db.Where("status = 1"), opts.FieldSelector)
	if err != nil {
		log.Errorf("user query error: %v", err)
		return nil, err // Return immediately if there's an error
	}

	// Apply pagination and execute the query
	d := query.
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

func userApplyFieldSelectors(query *gorm.DB, fieldSelector string) (*gorm.DB, error) {
	selector, err := fields.ParseSelector(fieldSelector)
	if err != nil {
		return nil, err
	}

	for _, req := range selector.Requirements() {
		switch req.Field {
		case "name":
			query = query.Where("name like ?", fmt.Sprintf("%%%s%%", req.Value))
		case "email":
			query = query.Where("email like ?", fmt.Sprintf("%%%s%%", req.Value))
		case "discordId":
			query = query.Where("discord_id = ?", req.Value)
		case "stripeId":
			query = query.Where("stripe_id = ?", req.Value)
		case "status":
			query = query.Where("status = ?", req.Value)
			// Add more cases as needed for other fields
		}
	}

	return query, nil
}
