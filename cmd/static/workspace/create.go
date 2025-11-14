package workspace

import (
	"fmt"

	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	flags "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/spf13/cobra"
)

func CreateCmd(parent *cobra.Command) *cobra.Command {
	var name string
	var nameFlag *flags.StrFlag
	var copyFlag *flags.StrFlag

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workspace",
		Long:  "Create a new workspace",
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

			if copyFlag.IsChanged() {
				err := workspace.Copy(*copyFlag.Value, *nameFlag.Value)
				if err != nil {
					return err
				}
				fmt.Println("Workspace copied successfully")
				return nil
			}

			err := workspace.Create(name)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}
			fmt.Println("Workspace created successfully")
			return nil
		},
	}

	nameFlag = flags.NewStr(cmd, "name", "", "Name of the workspace")
	copyFlag = flags.NewStr(cmd, "copy", "c", "Copy a workspace to a new workspace")

	return cmd
}
