package user

import (
	"errors"

	"github.com/go-playground/validator/v10"
	srvv1 "github.com/skeleton1231/gotal/internal/apiserver/service/v1"
	"github.com/skeleton1231/gotal/internal/apiserver/store"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
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

func validateUser(user *model.User) ([]string, error) {
	if err := validate.Struct(user); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return validationErrors, errors.New("validation error")
	}
	return nil, nil
}
