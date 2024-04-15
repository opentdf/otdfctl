package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

func policy_createSubjectConditionSet(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()
	var (
		ss      []*policy.SubjectSet
		ssBytes []byte
	)

	flagHelper := cli.NewFlagHelper(cmd)
	ssFlagJSON := flagHelper.GetOptionalString("subject-sets")
	ssFileJSON := flagHelper.GetOptionalString("subject-sets-file-json")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	// validate no flag conflicts
	if ssFileJSON == "" && ssFlagJSON == "" {
		cli.ExitWithError("At least one subject set must be provided ('--subject-sets', '--subject-sets-file-json')", nil)
	} else if ssFileJSON != "" && ssFlagJSON != "" {
		cli.ExitWithError("Only one of '--subject-sets' or '--subject-sets-file-json' can be provided", nil)
	}

	// read subject sets into bytes from either the flagged json file or json string
	if ssFileJSON != "" {
		jsonFile, err := os.Open(ssFileJSON)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to open file at path: %s", ssFileJSON), err)
		}
		defer jsonFile.Close()

		bytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", ssFileJSON), err)
		}
		ssBytes = bytes
	} else {
		ssBytes = []byte(ssFlagJSON)
	}

	if err := json.Unmarshal(ssBytes, &ss); err != nil {
		cli.ExitWithError("Error unmarshalling subject sets", err)
	}

	scs, err := h.CreateSubjectConditionSet(ss, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Error creating subject condition set", err)
	}

	var subjectSetsJSON []byte
	if subjectSetsJSON, err = json.Marshal(scs.SubjectSets); err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.Id},
		{"SubjectSets", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(scs.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, scs.Id, t, scs)
}

func policy_getSubjectConditionSet(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	scs, err := h.GetSubjectConditionSet(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
	}
	var subjectSetsJSON []byte
	if subjectSetsJSON, err = json.Marshal(scs.SubjectSets); err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.Id},
		{"SubjectSets", string(subjectSetsJSON)},
	}
	if mdRows := getMetadataRows(scs.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, scs.Id, t, scs)
}

func policy_listSubjectConditionSets(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	scsList, err := h.ListSubjectConditionSets()
	if err != nil {
		cli.ExitWithError("Error listing subject condition sets", err)
	}

	t := cli.NewTable()
	t.Headers("Id", "SubjectSets", "Labels", "Created At", "Updated At")
	for _, scs := range scsList {
		var subjectSetsJSON []byte
		if subjectSetsJSON, err = json.Marshal(scs.SubjectSets); err != nil {
			cli.ExitWithError("Error marshalling subject condition set", err)
		}
		metadata := cli.ConstructMetadata(scs.Metadata)
		rowCells := []string{scs.Id, string(subjectSetsJSON), metadata["Labels"], metadata["Created At"], metadata["Updated At"]}
		t.Row(rowCells...)
	}

	HandleSuccess(cmd, "", t, scsList)
}

func policy_updateSubjectConditionSet(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})
	ssFlagJSON := flagHelper.GetOptionalString("subject-sets")

	var ss []*policy.SubjectSet
	if ssFlagJSON != "" {
		if err := json.Unmarshal([]byte(ssFlagJSON), &ss); err != nil {
			cli.ExitWithError("Error unmarshalling subject sets", err)
		}
	}

	_, err := h.UpdateSubjectConditionSet(id, ss, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Error updating subject condition set", err)
	}

	scs, err := h.GetSubjectConditionSet(id)
	if err != nil {
		cli.ExitWithError("Error getting subject condition set", err)
	}

	var subjectSetsJSON []byte
	if subjectSetsJSON, err = json.Marshal(scs.SubjectSets); err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.Id},
		{"SubjectSets", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(scs.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, scs.Id, t, scs)
}

func policy_deleteSubjectConditionSet(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	scs, err := h.GetSubjectConditionSet(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
	}

	cli.ConfirmAction(cli.ActionDelete, "Subject Condition Set", id)

	if err := h.DeleteSubjectConditionSet(id); err != nil {
		cli.ExitWithError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
	}

	var subjectSetsJSON []byte
	if subjectSetsJSON, err = json.Marshal(scs.SubjectSets); err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.Id},
		{"SubjectSets", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(scs.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, scs.Id, t, scs)
}

var policy_subjectConditionSetsCmd *cobra.Command

func init() {
	createDoc := man.Docs.GetCommand("policy/subject-condition-sets/create",
		man.WithRun(policy_createSubjectConditionSet),
	)
	injectLabelFlags(&createDoc.Command, false)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("subject-sets").Name,
		createDoc.GetDocFlag("subject-sets").Shorthand,
		createDoc.GetDocFlag("subject-sets").Default,
		createDoc.GetDocFlag("subject-sets").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("subject-sets-file-json").Name,
		createDoc.GetDocFlag("subject-sets-file-json").Shorthand,
		createDoc.GetDocFlag("subject-sets-file-json").Default,
		createDoc.GetDocFlag("subject-sets-file-json").Description,
	)

	getDoc := man.Docs.GetCommand("policy/subject-condition-sets/get",
		man.WithRun(policy_getSubjectConditionSet),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	listDoc := man.Docs.GetDoc("policy/subject-condition-sets/list")
	listDoc.Run = policy_listSubjectConditionSets

	updateDoc := man.Docs.GetCommand("policy/subject-condition-sets/update",
		man.WithRun(policy_updateSubjectConditionSet),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("subject-sets").Name,
		updateDoc.GetDocFlag("subject-sets").Shorthand,
		updateDoc.GetDocFlag("subject-sets").Default,
		updateDoc.GetDocFlag("subject-sets").Description,
	)

	deleteDoc := man.Docs.GetCommand(
		"policy/subject-condition-sets/delete",
		man.WithRun(policy_deleteSubjectConditionSet),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)

	doc := man.Docs.GetCommand("policy/subject-condition-sets",
		man.WithSubcommands(
			createDoc,
			getDoc,
			updateDoc,
			deleteDoc,
		),
	)
	policy_subjectConditionSetsCmd = &doc.Command
	policyCmd.AddCommand(&doc.Command)
}

func getSubjectConditionSetOperatorFromChoice(choice string) (policy.SubjectMappingOperatorEnum, error) {
	switch choice {
	case "IN":
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN, nil
	case "NOT_IN":
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN, nil
	default:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED, fmt.Errorf("Unknown operator must be specified ['IN', 'NOT_IN']: %s", choice)
	}
}

func getSubjectConditionSetBooleanTypeFromChoice(choice string) (policy.ConditionBooleanTypeEnum, error) {
	switch choice {
	case "AND":
		return policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_AND, nil
	case "OR":
		return policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_OR, nil
	default:
		return policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_UNSPECIFIED, fmt.Errorf("Unknown boolean type must be specified ['AND', 'OR']: %s", choice)
	}
}

func getMarshaledSubjectSets(subjectSets string) ([]*policy.SubjectSet, error) {
	var ss []*policy.SubjectSet

	if err := json.Unmarshal([]byte(subjectSets), &ss); err != nil {
		return nil, err
	}

	return ss, nil
}
