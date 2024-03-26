package cmd

import (
	"fmt"
	"strings"

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
		Short: "Update a KAS grant",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)

			attr := flagHelper.GetOptionalString("attribute")
			val := flagHelper.GetOptionalString("value")
			kas := flagHelper.GetRequiredString("kas")
			// uri := flagHelper.GetOptionalString("uri")
			// local := flagHelper.GetOptionalString("public-key-local")
			// remote := flagHelper.GetOptionalString("public-key-remote")
			// labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			if kas == "" || (attr == "" && val == "") {
				cli.ExitWithError("Specify a key access server and an attribute id or attribute value if to update.", nil)
			}
			var (
				id     string
				header string
				// updated interface{}
				updated map[string]interface{}
			)

			// updated.kas_id = kas
			// updated := make(map[string]interface{})
			updated["kas_id"] = kas

			if attr != "" {
				_, err := h.UpdateKasGrantForAttribute(attr, kas)
				// akas, err := h.UpdateKasGrantForAttribute(attr, kas)
				// id = akas.AttributeId
				if err != nil {
					cli.ExitWithError("Could not update KAS grant for attribute", err)
				}
				id = attr
				header = "Attribute ID"
				updated["attribute_id"] = attr
			} else {
				_, err := h.UpdateKasGrantForValue(val, kas)
				// vkas, err := h.UpdateKasGrantForValue(val, kas)
				// id = vkas.ValueId
				if err != nil {
					cli.ExitWithError("Could not update KAS grant for attribute value", err)
				}
				id = val
				header = "Value ID"
				updated["value_id"] = val
				// updated.value_id = val
			}

			// updated, err := h.UpdateKasRegistryEntry(
			// 	id,
			// 	uri,
			// 	pubKey,
			// 	getMetadataMutable(labels),
			// 	getMetadataUpdateBehavior(),
			// )
			// if err != nil {
			// 	cli.ExitWithError("Could not update KAS registry entry", err)
			// }
			t := cli.NewTabular().
				Rows([][]string{
					{header, id},
					{"KAS ID", kas},
					// TODO: render labels [https://github.com/opentdf/tructl/issues/73]
				}...)
			HandleSuccess(cmd, id, t, updated)
		},
	}

	kasGrantsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a KAS grant",
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
