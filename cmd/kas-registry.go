package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/opentdf-v2-poc/sdk/common"
	"github.com/opentdf/opentdf-v2-poc/sdk/kasregistry"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	kasRegistry_crudCommands = []string{
		kasRegistrysCreateCmd.Use,
		kasRegistryGetCmd.Use,
		kasRegistrysListCmd.Use,
		kasRegistryUpdateCmd.Use,
		kasRegistryDeleteCmd.Use,
	}

	// KasRegistryCmd is the command for managing KAS registrations
	kasRegistryCmd = &cobra.Command{
		Use:   "kas-registry",
		Short: "Manage KAS Registry entries [ [" + strings.Join(kasRegistry_crudCommands, ", ") + "]",
		Long: `
	Manage KAS Registry entries within the platform.
	
	The Key Access Server (KAS) registry is a record of servers granting and maintaining public keys. The registry contains critical information like each server's uri, its public key (which can be either local or at a remote uri), and any metadata about the server.
	Key Access Servers grant keys for specified attributes and attribute values via Attribute Key Access Grants and Attribute Value Key Access Grants, which are managed separately from the registry.
	`,
	}

	kasRegistryGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a KAS registry entry by id",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			kasRegEntry, err := h.GetKasRegistryEntry(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find KAS registry entry (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			fmt.Println(cli.SuccessMessage("KAS registry entry found"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", kasRegEntry.Id},
						{"Metadata.Labels.Name", kasRegEntry.Metadata.Labels["name"]},
						{"Metadata.Description", kasRegEntry.Metadata.Description},
						{"URI", kasRegEntry.Uri},
					}...).Render(),
			)
		},
	}

	kasRegistrysListCmd = &cobra.Command{
		Use:   "list",
		Short: "List KAS registry entries",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			list, err := h.ListKasRegistryEntries()
			if err != nil {
				cli.ExitWithError("Could not get KAS registry entries", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "URI", "PublicKey", "MetaData")
			for _, kre := range list {
				t.Row(
					kre.Id,
					kre.Metadata.Labels["name"],
					kre.Metadata.Description,
					kre.Uri,
				)
			}
			fmt.Println(t.Render())
		},
	}

	kasRegistrysCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new KAS registry entry, i.e. 'https://example.com'",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			uri := flagHelper.GetRequiredString("uri")
			public_key := flagHelper.GetRequiredString("public_key")
			name := flagHelper.GetRequiredString("name")
			description := flagHelper.GetRequiredString("description")

			labelsMap := make(map[string]string)

			// check if a name has been passed
			if name != "" {
				labelsMap["name"] = name
			}

			// check if a description has been passed
			if description == "" {
				description = "No description provided"
			}

			createdKasRegEntry, err := h.CreateKasRegistryEntry(uri, &kasregistry.PublicKey{
				PublicKey: &kasregistry.PublicKey_Local{
					Local: public_key,
				},
			}, &common.MetadataMutable{
				Labels:      labelsMap,
				Description: description,
			})
			if err != nil {
				cli.ExitWithError("Could not create KAS registry entry", err)
			}

			fmt.Println(cli.SuccessMessage("KAS registry entry found"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", createdKasRegEntry.Id},
						{"Metadata.Labels.Name", createdKasRegEntry.Metadata.Labels["name"]},
						{"Metadata.Description", createdKasRegEntry.Metadata.Description},
						{"URI", createdKasRegEntry.Uri},
					}...).Render(),
			)
		},
	}

	// Update one KAS registry entry
	kasRegistryUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a KAS registry entry",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)

			id := flagHelper.GetRequiredString("id")
			uri := flagHelper.GetOptionalString("uri")
			public_key := flagHelper.GetOptionalString("public_key")
			name := flagHelper.GetOptionalString("name")
			description := flagHelper.GetOptionalString("description")

			// Initialize KeyAccessServer only if needed
			kas := &kasregistry.KeyAccessServerCreateUpdate{}

			// Check each optional parameter and set it if provided
			if uri != "" {
				kas.Uri = uri // Assuming Uri is now a *string
			}
			if public_key != "" {
				var publicKeyObj = kasregistry.PublicKey{
					PublicKey: &kasregistry.PublicKey_Local{
						Local: public_key,
					},
				}

				kas.PublicKey = &publicKeyObj
			}

			// create our metadata object
			var metaDataObj = common.MetadataMutable{}
			// we need to individually check if these parameters were passed, and add them to the update object
			if name != "" {
				metaDataObj.Labels = map[string]string{"name": name}
			}

			if description != "" {
				metaDataObj.Description = description
			}

			if description != "" || name != "" {
				kas.Metadata = &metaDataObj
			}

			req := &kasregistry.UpdateKeyAccessServerRequest{
				Id: id,
			}
			// set the kas update object on the request
			req.KeyAccessServer = kas

			// now lets make sure it is valid, did anything get passed as updated values?
			if kas.Uri != "" && kas.PublicKey != nil && kas.Metadata != nil {
				cli.ExitWithError("No values were passed to update. Please pass at least one value to update (E.G. 'uri', 'name', 'description', 'publicKey')", nil)
			}

			if _, err := h.UpdateKasRegistryEntry(
				id,
				req,
			); err != nil {
				cli.ExitWithError("Could not update KAS registry entry", err)
			}
			fmt.Println(cli.SuccessMessage(fmt.Sprintf("Namespace id: (%s) updated. Name set to (%s).", id, name)))
		},
	}

	kasRegistryDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a KAS registry entry by id",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			kasRegEntry, err := h.GetKasRegistryEntry(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find KAS registry entry (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmDelete("KAS Registry Entry", kasRegEntry.Metadata.Labels["name"])

			if err := h.DeleteKasRegistryEntry(id); err != nil {
				errMsg := fmt.Sprintf("Could not delete KAS registry entry (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			fmt.Println(cli.SuccessMessage("KAS Registry Entry deleted"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", kasRegEntry.Id},
						{"Name", kasRegEntry.Metadata.Labels["name"]},
					}...).Render(),
			)
		},
	}
)

func init() {
	policyCmd.AddCommand(kasRegistryCmd)

	kasRegistryCmd.AddCommand(kasRegistryGetCmd)
	kasRegistryGetCmd.Flags().StringP("id", "i", "", "Id of the KAS registry entry")

	kasRegistryCmd.AddCommand(kasRegistrysListCmd)

	kasRegistryCmd.AddCommand(kasRegistrysCreateCmd)
	kasRegistrysCreateCmd.Flags().StringP("description", "d", "", "The common description of the KAS registry entry")
	kasRegistrysCreateCmd.Flags().StringP("uri", "u", "", "The URI of the KAS registry entry")
	kasRegistrysCreateCmd.Flags().StringP("name", "n", "", "Name value of the KAS registry entry")
	kasRegistrysCreateCmd.Flags().StringP("public_key", "p", "", "The KAS Public Key")

	kasRegistryCmd.AddCommand(kasRegistryUpdateCmd)
	kasRegistryUpdateCmd.Flags().StringP("id", "i", "", "Id of the KAS registry entry")
	kasRegistryUpdateCmd.Flags().StringP("description", "d", "", "The common description of the KAS registry entry")
	kasRegistryUpdateCmd.Flags().StringP("uri", "u", "", "The URI of the KAS registry entry")
	kasRegistryUpdateCmd.Flags().StringP("name", "n", "", "Name value of the KAS registry entry")
	kasRegistryUpdateCmd.Flags().StringP("public_key", "p", "", "The KAS Public Key")

	kasRegistryCmd.AddCommand(kasRegistryDeleteCmd)
	kasRegistryDeleteCmd.Flags().StringP("id", "i", "", "Id of the KAS registry entry")
}

func init() {
	rootCmd.AddCommand(kasRegistryCmd)
}
