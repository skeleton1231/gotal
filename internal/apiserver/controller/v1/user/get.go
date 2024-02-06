package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/response"
)

func (u *UserController) Get(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
	}
	user, err := u.srv.Users().Get(c, id, model.GetOptions{})
	// user, err := store.Client().Users().Get(c, id, model.GetOptions{})
	if err != nil {
		response.WriteResponse(c, err, nil)
	}

	response.WriteResponse(c, nil, user)
}
