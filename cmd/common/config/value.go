package config

import "strconv"

type Value interface {
	String() string
	Bool() bool
	Int() int
}

type valuePtr struct {
	value any
}

func NewValue(value any) *valuePtr {
	return &valuePtr{value: value}
}

func (v *valuePtr) String() string {
	if v.value == nil {
		return ""
	}
	return anyToString(v.value)
}
func (v *valuePtr) Bool() bool {
	if v.value == nil {
		return false
	}
	return v.value.(bool)
}
func (v *valuePtr) Int() int {
	if v.value == nil {
		return 0
	}
	return v.value.(int)
}

func anyToString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}
