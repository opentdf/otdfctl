package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

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
				cli.ExitWithError(errMsg, err)
			}

			fmt.Println(cli.SuccessMessage("Namespace found"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", ns.Id},
						{"Name", ns.Name},
					}...).Render(),
			)
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
			fmt.Println(t.Render())
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

			created, err := h.CreateNamespace(name)
			if err != nil {
				cli.ExitWithError("Could not create namespace", err)
			}

			fmt.Println(cli.SuccessMessage("Namespace created"))
			fmt.Println(
				cli.NewTabular().Rows([][]string{
					{"Name", name},
					{"Id", created.Id},
				}...).Render(),
			)
		},
	}

	policy_namespaceDeleteCmd = &cobra.Command{
		Use:   "delete",
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
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmDelete("namespace", ns.Name)

			if err := h.DeleteNamespace(id); err != nil {
				errMsg := fmt.Sprintf("Could not delete namespace (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			fmt.Println(cli.SuccessMessage("Namespace deleted"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", ns.Id},
						{"Name", ns.Name},
					}...).Render(),
			)
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
			name := flagHelper.GetRequiredString("name")

			if _, err := h.UpdateNamespace(
				id,
				name,
			); err != nil {
				cli.ExitWithError("Could not update namespace", err)
			}
			fmt.Println(cli.SuccessMessage(fmt.Sprintf("Namespace id: (%s) updated. Name set to (%s).", id, name)))
		},
	}
)

func init() {
	policyCmd.AddCommand(namespacesCmd)

	policy_namespacesCmd.AddCommand(namespaceGetCmd)
	policy_namespaceGetCmd.Flags().StringP("id", "i", "", "Id of the namespace")

	policy_namespacesCmd.AddCommand(namespacesListCmd)

	policy_namespacesCmd.AddCommand(namespacesCreateCmd)
	policy_namespacesCreateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")

	policy_namespacesCmd.AddCommand(namespaceUpdateCmd)
	policy_namespaceUpdateCmd.Flags().StringP("id", "i", "", "Id of the namespace")
	policy_namespaceUpdateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")

	policy_namespacesCmd.AddCommand(namespaceDeleteCmd)
	policy_namespaceDeleteCmd.Flags().StringP("id", "i", "", "Id of the namespace")
}
