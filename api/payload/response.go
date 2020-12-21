package payload

import (
	"github.com/vanilla/go-pagination/api/dto"
	"github.com/vanilla/go-pagination/api/util"
	"gopkg.in/go-playground/validator.v8"
)

func GenerateValidationResponse(err error) (res dto.ValidationResponse) {
	res.Success = false

	var validations []dto.Validation

	// get validation
	validationError := err.(validator.ValidationErrors)

	for _, val := range validationError {
		// get field & rule
		field, rule := val.Field, val.Tag

		// create validation
		validation := dto.Validation{Field: field, Message: util.GenerateValidationMessage(field, rule)}

		// add validation
		validations = append(validations, validation)
	}

	// set validation
	res.Validations = validations
	return res
}
