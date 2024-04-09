package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/otdfctl/issues/73] is addressed

var (
	policy_namespacesCommands = []string{
		policy_namespacesCreateCmd.Use,
		policy_namespaceGetCmd.Use,
		policy_namespacesListCmd.Use,
		policy_namespaceUpdateCmd.Use,
		policy_namespaceDeactivateCmd.Use,
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
				errMsg := fmt.Sprintf("Failed to get namespace (%s)", id)
				cli.ExitWithError(errMsg, err)
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

			state := cli.GetState(cmd)
			list, err := h.ListNamespaces(state)
			if err != nil {
				cli.ExitWithError("Failed to list namespaces", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Name", "Active")
			for _, ns := range list {
				t.Row(
					ns.Id,
					ns.Name,
					strconv.FormatBool(ns.Active.GetValue()),
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
				cli.ExitWithError("Failed to create namespace", err)
			}

			t := cli.NewTabular().Rows([][]string{
				{"Name", name},
				{"Id", created.Id},
			}...)
			HandleSuccess(cmd, created.Id, t, created)
		},
	}

	policy_namespaceDeactivateCmd = &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate a namespace by id",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			ns, err := h.GetNamespace(id)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmAction(cli.ActionDeactivate, "namespace", ns.Name)

			d, err := h.DeactivateNamespace(id)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to deactivate namespace (%s)", id)
				cli.ExitWithError(errMsg, err)
			}

			t := cli.NewTabular().
				Rows([][]string{
					{"Id", ns.Id},
					{"Name", ns.Name},
				}...)
			HandleSuccess(cmd, ns.Id, t, d)
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
				cli.ExitWithError(fmt.Sprintf("Failed to update namespace (%s)", id), err)
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
	policy_namespacesListCmd.Flags().StringP("state", "s", "active", "Filter by state [active, inactive, any]")

	policy_namespacesCmd.AddCommand(policy_namespacesCreateCmd)
	policy_namespacesCreateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")
	injectLabelFlags(policy_namespacesCreateCmd, false)

	policy_namespacesCmd.AddCommand(policy_namespaceUpdateCmd)
	policy_namespaceUpdateCmd.Flags().StringP("id", "i", "", "Id of the namespace")
	injectLabelFlags(policy_namespaceUpdateCmd, true)

	policy_namespacesCmd.AddCommand(policy_namespaceDeactivateCmd)
	policy_namespaceDeactivateCmd.Flags().StringP("id", "i", "", "Id of the namespace")
}
