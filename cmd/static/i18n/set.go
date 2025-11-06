package i18n

import (
	"fmt"

	"github.com/magaluCloud/mgccli/i18n"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func setCmd() *cobra.Command {
	manager := i18n.GetInstance()
	return &cobra.Command{
		Use:   "set [código]",
		Short: manager.T("cli.i18n.set.short"),
		Long:  manager.T("cli.i18n.set.long"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			code := args[0]

			// Verificar se o idioma existe
			info, err := manager.GetLanguageInfo(code)
			if err != nil {
				return fmt.Errorf(manager.T("cli.i18n.set.error"), code)
			}

			// Definir o idioma
			if err := manager.SetLanguage(code); err != nil {
				return err
			}

			successColor := color.New(color.FgGreen, color.Bold)
			successColor.Printf(manager.T("cli.i18n.set.success"), info.NativeName, code)

			// Mostrar como persistir a configuração
			fmt.Println()
			noteColor := color.New(color.FgYellow)
			noteColor.Println(manager.T("cli.i18n.set.note"))
			fmt.Println(manager.T("cli.i18n.set.note_1"))
			fmt.Printf(manager.T("cli.i18n.set.note_2"), code)
			fmt.Println(manager.T("cli.i18n.set.note_3"))
			fmt.Printf(manager.T("cli.i18n.set.note_4"), code)

			return nil
		},
	}
}
