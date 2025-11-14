package validator

import (
	"reflect"

	integer "github.com/magaluCloud/mgccli/cmd/common/validator/int"
	str "github.com/magaluCloud/mgccli/cmd/common/validator/str"
)

type Validator interface {
	integer() error
	str() error
	Validate() error
}

type validator struct {
	value       any
	validateTag string
}

func NewValidator(value any, validateTag string) Validator {
	return &validator{value: value, validateTag: validateTag}
}

func (v *validator) Validate() error {
	//reflect on type of v.value and call the appropriate validator
	typ := reflect.TypeOf(v.value)
	switch typ.Kind() {
	case reflect.Int:
		return v.integer()
	case reflect.String:
		return v.str()
	}
	return nil
}

func (v *validator) integer() error {
	return integer.Validator(v.value.(int), v.validateTag)
}

func (v *validator) str() error {
	return str.Validator(v.value.(string), v.validateTag)
}
