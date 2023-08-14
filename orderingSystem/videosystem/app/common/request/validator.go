package request

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Getmessage() ValidatorMessages
}

type ValidatorMessages map[string]string

func GetErrorMsg(request interface{}, err error) string {
	if _, isValidatorError := err.(validator.ValidationErrors); isValidatorError {
		_, isValidator := request.(Validator)
		for _, v := range err.(validator.ValidationErrors) {
			if isValidator {

				if message, exit := request.(Validator).Getmessage()[v.Field()+":"+v.Tag()]; exit {
					return message
				}
			}
			return v.Error()
		}
	}
	return "Parameter error"
}
