package user

import (
	srvv1 "github.com/skeleton1231/gotal/internal/apiserver/service/v1" // 假设的v2服务包路径
	"github.com/skeleton1231/gotal/internal/apiserver/store"
)

type UserController struct {
	srv srvv1.Service
}

func NewUserController(store store.Factory) *UserController {
	return &UserController{
		srv: srvv1.NewService(store),
	}
}
