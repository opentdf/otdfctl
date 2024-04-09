package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/otdfctl/issues/73] is addressed

var (
	standardActions []string
	customActions   []string
)

func policy_getSubjectMapping(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	mapping, err := h.GetSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find subject mapping (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	var actionsJSON []byte
	if actionsJSON, err = json.Marshal(mapping.Actions); err != nil {
		cli.ExitWithError("Error marshalling subject mapping actions", err)
	}

	var subjectSetsJSON []byte
	if subjectSetsJSON, err = json.Marshal(mapping.SubjectConditionSet.SubjectSets); err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", mapping.Id},
		{"Subject AttrVal: Id", mapping.AttributeValue.Id},
		{"Subject AttrVal: Value", mapping.AttributeValue.Value},
		{"Actions", string(actionsJSON)},
		{"Subject Condition Set: Id", mapping.SubjectConditionSet.Id},
		{"Subject Condition Set", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(mapping.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, mapping.Id, t, mapping)
}

func policy_listSubjectMappings(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	list, err := h.ListSubjectMappings()
	if err != nil {
		cli.ExitWithError("Failed to get subject mappings", err)
	}

	t := cli.NewTable().Width(180)
	t.Headers("Id", "Subject AttrVal: Id", "Subject AttrVal: Value", "Actions", "Subject Condition Set: Id", "Subject Condition Set")
	for _, sm := range list {
		var actionsJSON []byte
		if actionsJSON, err = json.Marshal(sm.Actions); err != nil {
			cli.ExitWithError("Error marshalling subject mapping actions", err)
		}

		var subjectSetsJSON []byte
		if subjectSetsJSON, err = json.Marshal(sm.SubjectConditionSet.SubjectSets); err != nil {
			cli.ExitWithError("Error marshalling subject condition set", err)
		}

		rowCells := []string{
			sm.Id,
			sm.AttributeValue.Id,
			sm.AttributeValue.Value,
			string(actionsJSON),
			sm.SubjectConditionSet.Id,
			string(subjectSetsJSON),
		}
		t.Row(rowCells...)
	}
	HandleSuccess(cmd, "", t, list)
}

func policy_createSubjectMapping(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	attrValueId := flagHelper.GetRequiredString("attribute-value-id")
	standardActions := flagHelper.GetStringSlice("action-standard", standardActions, cli.FlagHelperStringSliceOptions{Min: 0})
	customActions := flagHelper.GetStringSlice("action-custom", customActions, cli.FlagHelperStringSliceOptions{Min: 0})
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})
	existingSCSId := flagHelper.GetOptionalString("subject-condition-set-id")
	// NOTE: labels within a new Subject Condition Set created on a SM creation are not supported
	newScsJSON := flagHelper.GetOptionalString("subject-condition-set-new")

	// validations
	if len(standardActions) == 0 && len(customActions) == 0 {
		cli.ExitWithError("At least one Standard or Custom Action [--action-standard, --action-custom] is required", nil)
	}
	if len(standardActions) > 0 {
		for _, a := range standardActions {
			a = strings.ToUpper(a)
			if a != "DECRYPT" && a != "TRANSMIT" {
				cli.ExitWithError(fmt.Sprintf("Invalid Standard Action: '%s'. Must be one of [DECRYPT, TRANSMIT].", a), nil)
			}
		}
	}

	actions := getFullActionsList(standardActions, customActions)

	var ss []*policy.SubjectSet
	var scs *subjectmapping.SubjectConditionSetCreate
	if newScsJSON != "" {
		if err := json.Unmarshal([]byte(newScsJSON), &ss); err != nil {
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
	if actionsJSON, err = json.Marshal(mapping.Actions); err != nil {
		cli.ExitWithError("Error marshalling subject mapping actions", err)
	}

	var subjectSetsJSON []byte
	if mapping.SubjectConditionSet != nil {
		if subjectSetsJSON, err = json.Marshal(mapping.SubjectConditionSet.SubjectSets); err != nil {
			cli.ExitWithError("Error marshalling subject condition set", err)
		}
	}

	rows := [][]string{
		{"Id", mapping.Id},
		{"Subject AttrVal: Id", mapping.AttributeValue.Id},
		{"Actions", string(actionsJSON)},
		{"Subject Condition Set: Id", mapping.SubjectConditionSet.Id},
		{"Subject Condition Set", string(subjectSetsJSON)},
		{"Attribute Value Id", mapping.AttributeValue.Id},
	}

	if mdRows := getMetadataRows(mapping.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, mapping.Id, t, mapping)
}

func policy_deleteSubjectMapping(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	sm, err := h.GetSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find subject mapping (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDelete, "subject mapping", sm.Id)

	deleted, err := h.DeleteSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to delete subject mapping (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	HandleSuccess(cmd, id, nil, deleted)
}

func policy_updateSubjectMapping(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")
	standardActions := flagHelper.GetStringSlice("action-standard", standardActions, cli.FlagHelperStringSliceOptions{Min: 0})
	customActions := flagHelper.GetStringSlice("action-custom", customActions, cli.FlagHelperStringSliceOptions{Min: 0})
	scsId := flagHelper.GetOptionalString("subject-condition-set-id")
	labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	if len(standardActions) > 0 {
		for _, a := range standardActions {
			a = strings.ToUpper(a)
			if a != "DECRYPT" && a != "TRANSMIT" {
				cli.ExitWithError(fmt.Sprintf("Invalid Standard Action: '%s'. Must be one of [ENCRYPT, TRANSMIT]. Other actions must be custom.", a), nil)
			}
		}
	}
	actions := getFullActionsList(standardActions, customActions)

	updated, err := h.UpdateSubjectMapping(
		id,
		scsId,
		actions,
		getMetadataMutable(labels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update subject mapping", err)
	}

	HandleSuccess(cmd, id, nil, updated)
}

func getSubjectMappingMappingActionEnumFromChoice(readable string) policy.Action_StandardAction {
	switch readable {
	case "DECRYPT":
		return policy.Action_STANDARD_ACTION_DECRYPT
	case "TRANSMIT":
		return policy.Action_STANDARD_ACTION_TRANSMIT
	default:
		return policy.Action_STANDARD_ACTION_UNSPECIFIED
	}
}

func getFullActionsList(standardActions, customActions []string) []*policy.Action {
	actions := []*policy.Action{}
	for _, a := range standardActions {
		actions = append(actions, &policy.Action{
			Value: &policy.Action_Standard{
				Standard: getSubjectMappingMappingActionEnumFromChoice(a),
			},
		})
	}
	for _, a := range customActions {
		actions = append(actions, &policy.Action{
			Value: &policy.Action_Custom{
				Custom: a,
			},
		})
	}
	return actions
}

func init() {
	getDoc := man.Docs.GetCommand("policy/subject-mappings/get",
		man.WithRun(policy_getSubjectMapping),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	listDoc := man.Docs.GetCommand("policy/subject-mappings/list",
		man.WithRun(policy_listSubjectMappings),
	)

	createDoc := man.Docs.GetCommand("policy/subject-mappings/create",
		man.WithRun(policy_createSubjectMapping),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("attribute-value-id").Name,
		createDoc.GetDocFlag("attribute-value-id").Shorthand,
		createDoc.GetDocFlag("attribute-value-id").Default,
		createDoc.GetDocFlag("attribute-value-id").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&standardActions,
		createDoc.GetDocFlag("action-standard").Name,
		createDoc.GetDocFlag("action-standard").Shorthand,
		[]string{},
		createDoc.GetDocFlag("action-standard").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&customActions,
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

	updateDoc := man.Docs.GetCommand("policy/subject-mappings/update",
		man.WithRun(policy_updateSubjectMapping),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().StringSliceVarP(&standardActions,
		updateDoc.GetDocFlag("action-standard").Name,
		updateDoc.GetDocFlag("action-standard").Shorthand,
		[]string{},
		updateDoc.GetDocFlag("action-standard").Description,
	)
	updateDoc.Flags().StringSliceVarP(&customActions,
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

	deleteDoc := man.Docs.GetCommand("policy/subject-mappings/delete",
		man.WithRun(policy_deleteSubjectMapping),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)

	doc := man.Docs.GetCommand("policy/subject-mappings",
		man.WithSubcommands(
			createDoc,
			getDoc,
			listDoc,
			updateDoc,
			deleteDoc,
		),
	)
	policy_subjectMappingCmd := &doc.Command
	policyCmd.AddCommand(policy_subjectMappingCmd)
}
