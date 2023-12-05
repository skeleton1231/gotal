package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/response"
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

	// Validate the user struct
	if err := validate.Struct(user); err != nil {
		// Handle validation errors
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		response.WriteResponse(c, errors.WithCode(code.ErrValidation, "Validation error"), validationErrors)
		return
	}

	defaultTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	user.Password, _ = common.Encrypt(user.Password)
	user.Status = 1
	user.EmailVerifiedAt = defaultTime
	user.TrialEndsAt = defaultTime

	if err := u.srv.Users().Create(c, &user, model.CreateOptions{}); err != nil {
		response.WriteResponse(c, err, nil)

		return
	}

	response.WriteResponse(c, nil, user)

}
