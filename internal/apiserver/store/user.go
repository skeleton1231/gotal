package store

import (
	"context"

	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *model.User, opts model.CreateOptions) error
	Update(ctx context.Context, user *model.User, opts model.UpdateOptions) error
	Delete(ctx context.Context, username string, opts model.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts model.DeleteOptions) error
	Get(ctx context.Context, username string, opts model.GetOptions) (*model.User, error)
	List(ctx context.Context, opts model.ListOptions) (*model.UserList, error)
}
