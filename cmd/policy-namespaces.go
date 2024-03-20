package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

var (
	policy_namespacesCommands = []string{
		policy_namespacesCreateCmd.Use,
		policy_namespaceGetCmd.Use,
		policy_namespacesListCmd.Use,
		policy_namespaceUpdateCmd.Use,
		policy_namespaceDeleteCmd.Use,
	}

	policy_namespacesCmd = &cobra.Command{
		Use:   "namespaces",
		Short: "Manage namespaces [" + strings.Join(policy_namespacesCommands, ", ") + "]",
		Long: `
Namespaces - commands to manage attribute namespaces within the platform.
		
Namespaces drive associations of attributes and their values and differentiate between them.
For example: "bob.com" and "alice.net" are different namespaces that may have the same
or different attributes tied to each.
`,
	}

	policy_namespaceGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a namespace by id",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			ns, err := h.GetNamespace(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find namespace (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
			}

			t := cli.NewTabular().
				Rows([][]string{
					{"Id", ns.Id},
					{"Name", ns.Name},
				}...)
			HandleSuccess(cmd, ns.Id, t, ns)
		},
	}

	policy_namespacesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List namespaces",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			list, err := h.ListNamespaces()
			if err != nil {
				cli.ExitWithError("Could not get namespaces", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Name")
			for _, ns := range list {
				t.Row(
					ns.Id,
					ns.Name,
				)
			}
			HandleSuccess(cmd, "", t, list)
		},
	}

	policy_namespacesCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace, i.e. 'https://example.com'",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			name := flagHelper.GetRequiredString("name")
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			created, err := h.CreateNamespace(name, getMetadataMutable(metadataLabels))
			if err != nil {
				cli.ExitWithError("Could not create namespace", err)
			}

			t := cli.NewTabular().Rows([][]string{
				{"Name", name},
				{"Id", created.Id},
			}...)
			HandleSuccess(cmd, created.Id, t, created)
		},
	}

	policy_namespaceDeleteCmd = &cobra.Command{
		Use:   "deactivate",
		Short: "Delete a namespace by id",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			ns, err := h.GetNamespace(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find namespace (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
			}

			cli.ConfirmDelete("namespace", ns.Name)

			if err := h.DeactivateNamespace(id); err != nil {
				errMsg := fmt.Sprintf("Could not deactivate namespace (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
			}

			t := cli.NewTabular().
				Rows([][]string{
					{"Id", ns.Id},
					{"Name", ns.Name},
				}...)
			HandleSuccess(cmd, ns.Id, t, ns)
		},
	}

	// Update one namespace
	policy_namespaceUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a namespace",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			ns, err := h.UpdateNamespace(
				id,
				getMetadataMutable(labels),
				getMetadataUpdateBehavior(),
			)
			if err != nil {
				cli.ExitWithError("Could not update namespace", err)
			}

			t := cli.NewTabular().Rows([][]string{
				{"Id", ns.Id},
				{"Name", ns.Name},
			}...)
			HandleSuccess(cmd, id, t, ns)
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_namespacesCmd)

	policy_namespacesCmd.AddCommand(policy_namespaceGetCmd)
	policy_namespaceGetCmd.Flags().StringP("id", "i", "", "Id of the namespace")

	policy_namespacesCmd.AddCommand(policy_namespacesListCmd)

	policy_namespacesCmd.AddCommand(policy_namespacesCreateCmd)
	policy_namespacesCreateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")
	policy_namespacesCreateCmd.Flags().StringSliceVarP(&metadataLabels, "label", "l", []string{}, "Optional metadata 'labels' in the format: key=value")

	policy_namespacesCmd.AddCommand(policy_namespaceUpdateCmd)
	policy_namespaceUpdateCmd.Flags().StringP("id", "i", "", "Id of the namespace")
	policy_namespaceUpdateCmd.Flags().StringSliceVarP(&metadataLabels, "label", "l", []string{}, "Optional new metadata 'labels' in the format: key=value")
	policy_namespaceUpdateCmd.Flags().BoolVar(&forceReplaceMetadataLabels, "force-replace-labels", false, "Destructively replace entire set of existing metadata 'labels' with any provided to this command.")

	policy_namespacesCmd.AddCommand(policy_namespaceDeleteCmd)
	policy_namespaceDeleteCmd.Flags().StringP("id", "i", "", "Id of the namespace")
}
