package workspace

import (
	"github.com/spf13/cobra"
)

func WorkspaceCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "workspace",
		Short:   "Workspace",
		Long:    "Workspace",
		GroupID: "settings",
	}

	cmd.AddCommand(CreateCmd(parent))
	cmd.AddCommand(DeleteCmd(parent))
	cmd.AddCommand(GetCmd(parent))
	cmd.AddCommand(ListCmd(parent))
	cmd.AddCommand(SetCmd(parent))

	parent.AddCommand(cmd)
}
