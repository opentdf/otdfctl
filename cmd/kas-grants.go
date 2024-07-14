package cmd

import (
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func policy_updateKasGrant(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)

	attr := flagHelper.GetOptionalString("attribute-id")
	val := flagHelper.GetOptionalString("value-id")
	kas := flagHelper.GetRequiredString("kas-id")
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
			cli.ExitWithError("Failed to update KAS grant for attribute", err)
		}
		id = attr
		header = "Attribute ID"
	} else {
		res, err = h.UpdateKasGrantForValue(val, kas)
		if err != nil {
			cli.ExitWithError("Failed to update KAS grant for attribute value", err)
		}
		id = val
		header = "Value ID"
	}

	t := cli.NewTabular([]string{header, id}, []string{"KAS ID", kas})
	HandleSuccess(cmd, id, t, res)
}

func policy_deleteKasGrant(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	attr := flagHelper.GetOptionalString("attribute-id")
	val := flagHelper.GetOptionalString("value-id")
	kas := flagHelper.GetRequiredString("kas-id")

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
			cli.ExitWithError("Failed to update KAS grant for attribute", err)
		}
		id = attr
		header = "Attribute ID"
	} else {
		_, err := h.DeleteKasGrantFromValue(val, kas)
		if err != nil {
			cli.ExitWithError("Failed to update KAS grant for attribute value", err)
		}
		id = val
		header = "Value ID"
	}

	t := cli.NewTabular([]string{header, id}, []string{"KAS ID", kas})
	HandleSuccess(cmd, id, t, res)
}

func init() {
	updateCmd := man.Docs.GetCommand("policy/kas-grants/update",
		man.WithRun(policy_updateKasGrant),
	)
	updateCmd.Flags().StringP(
		updateCmd.GetDocFlag("attribute-id").Name,
		updateCmd.GetDocFlag("attribute-id").Shorthand,
		updateCmd.GetDocFlag("attribute-id").Default,
		updateCmd.GetDocFlag("attribute-id").Description,
	)
	updateCmd.Flags().StringP(
		updateCmd.GetDocFlag("value-id").Name,
		updateCmd.GetDocFlag("value-id").Shorthand,
		updateCmd.GetDocFlag("value-id").Default,
		updateCmd.GetDocFlag("value-id").Description,
	)
	updateCmd.Flags().StringP(
		updateCmd.GetDocFlag("kas-id").Name,
		updateCmd.GetDocFlag("kas-id").Shorthand,
		updateCmd.GetDocFlag("kas-id").Default,
		updateCmd.GetDocFlag("kas-id").Description,
	)
	injectLabelFlags(&updateCmd.Command, true)

	deleteCmd := man.Docs.GetCommand("policy/kas-grants/delete",
		man.WithRun(policy_deleteKasGrant),
	)
	deleteCmd.Flags().StringP(
		deleteCmd.GetDocFlag("attribute-id").Name,
		deleteCmd.GetDocFlag("attribute-id").Shorthand,
		deleteCmd.GetDocFlag("attribute-id").Default,
		deleteCmd.GetDocFlag("attribute-id").Description,
	)
	deleteCmd.Flags().StringP(
		deleteCmd.GetDocFlag("value-id").Name,
		deleteCmd.GetDocFlag("value-id").Shorthand,
		deleteCmd.GetDocFlag("value-id").Default,
		deleteCmd.GetDocFlag("value-id").Description,
	)
	deleteCmd.Flags().StringP(
		deleteCmd.GetDocFlag("kas-id").Name,
		deleteCmd.GetDocFlag("kas-id").Shorthand,
		deleteCmd.GetDocFlag("kas-id").Default,
		deleteCmd.GetDocFlag("kas-id").Description,
	)

	cmd := man.Docs.GetCommand("policy/kas-grants",
		man.WithSubcommands(updateCmd, deleteCmd),
	)
	policyCmd.AddCommand(&cmd.Command)
}
