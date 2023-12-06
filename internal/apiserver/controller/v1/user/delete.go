package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/response"
	"github.com/skeleton1231/gotal/pkg/log"
)

func (u *UserController) Delete(c *gin.Context) {
	log.Record(c).Info("delete user function called.")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
	}

	if err := u.srv.Users().Delete(c, id, model.DeleteOptions{Unscoped: true}); err != nil {
		response.WriteResponse(c, err, nil)

		return
	}

	response.WriteResponse(c, nil, nil)
}
