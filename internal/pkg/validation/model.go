package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CheckModel[T any](model T) ([]string, error) {
	if err := validate.Struct(model); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return validationErrors, err
	}
	return nil, nil
}
