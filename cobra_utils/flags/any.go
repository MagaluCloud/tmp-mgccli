package cobrautils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

type AnyFlag struct {
	baseFlag
	Value *any
}

func (a *AnyFlag) Set(val string) error {
	if len(val) > 0 && (val[0] == '{' || val[0] == '[') {
		var jsonValue any
		if err := json.Unmarshal([]byte(val), &jsonValue); err == nil {
			*a.Value = jsonValue
			return nil
		}
	}

	if val == "true" {
		*a.Value = true
		return nil
	}
	if val == "false" {
		*a.Value = false
		return nil
	}

	if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
		*a.Value = intVal
		return nil
	}

	if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		*a.Value = floatVal
		return nil
	}

	*a.Value = val
	return nil
}

func (a *AnyFlag) String() string {
	if a.Value == nil {
		return ""
	}

	if val := *a.Value; val != nil {
		if b, err := json.Marshal(val); err == nil {
			if len(b) > 0 && (b[0] == '{' || b[0] == '[') {
				return string(b)
			}
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (a *AnyFlag) Type() string {
	return "any"
}

func NewAny(cmd *cobra.Command, name string, _ any, usage string) *AnyFlag {
	var value any = nil
	flag := &AnyFlag{baseFlag: baseFlag{cmd, name}, Value: &value}
	cmd.Flags().Var(flag, name, usage)
	return flag
}

func NewAnyP(cmd *cobra.Command, name string, shorthand string, _ any, usage string) *AnyFlag {
	var value any = nil
	flag := &AnyFlag{baseFlag: baseFlag{cmd, name}, Value: &value}
	cmd.Flags().VarP(flag, name, shorthand, usage)
	return flag
}
