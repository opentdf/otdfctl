package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/platform/protocol/go/kasregistry"
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
		Short: "Manage Key Access Server registrations [" + strings.Join(kasRegistry_crudCommands, ", ") + "]",
		Long: `
	Manage Key Access Server registrations within the platform.
	
	The Key Access Server (KAS) registry is a record of servers granting and maintaining public keys. The registry contains critical
	information like each server's uri, its public key (which can be either local or at a remote uri), and any metadata about the server.
	Key Access Servers grant keys for specified Attributes and their Values via Attribute Key Access Grants and Attribute Value
	Key Access Grants.
	`,
	}

	kasRegistryGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a registered Key Access Server by id",
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

			keyType := "Local"
			key := kas.PublicKey.GetLocal()
			if kas.PublicKey.GetRemote() != "" {
				keyType = "Remote"
				key = kas.PublicKey.GetRemote()
			}

			t := cli.NewTabular().
				Rows([][]string{
					{"Id", kas.Id},
					// TODO: render labels [https://github.com/opentdf/tructl/issues/73]
					{"URI", kas.Uri},
					{"PublicKey Type", keyType},
					{"PublicKey", key},
				}...)
			HandleSuccess(cmd, kas.Id, t, kas)
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
			t.Headers("Id", "URI", "PublicKey Location", "PublicKey")
			for _, kas := range list {
				keyType := "Local"
				key := kas.PublicKey.GetLocal()
				if kas.PublicKey.GetRemote() != "" {
					keyType = "Remote"
					key = kas.PublicKey.GetRemote()
				}

				t.Row(
					kas.Id,
					kas.Uri,
					keyType,
					key,
					// TODO: render labels [https://github.com/opentdf/tructl/issues/73]
				)
			}
			HandleSuccess(cmd, "", t, list)
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
			local := flagHelper.GetOptionalString("public-key-local")
			remote := flagHelper.GetOptionalString("public-key-remote")
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			if local == "" && remote == "" {
				e := fmt.Errorf("A public key is required. Please pass either a local or remote public key")
				cli.ExitWithError("Issue with create flags 'public-key-local' and 'public-key-remote': ", e)
			}

			key := &kasregistry.PublicKey{}
			keyType := "Local"
			if local != "" {
				if remote != "" {
					e := fmt.Errorf("Only one public key is allowed. Please pass either a local or remote public key but not both")
					cli.ExitWithError("Issue with create flags 'public-key-local' and 'public-key-remote': ", e)
				}
				key.PublicKey = &kasregistry.PublicKey_Local{Local: local}
			} else {
				keyType = "Remote"
				key.PublicKey = &kasregistry.PublicKey_Remote{Remote: remote}
			}

			created, err := h.CreateKasRegistryEntry(
				uri,
				key,
				getMetadataMutable(metadataLabels),
			)
			if err != nil {
				cli.ExitWithError("Could not create KAS registry entry", err)
			}

			t := cli.NewTabular().
				Rows([][]string{
					{"Id", created.Id},
					{"URI", created.Uri},
					{"PublicKey Type", keyType},
					{"PublicKey", local},
					// TODO: render labels [https://github.com/opentdf/tructl/issues/73]
				}...)

			HandleSuccess(cmd, created.Id, t, created)
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

	kasRegistryDeleteCmd = &cobra.Command{
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
	policyCmd.AddCommand(kasRegistryCmd)

	kasRegistryCmd.AddCommand(kasRegistryGetCmd)
	kasRegistryGetCmd.Flags().StringP("id", "i", "", "Id of the KAS registry entry")

	kasRegistryCmd.AddCommand(kasRegistrysListCmd)
	// TODO: active, inactive, any state querying [https://github.com/opentdf/tructl/issues/68]

	kasRegistryCmd.AddCommand(kasRegistrysCreateCmd)
	kasRegistrysCreateCmd.Flags().StringP("uri", "u", "", "The URI of the KAS registry entry")
	kasRegistrysCreateCmd.Flags().StringP("public-key-local", "p", "", "A local public key for the registered Key Access Server (KAS)")
	kasRegistrysCreateCmd.Flags().StringP("public-key-remote", "r", "", "A remote endpoint that provides a public key for the registered Key Access Server (KAS)")
	kasRegistrysCreateCmd.Flags().StringSliceVarP(&metadataLabels, "label", "l", []string{}, "Optional metadata 'labels' in the format: key=value")

	kasRegistryCmd.AddCommand(kasRegistryUpdateCmd)
	kasRegistryUpdateCmd.Flags().StringP("id", "i", "", "Id of the KAS registry entry")
	kasRegistryUpdateCmd.Flags().StringP("uri", "u", "", "The URI of the KAS registry entry")
	kasRegistryUpdateCmd.Flags().StringP("public-key-local", "p", "", "A local public key for the registered Key Access Server (KAS)")
	kasRegistryUpdateCmd.Flags().StringP("public-key-remote", "r", "", "A remote endpoint that serves a public key for the registered Key Access Server (KAS)")
	kasRegistryUpdateCmd.Flags().StringSliceVarP(&metadataLabels, "label", "l", []string{}, "Optional metadata 'labels' in the format: key=value")
	kasRegistryUpdateCmd.Flags().BoolVar(&forceReplaceMetadataLabels, "force-replace-labels", false, "Destructively replace entire set of existing metadata 'labels' with any provided to this command.")

	kasRegistryCmd.AddCommand(kasRegistryDeleteCmd)
	kasRegistryDeleteCmd.Flags().StringP("id", "i", "", "Id of the KAS registry entry")
}

func init() {
	rootCmd.AddCommand(kasRegistryCmd)
}
