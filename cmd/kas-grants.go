package cmd

import (
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

			if attr == "" && val == "" {
				cli.ExitWithError("Must specify and Attribute Definition id or Value id to update.", nil)
			}
			var (
				id     string
				header string
				res    interface{}
				err    error
			)

			if attr != "" {
				res, err = h.UpdateKasGrantForAttribute(attr, kas)
				if err != nil {
					cli.ExitWithError("Could not update KAS grant for attribute", err)
				}
				id = attr
				header = "Attribute ID"
			} else {
				res, err = h.UpdateKasGrantForValue(val, kas)
				if err != nil {
					cli.ExitWithError("Could not update KAS grant for attribute value", err)
				}
				id = val
				header = "Value ID"
			}

			t := cli.NewTabular().
				Rows([][]string{
					{header, id},
					{"KAS ID", kas},
				}...)
			HandleSuccess(cmd, id, t, res)
		},
	}

	kasGrantsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a KAS grant",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			attr := flagHelper.GetOptionalString("attribute")
			val := flagHelper.GetOptionalString("value")
			kas := flagHelper.GetRequiredString("kas")

			if attr == "" && val == "" {
				cli.ExitWithError("Must specify and Attribute Definition id or Value id to delete.", nil)
			}
			var (
				id     string
				header string
				res    interface{}
				err    error
			)

			cli.ConfirmAction(cli.ActionDelete, "KAS ID: ", kas)

			if attr != "" {
				res, err = h.DeleteKasGrantFromAttribute(attr, kas)
				if err != nil {
					cli.ExitWithError("Could not update KAS grant for attribute", err)
				}
				id = attr
				header = "Attribute ID"
			} else {
				_, err := h.DeleteKasGrantFromValue(val, kas)
				if err != nil {
					cli.ExitWithError("Could not update KAS grant for attribute value", err)
				}
				id = val
				header = "Value ID"
			}

			t := cli.NewTabular().
				Rows([][]string{
					{header, id},
					{"KAS ID", kas},
				}...)
			HandleSuccess(cmd, id, t, res)
		},
	}
)

func init() {
	policyCmd.AddCommand(kasGrantsCmd)

	kasGrantsCmd.AddCommand(kasGrantsUpdateCmd)
	kasGrantsUpdateCmd.Flags().StringP("attribute", "a", "", "Attribute Definition ID")
	kasGrantsUpdateCmd.Flags().StringP("value", "v", "", "Attribute Value ID")
	kasGrantsUpdateCmd.Flags().StringP("kas", "k", "", "Key Access Server (KAS) ID")
	injectLabelFlags(kasGrantsUpdateCmd, true)

	kasGrantsCmd.AddCommand(kasGrantsDeleteCmd)
	kasGrantsDeleteCmd.Flags().StringP("attribute", "a", "", "Attribute Definition ID")
	kasGrantsDeleteCmd.Flags().StringP("value", "v", "", "Attribute Value ID")
	kasGrantsDeleteCmd.Flags().StringP("kas", "k", "", "Key Access Server (KAS) ID")
	rootCmd.AddCommand(kasGrantsCmd)
}
