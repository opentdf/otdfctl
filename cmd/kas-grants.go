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
	assignCmd := man.Docs.GetCommand("policy/kas-grants/assign",
		man.WithRun(policy_updateKasGrant),
	)
	assignCmd.Flags().StringP(
		assignCmd.GetDocFlag("attribute-id").Name,
		assignCmd.GetDocFlag("attribute-id").Shorthand,
		assignCmd.GetDocFlag("attribute-id").Default,
		assignCmd.GetDocFlag("attribute-id").Description,
	)
	assignCmd.Flags().StringP(
		assignCmd.GetDocFlag("value-id").Name,
		assignCmd.GetDocFlag("value-id").Shorthand,
		assignCmd.GetDocFlag("value-id").Default,
		assignCmd.GetDocFlag("value-id").Description,
	)
	assignCmd.Flags().StringP(
		assignCmd.GetDocFlag("kas-id").Name,
		assignCmd.GetDocFlag("kas-id").Shorthand,
		assignCmd.GetDocFlag("kas-id").Default,
		assignCmd.GetDocFlag("kas-id").Description,
	)
	injectLabelFlags(&assignCmd.Command, true)

	removeCmd := man.Docs.GetCommand("policy/kas-grants/remove",
		man.WithRun(policy_deleteKasGrant),
	)
	removeCmd.Flags().StringP(
		removeCmd.GetDocFlag("attribute-id").Name,
		removeCmd.GetDocFlag("attribute-id").Shorthand,
		removeCmd.GetDocFlag("attribute-id").Default,
		removeCmd.GetDocFlag("attribute-id").Description,
	)
	removeCmd.Flags().StringP(
		removeCmd.GetDocFlag("value-id").Name,
		removeCmd.GetDocFlag("value-id").Shorthand,
		removeCmd.GetDocFlag("value-id").Default,
		removeCmd.GetDocFlag("value-id").Description,
	)
	removeCmd.Flags().StringP(
		removeCmd.GetDocFlag("kas-id").Name,
		removeCmd.GetDocFlag("kas-id").Shorthand,
		removeCmd.GetDocFlag("kas-id").Default,
		removeCmd.GetDocFlag("kas-id").Description,
	)

	cmd := man.Docs.GetCommand("policy/kas-grants",
		man.WithSubcommands(assignCmd, removeCmd),
	)
	policyCmd.AddCommand(&cmd.Command)
}
