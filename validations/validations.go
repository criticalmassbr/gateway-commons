package validations

import (
	"net/mail"
	"reflect"

	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -destination=../mocks/validator.go -package=mocks github.com/go-playground/validator/v10 FieldLevel

func rfc5322EmailValidator(fl validator.FieldLevel) bool {
	st := fl.Field()

	if st.Kind() != reflect.String {
		return false
	}

	_, err := mail.ParseAddress(st.String())
	return err == nil
}

func New() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("emailRFC5322", rfc5322EmailValidator)

	return validate
}
