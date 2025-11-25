package workspace

import (
	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func SelectCmd(parent *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "select [name]",
		Short: "Select a workspace",
		Long:  "Select a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace := parent.Context().Value(cmdutils.CXT_WORKSPACE_KEY).(workspace.Workspace)
			list, err := workspace.List()
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}

			options := []huh.Option[string]{}
			for _, w := range list {
				options = append(options, huh.NewOption(w.Name(), w.Name()))
			}

			selectedWorkspace := huh.NewSelect[string]()
			selectedWorkspace.Options(options...)
			err = selectedWorkspace.Run()
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}
			var name string
			selectedWorkspace.Value(&name)
			err = workspace.Set(name)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}
			return nil
		},
	}
}
