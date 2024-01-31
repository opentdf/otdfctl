package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	namespacesCommands = []string{
		namespacesCreateCmd.Use,
		namespaceGetCmd.Use,
		namespacesListCmd.Use,
		namespaceUpdateCmd.Use,
		namespaceDeleteCmd.Use,
	}

	namespacesCmd = &cobra.Command{
		Use:   "namespaces",
		Short: "Manage namespaces [" + strings.Join(namespacesCommands, ", ") + "]",
		Long: `
Namespaces - commands to manage attribute namespaces within the platform.
		
Namespaces drive associations of attributes and their values and differentiate between them.
For example: "bob.com" and "alice.net" are different namespaces that may have the same
or different attributes tied to each.
`,
	}

	namespaceGetCmd = &cobra.Command{
		Use:   "get <id>",
		Short: "Get a namespace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			id, err := strconv.Atoi(args[0])
			if err != nil {
				cli.ExitWithError("Invalid ID", err)
			}

			ns, err := h.GetNamespace(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find namespace (%d)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			fmt.Println(cli.SuccessMessage("Namespace found"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", strconv.Itoa(int(ns.Id))},
						{"Name", ns.Name},
					}...).Render(),
			)
		},
	}

	namespacesListCmd = &cobra.Command{
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
					strconv.Itoa(int(ns.Id)),
					ns.Name,
				)
			}
			fmt.Println(t.Render())
		},
	}

	namespacesCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
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

	namespaceDeleteCmd = &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a namespace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(cli.ErrorMessage("Invalid ID", err))
				os.Exit(1)
			}
			ns, err := h.GetNamespace(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find namespace (%d)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmDelete("namespace", ns.Name)

			if err := h.DeleteNamespace(id); err != nil {
				errMsg := fmt.Sprintf("Could not delete namespace (%d)", id)
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
	namespaceUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a namespace",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)

			id := flagHelper.GetRequiredInt32("id")
			name := flagHelper.GetRequiredString("name")

			if _, err := h.UpdateNamespace(
				id,
				name,
			); err != nil {
				cli.ExitWithError("Could not update namespace", err)
				return
			} else {
				fmt.Println(cli.SuccessMessage(fmt.Sprintf("Namespace id: (%d) updated. Name set to (%s).", id, name)))
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(namespacesCmd)

	namespacesCmd.AddCommand(namespaceGetCmd)

	namespacesCmd.AddCommand(namespacesListCmd)

	namespacesCmd.AddCommand(namespacesCreateCmd)
	namespacesCreateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")

	namespacesCmd.AddCommand(namespaceUpdateCmd)
	namespaceUpdateCmd.Flags().Int32P("id", "i", 0, "Id of the namespace")
	namespaceUpdateCmd.Flags().StringP("name", "n", "", "Name value of the namespace")

	namespacesCmd.AddCommand(namespaceDeleteCmd)
}
