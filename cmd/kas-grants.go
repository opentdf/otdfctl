package cmd

import (
	"errors"
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var forceFlagValue = false

func policy_assignKasGrant(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	nsID := flagHelper.GetOptionalString("namespace-id")
	attrID := flagHelper.GetOptionalString("attribute-id")
	valID := flagHelper.GetOptionalString("value-id")
	kasID := flagHelper.GetRequiredString("kas-id")

	count := 0
	for _, v := range []string{nsID, attrID, valID} {
		if v != "" {
			count++
		}
	}
	if count != 1 {
		cli.ExitWithError("Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign", errors.New("invalid flag values"))
	}

	var (
		id    string
		res   interface{}
		err   error
		rowID []string
	)

	kas, err := h.GetKasRegistryEntry(kasID)
	if err != nil || kas == nil {
		cli.ExitWithError("Failed to get registered KAS", err)
	}

	ctx := cmd.Context()
	if nsID != "" {
		res, err = h.AssignKasGrantToNamespace(ctx, nsID, kasID)
		if err != nil {
			cli.ExitWithError("Failed to assign KAS Grant for Namespace", err)
		}
		rowID = []string{"Namespace ID", nsID}
	} else if attrID != "" {
		res, err = h.AssignKasGrantToAttribute(ctx, attrID, kasID)
		if err != nil {
			cli.ExitWithError("Failed to assign KAS Grant for Attribute Definition", err)
		}
		rowID = []string{"Attribute ID", attrID}
	} else {
		res, err = h.AssignKasGrantToValue(ctx, valID, kasID)
		if err != nil {
			cli.ExitWithError("Failed to assign KAS Grant for Attribute Value", err)
		}
		rowID = []string{"Value ID", valID}
	}

	t := cli.NewTabular(rowID, []string{"KAS ID", kasID}, []string{"Granted KAS URI", kas.GetUri()})
	HandleSuccess(cmd, id, t, res)
}

func policy_unassignKasGrant(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	nsID := flagHelper.GetOptionalString("namespace-id")
	attrID := flagHelper.GetOptionalString("attribute-id")
	valID := flagHelper.GetOptionalString("value-id")
	kasID := flagHelper.GetRequiredString("kas-id")
	force := flagHelper.GetOptionalBool("force")

	count := 0
	for _, v := range []string{nsID, attrID, valID} {
		if v != "" {
			count++
		}
	}
	if count != 1 {
		cli.ExitWithError("Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign", errors.New("invalid flag values"))
	}
	var (
		res     interface{}
		err     error
		confirm string
		rowID   []string
		rowFQN  []string
	)

	kas, err := h.GetKasRegistryEntry(kasID)
	if err != nil || kas == nil {
		cli.ExitWithError("Failed to get registered KAS", err)
	}
	kasURI := kas.GetUri()

	ctx := cmd.Context()
	if nsID != "" {
		ns, err := h.GetNamespace(nsID)
		if err != nil || ns == nil {
			cli.ExitWithError("Failed to get namespace definition", err)
		}
		confirm = fmt.Sprintf("the grant to namespace FQN (%s) of KAS URI", ns.GetFqn())
		cli.ConfirmAction(cli.ActionDelete, confirm, kasURI, force)
		res, err = h.DeleteKasGrantFromNamespace(ctx, nsID, kasID)
		if err != nil {
			cli.ExitWithError("Failed to update KAS grant for namespace", err)
		}

		rowID = []string{"Namespace ID", nsID}
		rowFQN = []string{"Namespace FQN", ns.GetFqn()}
	} else if attrID != "" {
		attr, err := h.GetAttribute(attrID)
		if err != nil || attr == nil {
			cli.ExitWithError("Failed to get attribute definition", err)
		}
		confirm = fmt.Sprintf("the grant to attribute FQN (%s) of KAS URI", attr.GetFqn())
		cli.ConfirmAction(cli.ActionDelete, confirm, kasURI, force)
		res, err = h.DeleteKasGrantFromAttribute(ctx, attrID, kasID)
		if err != nil {
			cli.ExitWithError("Failed to update KAS grant for attribute", err)
		}

		rowID = []string{"Attribute ID", attrID}
		rowFQN = []string{"Attribute FQN", attr.GetFqn()}
	} else {
		val, err := h.GetAttributeValue(valID)
		if err != nil || val == nil {
			cli.ExitWithError("Failed to get attribute value", err)
		}
		confirm = fmt.Sprintf("the grant to attribute value FQN (%s) of KAS URI", val.GetFqn())
		cli.ConfirmAction(cli.ActionDelete, confirm, kasURI, force)
		_, err = h.DeleteKasGrantFromValue(ctx, valID, kasID)
		if err != nil {
			cli.ExitWithError("Failed to update KAS grant for attribute value", err)
		}
		rowID = []string{"Value ID", valID}
		rowFQN = []string{"Value FQN", val.GetFqn()}
	}

	t := cli.NewTabular(rowID, rowFQN,
		[]string{"KAS ID", kasID},
		[]string{"Unassigned Granted KAS URI", kasURI},
	)
	HandleSuccess(cmd, "", t, res)
}

func init() {
	assignCmd := man.Docs.GetCommand("policy/kas-grants/assign",
		man.WithRun(policy_assignKasGrant),
	)
	assignCmd.Flags().StringP(
		assignCmd.GetDocFlag("namespace-id").Name,
		assignCmd.GetDocFlag("namespace-id").Shorthand,
		assignCmd.GetDocFlag("namespace-id").Default,
		assignCmd.GetDocFlag("namespace-id").Description,
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

	unassignCmd := man.Docs.GetCommand("policy/kas-grants/unassign",
		man.WithRun(policy_unassignKasGrant),
	)
	unassignCmd.Flags().StringP(
		unassignCmd.GetDocFlag("namespace-id").Name,
		unassignCmd.GetDocFlag("namespace-id").Shorthand,
		unassignCmd.GetDocFlag("namespace-id").Default,
		unassignCmd.GetDocFlag("namespace-id").Description,
	)
	unassignCmd.Flags().StringP(
		unassignCmd.GetDocFlag("attribute-id").Name,
		unassignCmd.GetDocFlag("attribute-id").Shorthand,
		unassignCmd.GetDocFlag("attribute-id").Default,
		unassignCmd.GetDocFlag("attribute-id").Description,
	)
	unassignCmd.Flags().StringP(
		unassignCmd.GetDocFlag("value-id").Name,
		unassignCmd.GetDocFlag("value-id").Shorthand,
		unassignCmd.GetDocFlag("value-id").Default,
		unassignCmd.GetDocFlag("value-id").Description,
	)
	unassignCmd.Flags().StringP(
		unassignCmd.GetDocFlag("kas-id").Name,
		unassignCmd.GetDocFlag("kas-id").Shorthand,
		unassignCmd.GetDocFlag("kas-id").Default,
		unassignCmd.GetDocFlag("kas-id").Description,
	)
	unassignCmd.Flags().BoolVar(
		&forceFlagValue,
		unassignCmd.GetDocFlag("force").Name,
		false,
		unassignCmd.GetDocFlag("force").Description,
	)

	cmd := man.Docs.GetCommand("policy/kas-grants",
		man.WithSubcommands(assignCmd, unassignCmd),
	)
	policyCmd.AddCommand(&cmd.Command)
}
