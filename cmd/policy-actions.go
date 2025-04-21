package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
	"github.com/spf13/cobra"
)

func policy_getAction(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	mapping, err := h.GetSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find subject mapping (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	var actionsJSON []byte
	if actionsJSON, err = json.Marshal(mapping.GetActions()); err != nil {
		cli.ExitWithError("Error marshalling subject mapping actions", err)
	}

	var subjectSetsJSON []byte
	if subjectSetsJSON, err = json.Marshal(mapping.GetSubjectConditionSet().GetSubjectSets()); err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", mapping.GetId()},
		{"Attribute Value: Id", mapping.GetAttributeValue().GetId()},
		{"Attribute Value: Value", mapping.GetAttributeValue().GetValue()},
		{"Actions", string(actionsJSON)},
		{"Subject Condition Set: Id", mapping.GetSubjectConditionSet().GetId()},
		{"Subject Condition Set", string(subjectSetsJSON)},
	}
	if mdRows := getMetadataRows(mapping.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, mapping.GetId(), t, mapping)
}

func policy_listActions(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	list, page, err := h.ListSubjectMappings(limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to get subject mappings", err)
	}
	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("value_id", "Attribute Value Id", cli.FlexColumnWidthFour),
		table.NewFlexColumn("value_fqn", "Attibribute Value FQN", cli.FlexColumnWidthFour),
		table.NewFlexColumn("actions", "Actions", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("subject_condition_set_id", "Subject Condition Set: Id", cli.FlexColumnWidthFour),
		table.NewFlexColumn("subject_condition_set", "Subject Condition Set", cli.FlexColumnWidthThree),
	)
	rows := []table.Row{}
	for _, sm := range list {
		var actionsJSON []byte
		if actionsJSON, err = json.Marshal(sm.GetActions()); err != nil {
			cli.ExitWithError("Error marshalling subject mapping actions", err)
		}

		var subjectSetsJSON []byte
		if subjectSetsJSON, err = json.Marshal(sm.GetSubjectConditionSet().GetSubjectSets()); err != nil {
			cli.ExitWithError("Error marshalling subject condition set", err)
		}

		rows = append(rows, table.NewRow(table.RowData{
			"id":                       sm.GetId(),
			"value_id":                 sm.GetAttributeValue().GetId(),
			"value_fqn":                sm.GetAttributeValue().GetFqn(),
			"actions":                  string(actionsJSON),
			"subject_condition_set_id": sm.GetSubjectConditionSet().GetId(),
			"subject_condition_set":    string(subjectSetsJSON),
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, list)
}

func policy_createAction(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	attrValueId := c.Flags.GetRequiredID("attribute-value-id")
	actionsStandard = c.Flags.GetStringSlice("action-standard", actionsStandard, cli.FlagsStringSliceOptions{Min: 0})
	actionsCustom = c.Flags.GetStringSlice("action-custom", actionsCustom, cli.FlagsStringSliceOptions{Min: 0})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})
	existingSCSId := c.Flags.GetOptionalID("subject-condition-set-id")
	// NOTE: labels within a new Subject Condition Set created on a SM creation are not supported
	newScsJSON := c.Flags.GetOptionalString("subject-condition-set-new")

	// validations
	if len(actionsStandard) == 0 && len(actionsCustom) == 0 {
		cli.ExitWithError("At least one Standard or Custom Action [--action-standard, --action-custom] is required", nil)
	}
	if len(actionsStandard) > 0 {
		for _, a := range actionsStandard {
			a = strings.ToUpper(a)
			if a != actionStandardDecrypt && a != actionStandardTransmit {
				cli.ExitWithError(fmt.Sprintf("Invalid Standard Action: '%s'. Must be one of [DECRYPT, TRANSMIT].", a), nil)
			}
		}
	}
	if existingSCSId == "" && newScsJSON == "" {
		cli.ExitWithError("At least one Subject Condition Set flag [--subject-condition-set-id, --subject-condition-set-new] must be provided", nil)
	}

	actions := getFullActionsList(actionsStandard, actionsCustom)

	var scs *subjectmapping.SubjectConditionSetCreate
	if newScsJSON != "" {
		ss, err := unmarshalSubjectSetsProto([]byte(newScsJSON))
		if err != nil {
			cli.ExitWithError("Error unmarshalling subject sets", err)
		}
		scs = &subjectmapping.SubjectConditionSetCreate{
			SubjectSets: ss,
		}
	}

	mapping, err := h.CreateNewSubjectMapping(attrValueId, actions, existingSCSId, scs, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create subject mapping", err)
	}

	var actionsJSON []byte
	if actionsJSON, err = json.Marshal(mapping.GetActions()); err != nil {
		cli.ExitWithError("Error marshalling subject mapping actions", err)
	}

	var subjectSetsJSON []byte
	if mapping.GetSubjectConditionSet() != nil {
		if subjectSetsJSON, err = json.Marshal(mapping.GetSubjectConditionSet().GetSubjectSets()); err != nil {
			cli.ExitWithError("Error marshalling subject condition set", err)
		}
	}

	rows := [][]string{
		{"Id", mapping.GetId()},
		{"Attribute Value Id", mapping.GetAttributeValue().GetId()},
		{"Actions", string(actionsJSON)},
		{"Subject Condition Set: Id", mapping.GetSubjectConditionSet().GetId()},
		{"Subject Condition Set", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(mapping.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, mapping.GetId(), t, mapping)
}

func policy_deleteAction(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetOptionalBool("force")

	sm, err := h.GetSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find subject mapping (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDelete, "subject mapping", sm.GetId(), force)

	deleted, err := h.DeleteSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to delete subject mapping (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	rows := [][]string{{"Id", sm.GetId()}}
	if mdRows := getMetadataRows(deleted.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, deleted)
}

func policy_updateAction(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	actionsStandard = c.Flags.GetStringSlice("action-standard", actionsStandard, cli.FlagsStringSliceOptions{Min: 0})
	actionsCustom = c.Flags.GetStringSlice("action-custom", actionsCustom, cli.FlagsStringSliceOptions{Min: 0})
	scsId := c.Flags.GetOptionalID("subject-condition-set-id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if len(actionsStandard) > 0 {
		for _, a := range actionsStandard {
			a = strings.ToUpper(a)
			if a != actionStandardDecrypt && a != actionStandardTransmit {
				cli.ExitWithError(fmt.Sprintf("Invalid Standard Action: '%s'. Must be one of [ENCRYPT, TRANSMIT]. Other actions must be custom.", a), nil)
			}
		}
	}
	actions := getFullActionsList(actionsStandard, actionsCustom)

	updated, err := h.UpdateSubjectMapping(
		id,
		scsId,
		actions,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update subject mapping", err)
	}
	rows := [][]string{{"Id", id}}
	if mdRows := getMetadataRows(updated.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, id, t, updated)
}

func init() {
	getDoc := man.Docs.GetCommand("policy/actions/get",
		man.WithRun(policy_getAction),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	listDoc := man.Docs.GetCommand("policy/actions/list",
		man.WithRun(policy_listActions),
	)
	injectListPaginationFlags(listDoc)

	createDoc := man.Docs.GetCommand("policy/actions/create",
		man.WithRun(policy_createAction),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("attribute-value-id").Name,
		createDoc.GetDocFlag("attribute-value-id").Shorthand,
		createDoc.GetDocFlag("attribute-value-id").Default,
		createDoc.GetDocFlag("attribute-value-id").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&actionsStandard,
		createDoc.GetDocFlag("action-standard").Name,
		createDoc.GetDocFlag("action-standard").Shorthand,
		[]string{},
		createDoc.GetDocFlag("action-standard").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&actionsCustom,
		createDoc.GetDocFlag("action-custom").Name,
		createDoc.GetDocFlag("action-custom").Shorthand,
		[]string{},
		createDoc.GetDocFlag("action-custom").Description,
	)
	createDoc.Flags().String(
		createDoc.GetDocFlag("subject-condition-set-id").Name,
		createDoc.GetDocFlag("subject-condition-set-id").Default,
		createDoc.GetDocFlag("subject-condition-set-id").Description,
	)
	createDoc.Flags().String(
		createDoc.GetDocFlag("subject-condition-set-new").Name,
		createDoc.GetDocFlag("subject-condition-set-new").Default,
		createDoc.GetDocFlag("subject-condition-set-new").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	updateDoc := man.Docs.GetCommand("policy/actions/update",
		man.WithRun(policy_updateAction),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().StringSliceVarP(
		&actionsStandard,
		updateDoc.GetDocFlag("action-standard").Name,
		updateDoc.GetDocFlag("action-standard").Shorthand,
		[]string{},
		updateDoc.GetDocFlag("action-standard").Description,
	)
	updateDoc.Flags().StringSliceVarP(
		&actionsCustom,
		updateDoc.GetDocFlag("action-custom").Name,
		updateDoc.GetDocFlag("action-custom").Shorthand,
		[]string{},
		updateDoc.GetDocFlag("action-custom").Description,
	)
	updateDoc.Flags().String(
		updateDoc.GetDocFlag("subject-condition-set-id").Name,
		updateDoc.GetDocFlag("subject-condition-set-id").Default,
		updateDoc.GetDocFlag("subject-condition-set-id").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	deleteDoc := man.Docs.GetCommand("policy/actions/delete",
		man.WithRun(policy_deleteAction),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)
	deleteDoc.Flags().Bool(
		deleteDoc.GetDocFlag("force").Name,
		false,
		deleteDoc.GetDocFlag("force").Description,
	)

	policy_ActionsDoc := man.Docs.GetCommand("policy/actions",
		man.WithSubcommands(
			getDoc,
			listDoc,
			createDoc,
			updateDoc,
			deleteDoc,
		),
	)
	policyCmd.AddCommand(&policy_ActionsDoc.Command)
}
