package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New()

	// Use json tag names in error messages instead of Go field names.
	// e.g. "Email" → "email", "NewPassword" → "new_password"
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Check validates a struct and returns a map of field → human message.
// Returns nil if the struct is valid.
func Check(s any) map[string]string {
	if err := v.Struct(s); err != nil {
		errs := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = message(e)
		}
		return errs
	}
	return nil
}

func message(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "required"
	case "email":
		return "must be a valid email address"
	case "min":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters", e.Param())
		}
		return fmt.Sprintf("must be >= %s", e.Param())
	case "max":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("must be at most %s characters", e.Param())
		}
		return fmt.Sprintf("must be <= %s", e.Param())
	case "gte":
		return fmt.Sprintf("must be >= %s", e.Param())
	case "gt":
		return fmt.Sprintf("must be > %s", e.Param())
	default:
		return fmt.Sprintf("failed '%s' validation", e.Tag())
	}
}
