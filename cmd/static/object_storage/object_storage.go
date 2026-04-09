package objectstorage

import (
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ObjectStorageCmd cria e configura o comando de object storage
func ObjectStorageCmd(parent *cobra.Command) {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "object-storage",
		Short:   manager.T("cli.object_storage.short"),
		Long:    manager.T("cli.object_storage.long"),
		Aliases: []string{"object", "objects", "objs", "os"},
		GroupID: "products",
	}

	// Adicionar subcomandos
	cmd.AddCommand(buckets.BucketsCommand(parent.Context()))

	parent.AddCommand(cmd)
}
