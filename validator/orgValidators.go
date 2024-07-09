package validator

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (r OrgRegisterReq) Validate() ([]byte, error) {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 250)),
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

func (r OrgAddUserReq) Validate() ([]byte, error) {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.Required, validation.Length(3, 250)),
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
