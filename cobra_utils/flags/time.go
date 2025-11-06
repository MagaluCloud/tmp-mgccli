package cobrautils

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// TimeFlag permite parsear valores de data/hora em múltiplos formatos
// Suporta: RFC3339, ISO8601, timestamps Unix (segundos/milissegundos), e formatos customizados
type TimeFlag struct {
	baseFlag
	Value   *time.Time
	Formats []string // Formatos customizados adicionais para tentar
}

// Formatos padrão suportados
var defaultTimeFormats = []string{
	time.RFC3339,           // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano,       // "2006-01-02T15:04:05.999999999Z07:00"
	"2006-01-02T15:04:05",  // ISO8601 sem timezone
	"2006-01-02 15:04:05",  // Formato comum
	"2006-01-02T15:04:05Z", // ISO8601 UTC
	time.DateTime,          // "2006-01-02 15:04:05" (Go 1.20+)
	time.DateOnly,          // "2006-01-02" (Go 1.20+)
	time.TimeOnly,          // "15:04:05" (Go 1.20+)
	"2006-01-02",           // Data apenas
	"15:04:05",             // Hora apenas
	"15:04",                // Hora:minuto
	"02/01/2006",           // Formato brasileiro DD/MM/YYYY
	"02/01/2006 15:04:05",  // Formato brasileiro com hora
	time.RFC1123,           // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,          // "Mon, 02 Jan 2006 15:04:05 -0700"
}

// Set faz o parse do valor informado
func (t *TimeFlag) Set(val string) error {
	// Primeiro, tenta parsear como timestamp Unix (segundos)
	if ts, err := parseUnixTimestamp(val); err == nil {
		*t.Value = ts
		return nil
	}

	// Tenta formatos customizados primeiro (se fornecidos)
	for _, format := range t.Formats {
		if parsed, err := time.Parse(format, val); err == nil {
			*t.Value = parsed
			return nil
		}
	}

	// Tenta formatos padrão
	for _, format := range defaultTimeFormats {
		if parsed, err := time.Parse(format, val); err == nil {
			*t.Value = parsed
			return nil
		}
	}

	return fmt.Errorf("invalid time format: %s (supported: RFC3339, ISO8601, Unix timestamp, YYYY-MM-DD, etc)", val)
}

// String retorna a representação em string do valor
func (t *TimeFlag) String() string {
	if t.Value == nil || t.Value.IsZero() {
		return ""
	}
	return t.Value.Format(time.RFC3339)
}

// Type retorna o tipo para pflag
func (t *TimeFlag) Type() string {
	return "time"
}

// parseUnixTimestamp tenta parsear como timestamp Unix (segundos ou milissegundos)
func parseUnixTimestamp(val string) (time.Time, error) {
	var timestamp int64

	// Tenta parsear como inteiro
	_, err := fmt.Sscanf(val, "%d", &timestamp)
	if err != nil {
		return time.Time{}, err
	}

	// Se o timestamp é muito grande, assume milissegundos
	if timestamp > 1e12 {
		return time.UnixMilli(timestamp), nil
	}

	// Caso contrário, assume segundos
	return time.Unix(timestamp, 0), nil
}

// NewTimeFlag cria uma flag de tempo
func NewTime(cmd *cobra.Command, name string, usage string, _ string) *TimeFlag {
	value := &time.Time{}
	tv := &TimeFlag{
		baseFlag: baseFlag{cmd, name},
		Value:    value,
	}
	cmd.Flags().Var(tv, name, usage)
	return tv
}

// NewTimeFlagP cria uma flag de tempo com shorthand
func NewTimeP(cmd *cobra.Command, name string, shorthand string, usage string, _ string) *TimeFlag {
	value := &time.Time{}
	tv := &TimeFlag{
		baseFlag: baseFlag{cmd, name},
		Value:    value,
	}
	cmd.Flags().VarP(tv, name, shorthand, usage)
	return tv
}

// NewTimeFlagWithFormats cria uma flag de tempo com formatos customizados adicionais
func NewTimeWithFormats(cmd *cobra.Command, name string, usage string, formats []string) *TimeFlag {
	value := &time.Time{}
	tv := &TimeFlag{
		baseFlag: baseFlag{cmd, name},
		Value:    value,
		Formats:  formats,
	}
	cmd.Flags().Var(tv, name, usage)
	return tv
}

// NewTimeFlagWithFormatsP cria uma flag de tempo com formatos customizados e shorthand
func NewTimeWithFormatsP(cmd *cobra.Command, name string, shorthand string, usage string, formats []string) *TimeFlag {
	value := &time.Time{}
	tv := &TimeFlag{
		baseFlag: baseFlag{cmd, name},
		Value:    value,
		Formats:  formats,
	}
	cmd.Flags().VarP(tv, name, shorthand, usage)
	return tv
}
