package user

import (
	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/response"
	"github.com/skeleton1231/gotal/pkg/log"
)

// List list the users in the storage.
// Only administrator can call this function.
func (u *UserController) List(c *gin.Context) {
	log.Record(c).Info("list user function called.")

	var r model.ListOptions
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	users, err := u.srv.Users().List(c, r)
	if err != nil {
		response.WriteResponse(c, err, nil)

		return
	}

	response.WriteResponse(c, nil, users)
}
