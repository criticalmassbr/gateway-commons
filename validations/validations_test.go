package validations

import (
	"reflect"
	"testing"

	"github.com/criticalmassbr/gateway-commons/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRfc5322EmailValidator(t *testing.T) {
	validator := New()

	emptyEmail := ""
	validEmail := "Somos Dialog <jhon@gmail.com>"
	invalidEmail := "12345"

	expect := assert.New(t)

	type Email struct {
		Email string `validate:"emailRFC5322"`
	}

	expect.Nil(validator.Struct(Email{Email: validEmail}))
	expect.NotNil(validator.Struct(Email{Email: emptyEmail}))
	expect.NotNil(validator.Struct(Email{Email: invalidEmail}))

	type RequiredEmail struct {
		Email string `validate:"required,emailRFC5322"`
	}

	expect.Nil(validator.Struct(RequiredEmail{Email: validEmail}))
	expect.NotNil(validator.Struct(RequiredEmail{Email: emptyEmail}))
	expect.NotNil(validator.Struct(RequiredEmail{Email: invalidEmail}))

	type OptionalEmail struct {
		Email string `validate:"omitempty,emailRFC5322"`
	}

	expect.Nil(validator.Struct(OptionalEmail{Email: validEmail}))
	expect.Nil(validator.Struct(OptionalEmail{Email: emptyEmail}))
	expect.NotNil(validator.Struct(OptionalEmail{Email: invalidEmail}))

	type InvalidEmailField struct {
		Email int `validate:"emailRFC5322"`
	}

	expect.NotNil(validator.Struct(InvalidEmailField{Email: 1}))
	expect.NotNil(validator.Struct(InvalidEmailField{Email: 1}))
	expect.NotNil(validator.Struct(InvalidEmailField{Email: 1}))

	crtl := gomock.NewController(t)
	flMock := mocks.NewMockFieldLevel(crtl)

	flMock.EXPECT().Field().Return(reflect.ValueOf(1))
	expect.False(rfc5322EmailValidator(flMock), "should have failed as reflect returned a non string field")

	flMock.EXPECT().Field().Return(reflect.ValueOf(""))
	expect.False(rfc5322EmailValidator(flMock), "should have failed as reflect returned a non email string")

	flMock.EXPECT().Field().Return(reflect.ValueOf(validEmail))
	expect.True(rfc5322EmailValidator(flMock))
}
