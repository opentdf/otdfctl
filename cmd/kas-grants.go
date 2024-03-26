package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/platform/protocol/go/kasregistry"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	kasGrants_crudCommands = []string{
		kasGrantsUpdateCmd.Use,
		kasGrantsDeleteCmd.Use,
	}

	// KasGrantsCmd is the command for managing KAS grants
	kasGrantsCmd = &cobra.Command{
		Use:   "kas-grants",
		Short: "Manage Key Access Server grants [" + strings.Join(kasGrants_crudCommands, ", ") + "]",
	}

	// Update one KAS registry entry
	kasGrantsUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a KAS registry entry",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)

			id := flagHelper.GetRequiredString("id")
			uri := flagHelper.GetOptionalString("uri")
			local := flagHelper.GetOptionalString("public-key-local")
			remote := flagHelper.GetOptionalString("public-key-remote")
			labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			if local == "" && remote == "" && len(labels) == 0 && uri == "" {
				cli.ExitWithError("No values were passed to update. Please pass at least one value to update (E.G. 'uri', 'public-key-local', 'public-key-remote', 'label')", nil)
			}

			// TODO: should update of a type of key be a dangerous mutation or cause a need for confirmation in the CLI?
			var pubKey *kasregistry.PublicKey
			if local != "" && remote != "" {
				e := fmt.Errorf("Only one public key is allowed. Please pass either a local or remote public key but not both")
				cli.ExitWithError("Issue with update flags 'public-key-local' and 'public-key-remote': ", e)
			} else if local != "" {
				pubKey = &kasregistry.PublicKey{PublicKey: &kasregistry.PublicKey_Local{Local: local}}
			} else if remote != "" {
				pubKey = &kasregistry.PublicKey{PublicKey: &kasregistry.PublicKey_Remote{Remote: remote}}
			}

			updated, err := h.UpdateKasRegistryEntry(
				id,
				uri,
				pubKey,
				getMetadataMutable(labels),
				getMetadataUpdateBehavior(),
			)
			if err != nil {
				cli.ExitWithError("Could not update KAS registry entry", err)
			}
			t := cli.NewTabular().
				Rows([][]string{
					{"Id", id},
					{"URI", uri},
					// TODO: render labels [https://github.com/opentdf/tructl/issues/73]
				}...)
			HandleSuccess(cmd, id, t, updated)
		},
	}

	kasGrantsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a KAS registry entry by id",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			kas, err := h.GetKasRegistryEntry(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find KAS registry entry (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
			}

			cli.ConfirmDelete("KAS Registry Entry: ", id)

			if err := h.DeleteKasRegistryEntry(id); err != nil {
				errMsg := fmt.Sprintf("Could not delete KAS registry entry (%s)", id)
				cli.ExitWithError(errMsg, err)
			}

			t := cli.NewTabular().
				Rows([][]string{
					{"Id", kas.Id},
					{"URI", kas.Uri},
				}...)

			HandleSuccess(cmd, kas.Id, t, kas)
		},
	}
)

func init() {
	policyCmd.AddCommand(kasGrantsCmd)

	kasGrantsCmd.AddCommand(kasGrantsUpdateCmd)
	kasGrantsUpdateCmd.Flags().StringP("attribute", "a", "", "attribute id")
	kasGrantsUpdateCmd.Flags().StringP("value", "v", "", "attribute value id")
	kasGrantsUpdateCmd.Flags().StringP("kas", "k", "", "kas id")
	injectLabelFlags(kasGrantsUpdateCmd, true)

	kasGrantsCmd.AddCommand(kasGrantsDeleteCmd)
	kasGrantsDeleteCmd.Flags().StringP("attribute", "a", "", "attribute id")
	kasGrantsDeleteCmd.Flags().StringP("value", "v", "", "attribute value id")
	kasGrantsDeleteCmd.Flags().StringP("kas", "k", "", "kas id")
}

func init() {
	rootCmd.AddCommand(kasGrantsCmd)
}
