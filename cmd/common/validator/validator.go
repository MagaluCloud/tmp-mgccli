package validator

import (
	"reflect"

	integer "github.com/magaluCloud/mgccli/cmd/common/validator/int"
)

type Validator interface {
	Integer() error
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
		return v.Integer()
	}
	return nil
}

func (v *validator) Integer() error {
	return integer.Validator(v.value.(int), v.validateTag)
}
