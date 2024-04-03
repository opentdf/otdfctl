package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/man"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

var (
	standardActions []string
	customActions   []string
)

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

func runGetSubjectMappings(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	mapping, err := h.GetSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Could not find subject mapping (%s)", id)
		cli.ExitWithNotFoundError(errMsg, err)
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

func runListSubjectMappings(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	list, err := h.ListSubjectMappings()
	if err != nil {
		cli.ExitWithError("Could not get subject mappings", err)
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

func runCreateSubjectMappings(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	attrValueId := flagHelper.GetRequiredString("attribute-value-id")
	standardActions := flagHelper.GetStringSlice("action-standard", standardActions, cli.FlagHelperStringSliceOptions{Min: 0})
	customActions := flagHelper.GetStringSlice("action-custom", customActions, cli.FlagHelperStringSliceOptions{Min: 0})
	existingSCSId := flagHelper.GetOptionalString("subject-condition-set-id")
	// TODO: do we need to support creating a SM & SCS simultaneously? If so, it gets more complex.
	// newScs := flagHelper.GetOptionalString("new-subject-condition-set")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	// validations
	if len(standardActions) == 0 && len(customActions) == 0 {
		cli.ExitWithError("At least one Standard or Custom Action [--action-standard, --action-custom] is required", nil)
	}
	if len(standardActions) > 0 {
		for _, a := range standardActions {
			a = strings.ToUpper(a)
			if a != "DECRYPT" && a != "TRANSMIT" {
				cli.ExitWithError(fmt.Sprintf("Invalid Standard Action: '%s'. Must be one of [ENCRYPT, TRANSMIT].", a), nil)
			}
		}
	}
	actions := getFullActionsList(standardActions, customActions)

	mapping, err := h.CreateNewSubjectMapping(attrValueId, actions, existingSCSId, nil, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Could not create subject mapping", err)
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

func runUpdateSubjectMappings(cmd *cobra.Command, args []string) {
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
		cli.ExitWithError("Could not update subject mapping", err)
	}

	HandleSuccess(cmd, id, nil, updated)
}

func runDeleteSubjectMappings(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	sm, err := h.GetSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Could not find subject mapping (%s)", id)
		cli.ExitWithNotFoundError(errMsg, err)
	}

	cli.ConfirmDelete("subject mapping", sm.Id)

	deleted, err := h.DeleteSubjectMapping(id)
	if err != nil {
		errMsg := fmt.Sprintf("Could not delete subject mapping (%s)", id)
		cli.ExitWithNotFoundError(errMsg, err)
	}
	HandleSuccess(cmd, id, nil, deleted)
}

func init() {
	cmd := man.Docs.GetDoc("policy/subject-mappings")

	createCmd := man.Docs.GetDoc("policy/subject-mappings/create")
	createCmd.Run = runCreateSubjectMappings

	listCmd := man.Docs.GetDoc("policy/subject-mappings/list")
	listCmd.Run = runListSubjectMappings

	updateCmd := man.Docs.GetDoc("policy/subject-mappings/update")
	updateCmd.Run = runUpdateSubjectMappings

	deleteCmd := man.Docs.GetDoc("policy/subject-mappings/delete")
	deleteCmd.Run = runDeleteSubjectMappings

	getCmd := man.Docs.GetDoc("policy/subject-mappings/get")
	getCmd.Run = runGetSubjectMappings

	cmd.Short = cmd.GetShort([]string{
		createCmd.Use,
		getCmd.Use,
		listCmd.Use,
		updateCmd.Use,
		deleteCmd.Use,
	})

	policyCmd.AddCommand(&cmd.Command)

	cmd.AddCommand(&getCmd.Command)
	getCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")

	cmd.AddCommand(&listCmd.Command)

	cmd.AddCommand(&createCmd.Command)
	createCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the mapped Attribute Value")
	createCmd.Flags().StringSliceVarP(&standardActions, "action-standard", "s", []string{}, "Standard Action: [DECRYPT, TRANSMIT]")
	createCmd.Flags().StringSliceVarP(&customActions, "action-custom", "c", []string{}, "Custom Action")
	createCmd.Flags().String("subject-condition-set-id", "", "Pre-existing Subject Condition Set Id")
	// TODO: do we need to support creating a SM & SCS simultaneously? If so, it gets more complex.
	// policy_subject_mappingCreateCmd.Flags().StringP("new-subject-condition-set", "scs", "", "New Subject Condition Set (optional)")
	injectLabelFlags(&createCmd.Command, false)

	cmd.AddCommand(&updateCmd.Command)
	updateCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
	updateCmd.Flags().StringSliceVarP(&standardActions, "action-standard", "s", []string{}, "Standard Action: [DECRYPT, TRANSMIT]. Note: destructively replaces existing Actions.")
	updateCmd.Flags().StringSliceVarP(&customActions, "action-custom", "c", []string{}, "Custom Action. Note: destructively replaces existing Actions.")
	updateCmd.Flags().String("subject-condition-set-id", "", "Updated Subject Condition Set Id")
	injectLabelFlags(&updateCmd.Command, true)

	cmd.AddCommand(&deleteCmd.Command)
	deleteCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
}
