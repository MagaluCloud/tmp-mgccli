package static

import (
	"github.com/magaluCloud/mgccli/cmd/static/auth"
	"github.com/magaluCloud/mgccli/cmd/static/config"
	"github.com/magaluCloud/mgccli/cmd/static/workspace"
	"github.com/spf13/cobra"
)

func RootStatic(parent *cobra.Command) {
	config.ConfigCmd(parent)
	auth.AuthCmd(parent)
	workspace.WorkspaceCmd(parent)
}
