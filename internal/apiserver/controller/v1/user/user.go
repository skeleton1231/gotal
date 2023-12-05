package user

import (
	"github.com/go-playground/validator/v10"
	srvv1 "github.com/skeleton1231/gotal/internal/apiserver/service/v1"
	"github.com/skeleton1231/gotal/internal/apiserver/store"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserController struct {
	srv srvv1.Service
}

func NewUserController(store store.Factory) *UserController {
	return &UserController{
		srv: srvv1.NewService(store),
	}
}
