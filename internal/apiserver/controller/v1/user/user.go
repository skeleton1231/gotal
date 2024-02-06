package user

import (
	srvv1 "github.com/skeleton1231/gotal/internal/apiserver/service/v1"
	srvv2 "github.com/skeleton1231/gotal/internal/apiserver/service/v2" // 假设的v2服务包路径
	"github.com/skeleton1231/gotal/internal/apiserver/store"
)

type UserController struct {
	srv    srvv1.Service
	rpcSrv srvv2.RpcService // 添加gRPC服务的引用

}

func NewUserController(store store.Factory, rpcSrv srvv2.RpcService) *UserController {
	return &UserController{
		srv:    srvv1.NewService(store),
		rpcSrv: rpcSrv, // 新增gRPC服务引用

	}
}
