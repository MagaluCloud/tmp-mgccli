package workspace

import (
	"fmt"

	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	flags "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/spf13/cobra"
)

func SetCmd(parent *cobra.Command) *cobra.Command {
	var name string
	var nameFlag *flags.StrFlag

	cmd := &cobra.Command{
		Use:   "set [name]",
		Short: "Set a workspace",
		Long:  "Set a workspace",

		RunE: func(cmd *cobra.Command, args []string) error {
			workspace := parent.Context().Value(cmdutils.CXT_WORKSPACE_KEY).(workspace.Workspace)
			if len(args) > 0 {
				cmd.Flags().Set("name", args[0])
			}
			if nameFlag.IsChanged() {
				name = *nameFlag.Value
			} else {
				return cmdutils.NewCliError("workspace name is required")
			} // CobraFlagsAssign

			err := workspace.Set(name)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}
			fmt.Printf("Workspace %s set successfully\n", name)
			return nil
		},
	}
	nameFlag = flags.NewStr(cmd, "name", "", "Name of the workspace")
	return cmd
}
