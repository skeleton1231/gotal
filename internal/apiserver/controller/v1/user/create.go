package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/response"
	"github.com/skeleton1231/gotal/internal/pkg/validation"
	"github.com/skeleton1231/gotal/pkg/log"
	"github.com/skeleton1231/gotal/pkg/util/common"
)

func (u *UserController) Create(c *gin.Context) {
	log.Record(c).Info("user create function called.")

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	user.Password, _ = common.Encrypt(c.Param("password"))
	user.Status = 1
	defaultTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	user.EmailVerifiedAt = defaultTime
	user.TrialEndsAt = defaultTime

	if validationErrors, err := validation.CheckModel(&user); err != nil {
		response.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), validationErrors)
		return
	}

	if err := u.srv.Users().Create(c, &user, model.CreateOptions{}); err != nil {
		response.WriteResponse(c, err, nil)

		return
	}

	token, _, err := generateJWTToken(&user)
	if err != nil {
		log.Errorf("generateJWTToken error is %s", err.Error())
	}
	user.Token = token
	response.WriteResponse(c, nil, user)

}
