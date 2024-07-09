package validator

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (r LoginReq) Validate() ([]byte, error) {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 30)),
	)

	if err == nil {
		return nil, nil
	}

	var FieldErrs []FieldErr

	if ve, ok := err.(validation.Errors); ok {
		for field, err := range ve {
			FieldErrs = append(FieldErrs, FieldErr{
				Field:   field,
				Message: err.Error(),
			})
		}
	}

	// marshal the map to JSON
	erroJSON, JsonErr := json.Marshal(FieldErrs)
	if JsonErr != nil {
		return nil, JsonErr
	}

	return erroJSON, nil
}
