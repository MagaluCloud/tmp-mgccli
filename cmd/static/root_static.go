package static

import (
	"github.com/magaluCloud/mgccli/cmd/static/auth"
	"github.com/magaluCloud/mgccli/cmd/static/config"
	"github.com/magaluCloud/mgccli/cmd/static/i18n"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func RootStatic(parent *cobra.Command, sdkCoreConfig sdk.CoreClient) {
	i18n.I18nCmd(parent)
	config.ConfigCmd(parent, sdkCoreConfig)
	auth.AuthCmd(parent.Context(), parent, sdkCoreConfig)
}
