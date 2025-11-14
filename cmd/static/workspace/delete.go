package workspace

import (
	"fmt"

	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	flags "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/spf13/cobra"
)

func DeleteCmd(parent *cobra.Command) *cobra.Command {
	var name string
	var nameFlag *flags.StrFlag
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a workspace",
		Long:  "Delete a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace := parent.Context().Value(cmdutils.CXT_WORKSPACE_KEY).(workspace.Workspace)
			if len(args) > 0 {
				cmd.Flags().Set("name", args[0])
			}
			if nameFlag.IsChanged() {
				name = *nameFlag.Value
			} else {
				return cmdutils.NewCliError("workspace name is required")
			}
			err := workspace.Delete(name)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}
			fmt.Println("Workspace deleted successfully")
			return nil
		},
	}
	nameFlag = flags.NewStr(cmd, "name", "", "Name of the workspace")
	return cmd
}
