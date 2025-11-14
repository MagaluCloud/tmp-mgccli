package workspace

import (
	"fmt"

	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func GetCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a workspace",
		Long:  "Get a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace := parent.Context().Value(cmdutils.CXT_WORKSPACE_KEY).(workspace.Workspace)
			fmt.Printf("Current workspace: %s\n", workspace.Get().Current().Name())
			return nil
		},
	}
	return cmd
}
