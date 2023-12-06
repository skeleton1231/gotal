package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/response"
	"github.com/skeleton1231/gotal/internal/pkg/validation"
	"github.com/skeleton1231/gotal/pkg/log"
)

func (u *UserController) Update(c *gin.Context) {
	log.Record(c).Info("update user function called.")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
	}

	var r model.User

	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	user, err := u.srv.Users().Get(c, id, model.GetOptions{})
	if err != nil {
		response.WriteResponse(c, err, nil)

		return
	}

	user.Email = r.Email
	user.Extend = r.Extend

	if validationErrors, err := validation.CheckModel(&r); err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), validationErrors)
		return
	}

	// Save changed fields.
	if err := u.srv.Users().Update(c, user, model.UpdateOptions{}); err != nil {
		response.WriteResponse(c, err, nil)

		return
	}

	response.WriteResponse(c, nil, user)
}
