package workspace

import (
	"fmt"

	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func ListCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workspaces",
		Long:  "List workspaces",

		RunE: func(cmd *cobra.Command, args []string) error {
			workspace := parent.Context().Value(cmdutils.CXT_WORKSPACE_KEY).(workspace.Workspace)
			workspaces, err := workspace.List()
			if err != nil {
				return err
			}
			for _, workspace := range workspaces {
				fmt.Printf("Name: %s\n", workspace.Name())
			}
			return nil
		},
	}
	return cmd
}
