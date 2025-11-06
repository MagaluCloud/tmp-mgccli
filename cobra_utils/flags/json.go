package cobrautils

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// JSONValue é um tipo genérico para parsear flags JSON usando pflag/cobra
// Exemplo de uso:
//
//	var myFlag cobrautils.JSONValue[MyType]
//	flags.Var(&myFlag, "myflag", "Some JSON input")
type JSONValue[T any] struct {
	baseFlag
	Value *T
}

// Set faz o parse do JSON recebido na flag para o tipo T
func (j *JSONValue[T]) Set(val string) error {
	if err := json.Unmarshal([]byte(val), j.Value); err != nil {
		return fmt.Errorf("invalid JSON for flag: %w", err)
	}
	return nil
}

// String serializa o valor atual para JSON
func (j *JSONValue[T]) String() string {
	b, err := json.Marshal(*j.Value)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// Type retorna o nome do tipo para pflag
func (j *JSONValue[T]) Type() string {
	return "json"
}

func NewJSONValue[T any](cmd *cobra.Command, name string, usage string) *JSONValue[T] {
	var value *JSONValue[T] = new(JSONValue[T])
	cmd.Flags().Var(value, name, usage)
	return &JSONValue[T]{baseFlag: baseFlag{cmd, name}, Value: value.Value}
}

func NewJSONValueP[T any](cmd *cobra.Command, name string, shorthand string, usage string) *JSONValue[T] {
	var value *JSONValue[T] = new(JSONValue[T])
	cmd.Flags().VarP(value, name, shorthand, usage)
	return &JSONValue[T]{baseFlag: baseFlag{cmd, name}, Value: value.Value}
}
