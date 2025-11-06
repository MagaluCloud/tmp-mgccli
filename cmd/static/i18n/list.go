package i18n

import (
	"fmt"

	"github.com/magaluCloud/mgccli/i18n"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	manager := i18n.GetInstance()
	return &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.i18n.list.short"),
		Long:  manager.T("cli.i18n.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			languages := manager.GetAvailableLanguages()
			currentLang := manager.GetLanguage()

			if len(languages) == 0 {
				fmt.Println(manager.T("cli.i18n.list.no_languages"))
				return nil
			}

			headerColor := color.New(color.FgCyan, color.Bold)
			headerColor.Println(manager.T("cli.i18n.list.title"))
			fmt.Println()

			for _, code := range languages {
				info, err := manager.GetLanguageInfo(code)
				if err != nil {
					continue
				}

				// Destacar idioma atual
				if code == currentLang {
					currentColor := color.New(color.FgGreen, color.Bold)
					currentColor.Printf("  ✓ %s (%s)\n", info.NativeName, code)
				} else {
					fmt.Printf("    %s (%s)\n", info.NativeName, code)
				}
			}

			fmt.Println()
			noteColor := color.New(color.FgYellow)
			noteColor.Printf(manager.T("cli.i18n.list.current_language"), currentLang)
			noteColor.Println(manager.T("cli.i18n.list.use_set"))

			return nil
		},
	}
}
