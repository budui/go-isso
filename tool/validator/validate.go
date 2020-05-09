package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// tValidator validate struct
type tValidator struct {
	v *validator.Validate
	t ut.Translator
}

var vd *tValidator

// newtValidator return a valid validator
func newtValidator() *tValidator {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, trans)

	return &tValidator{
		v: validate,
		t: trans,
	}
}

// Validate run the Validator
func (vd *tValidator) validate(data interface{}) error {
	if err := vd.v.Struct(data); err != nil {
		var errorString strings.Builder
		errs := err.(validator.ValidationErrors)

		errMsgs := errs.Translate(vd.t)

		for k, v := range errMsgs {
			errorString.WriteString(fmt.Sprintf("%s: %s", k, v))
		}

		return errors.New(errorString.String())
	}
	return nil
}

// Validate the struct
func Validate(data interface{}) error {
	if vd == nil {
		vd = newtValidator()
	}
	return vd.validate(data)
}
