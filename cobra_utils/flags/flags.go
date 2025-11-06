package cobrautils

import "github.com/spf13/cobra"

type baseFlag struct {
	cmd  *cobra.Command
	name string
}

type StrFlag struct {
	baseFlag
	Value *string
}

type Int64Flag struct {
	baseFlag
	Value *int64
}

type BoolFlag struct {
	baseFlag
	Value *bool
}

type IntFlag struct {
	baseFlag
	Value *int
}

type Float64Flag struct {
	baseFlag
	Value *float64
}

type StrSliceFlag struct {
	baseFlag
	Value *[]string
}

type StrMapFlag struct {
	baseFlag
	Value *map[string]string
}

func (f *baseFlag) IsChanged() bool {
	return f.cmd.Flags().Changed(f.name)
}

func NewStr(cmd *cobra.Command, name string, defaultValue string, usage string) *StrFlag {
	var value *string = new(string)
	cmd.Flags().StringVar(value, name, defaultValue, usage)
	return &StrFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewStrP(cmd *cobra.Command, name string, shorthand string, defaultValue string, usage string) *StrFlag {
	var value *string = new(string)
	cmd.Flags().StringVarP(value, name, shorthand, defaultValue, usage)
	return &StrFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewInt64(cmd *cobra.Command, name string, defaultValue int64, usage string) *Int64Flag {
	var value *int64 = new(int64)
	cmd.Flags().Int64Var(value, name, defaultValue, usage)
	return &Int64Flag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewInt64P(cmd *cobra.Command, name string, shorthand string, defaultValue int64, usage string) *Int64Flag {
	var value *int64 = new(int64)
	cmd.Flags().Int64VarP(value, name, shorthand, defaultValue, usage)
	return &Int64Flag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewBool(cmd *cobra.Command, name string, defaultValue bool, usage string) *BoolFlag {
	var value *bool = new(bool)
	cmd.Flags().BoolVar(value, name, defaultValue, usage)
	return &BoolFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewBoolP(cmd *cobra.Command, name string, shorthand string, defaultValue bool, usage string) *BoolFlag {
	var value *bool = new(bool)
	cmd.Flags().BoolVarP(value, name, shorthand, defaultValue, usage)
	return &BoolFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewInt(cmd *cobra.Command, name string, defaultValue int, usage string) *IntFlag {
	var value *int = new(int)
	cmd.Flags().IntVar(value, name, defaultValue, usage)
	return &IntFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewIntP(cmd *cobra.Command, name string, shorthand string, defaultValue int, usage string) *IntFlag {
	var value *int = new(int)
	cmd.Flags().IntVarP(value, name, shorthand, defaultValue, usage)
	return &IntFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewStrSlice(cmd *cobra.Command, name string, defaultValue []string, usage string) *StrSliceFlag {
	var value *[]string = new([]string)
	cmd.Flags().StringSliceVar(value, name, defaultValue, usage)
	return &StrSliceFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewStrSliceP(cmd *cobra.Command, name string, shorthand string, defaultValue []string, usage string) *StrSliceFlag {
	var value *[]string = new([]string)
	cmd.Flags().StringSliceVarP(value, name, shorthand, defaultValue, usage)
	return &StrSliceFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewStrMap(cmd *cobra.Command, name string, defaultValue map[string]string, usage string) *StrMapFlag {
	var value *map[string]string = new(map[string]string)
	cmd.Flags().StringToStringVar(value, name, defaultValue, usage)
	return &StrMapFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewStrMapP(cmd *cobra.Command, name string, shorthand string, defaultValue map[string]string, usage string) *StrMapFlag {
	var value *map[string]string = new(map[string]string)
	cmd.Flags().StringToStringVarP(value, name, shorthand, defaultValue, usage)
	return &StrMapFlag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewFloat64(cmd *cobra.Command, name string, defaultValue float64, usage string) *Float64Flag {
	var value *float64 = new(float64)
	cmd.Flags().Float64Var(value, name, defaultValue, usage)
	return &Float64Flag{baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewFloat64P(cmd *cobra.Command, name string, shorthand string, defaultValue float64, usage string) *Float64Flag {
	var value *float64 = new(float64)
	cmd.Flags().Float64VarP(value, name, shorthand, defaultValue, usage)
	return &Float64Flag{baseFlag: baseFlag{cmd, name}, Value: value}
}
