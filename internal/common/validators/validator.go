package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field '%s' failed on '%s' condition", e.Field(), e.Tag())
		}
	}
	return nil
}
