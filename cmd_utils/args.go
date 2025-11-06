package cmdutils

import (
	"fmt"
	"os"
	"strings"
)

type (
	argsParser struct {
		allArgs []string
	}

	ArgsParser interface {
		FullProgramPath() string
		AllArgs() []string
		GetValue(key string) (string, bool, error)
		GetValueWithDefault(key string, defaultValue string) (string, bool, error)
	}
)

func NewArgsParser() ArgsParser {
	return &argsParser{allArgs: os.Args[1:]}
}

func (o *argsParser) FullProgramPath() string {
	return os.Args[0]
}

func (o *argsParser) AllArgs() []string {
	if o.allArgs == nil {
		o.allArgs = os.Args[1:]
	}
	return o.allArgs
}

func (o *argsParser) GetValue(key string) (string, bool, error) {
	for i, arg := range o.allArgs {
		isPresent := false
		if o.keyIsValid(arg) && (o.cutPrefix(arg) == key || strings.HasPrefix(o.cutPrefix(arg), key+"=")) {
			isPresent = true
			value, err := o.extractValue(arg)
			if err == nil {
				return value, isPresent, nil
			}

			// get next arg as key
			if i+1 < len(o.allArgs) {
				key, err := o.extractKey(o.allArgs[i+1])
				if err == nil {
					return key, isPresent, nil
				}
			}

			if len(o.allArgs) == i+1 {
				if key == "debug" {
					o.ApplyValue(key, "debug")
					return "debug", isPresent, nil
				}
				return "", isPresent, nil
			}

			return "", isPresent, fmt.Errorf("invalid key-value pair: %s", arg)

		}
	}
	return "", false, fmt.Errorf("key not found: %s", key)
}
func (o *argsParser) ApplyValue(key string, value string) error {
	for i, arg := range o.allArgs {
		if o.keyIsValid(arg) && (o.cutPrefix(arg) == key || strings.HasPrefix(o.cutPrefix(arg), key+"=")) {
			o.allArgs[i] = fmt.Sprintf("%s=%s", key, value)
			return nil
		}
	}
	return fmt.Errorf("key not found: %s", key)
}
func (o *argsParser) GetValueWithDefault(key string, defaultValue string) (string, bool, error) {
	value, isPresent, err := o.GetValue(key)
	if err == nil {
		return value, isPresent, nil
	}
	return defaultValue, false, nil
}

func (o *argsParser) extractValue(keyValue string) (string, error) {
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid key-value pair: %s", keyValue)
	}
	return parts[1], nil
}

func (o *argsParser) keyIsValid(key string) bool {
	return strings.HasPrefix(key, "--") || strings.HasPrefix(key, "-")
}

func (o *argsParser) cutPrefix(key string) string {
	if result, ok := strings.CutPrefix(key, "--"); ok {
		return result
	}
	if result, ok := strings.CutPrefix(key, "-"); ok {
		return result
	}
	return key
}

func (o *argsParser) extractKey(keyValue string) (string, error) {
	if o.keyIsValid(keyValue) {
		return "", fmt.Errorf("invalid key: %s", keyValue)
	}
	return keyValue, nil
}
