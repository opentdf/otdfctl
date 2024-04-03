package cmd

import (
	"fmt"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/man"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

var (
	policy_namespacesCommands = []string{
		policy_namespacesCreateCmd.Use,
		policy_namespaceGetCmd.Use,
		policy_namespacesListCmd.Use,
		policy_namespaceUpdateCmd.Use,
		policy_namespaceDeactivateCmd.Use,
	}

	policy_namespacesCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy/attributes/namespaces").Use,
		Short: man.Docs.GetDoc("policy/attributes/namespaces").GetShort(policy_namespacesCommands),
		Long:  man.Docs.GetDoc("policy/attributes/namespaces").Long,
	}

	policy_namespaceGetCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy/attributes/namespaces/get").Use,
		Short: man.Docs.GetDoc("policy/attributes/namespaces/get").Short,
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
		Use:   man.Docs.GetDoc("policy/attributes/namespaces/list").Use,
		Short: man.Docs.GetDoc("policy/attributes/namespaces/list").Short,
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			list, err := h.ListNamespaces()
			if err != nil {
				cli.ExitWithError("Failed to list namespaces", err)
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
		Use:   man.Docs.GetDoc("policy/attributes/namespaces/create").Use,
		Short: man.Docs.GetDoc("policy/attributes/namespaces/create").Short,
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
		Use:   man.Docs.GetDoc("policy/attributes/namespaces/deactivate").Use,
		Short: man.Docs.GetDoc("policy/attributes/namespaces/deactivate").Short,
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
		Use:   man.Docs.GetDoc("policy/attributes/namespaces/update").Use,
		Short: man.Docs.GetDoc("policy/attributes/namespaces/update").Short,
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

	policy_namespacesCmd.AddCommand(policy_namespacesCreateCmd)
	policy_namespacesCreateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")
	injectLabelFlags(policy_namespacesCreateCmd, false)

	policy_namespacesCmd.AddCommand(policy_namespaceUpdateCmd)
	policy_namespaceUpdateCmd.Flags().StringP("id", "i", "", "Id of the namespace")
	injectLabelFlags(policy_namespaceUpdateCmd, true)

	policy_namespacesCmd.AddCommand(policy_namespaceDeactivateCmd)
	policy_namespaceDeactivateCmd.Flags().StringP("id", "i", "", "Id of the namespace")
}
