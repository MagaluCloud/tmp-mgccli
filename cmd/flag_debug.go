package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

const debugLevelFlag = "debug"

func addLogDebugFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().String(
		debugLevelFlag,
		"debug",
		`Display detailed log information at the debug level`,
	)
}

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"

	LevelDebugLevel = -4
	LevelInfoLevel  = 0
	LevelWarnLevel  = 4
	LevelErrorLevel = 8
)

func parseDebugLevel(result string) int {
	switch strings.ToLower(result) {
	case LevelDebug:
		return LevelDebugLevel
	case LevelInfo:
		return LevelInfoLevel
	case LevelWarn:
		return LevelWarnLevel
	case LevelError:
		return LevelErrorLevel
	}

	return LevelDebugLevel
}
