package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// Helper to unmarshal SubjectSets from JSON (stored as JSONB in the database column)
func unmarshalSubjectSetsProto(conditionJSON []byte) ([]*policy.SubjectSet, error) {
	var (
		raw []json.RawMessage
		ss  []*policy.SubjectSet
	)
	if err := json.Unmarshal(conditionJSON, &raw); err != nil {
		return nil, err
	}

	for _, r := range raw {
		s := policy.SubjectSet{}
		if err := protojson.Unmarshal(r, &s); err != nil {
			return nil, err
		}
		ss = append(ss, &s)
	}

	return ss, nil
}

// Helper to marshal SubjectSets into JSON (stored as JSONB in the database column)
func marshalSubjectSetsProto(subjectSet []*policy.SubjectSet) ([]byte, error) {
	var raw []json.RawMessage
	for _, ss := range subjectSet {
		b, err := protojson.Marshal(ss)
		if err != nil {
			return nil, err
		}
		raw = append(raw, b)
	}
	return json.Marshal(raw)
}

func policy_createSubjectConditionSet(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()
	var ssBytes []byte

	ssFlagJSON := c.Flags.GetOptionalString("subject-sets")
	ssFileJSON := c.Flags.GetOptionalString("subject-sets-file-json")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

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

		bytes, err := io.ReadAll(jsonFile)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", ssFileJSON), err)
		}
		ssBytes = bytes
	} else {
		ssBytes = []byte(ssFlagJSON)
	}

	ss, err := unmarshalSubjectSetsProto(ssBytes)
	if err != nil {
		cli.ExitWithError("Error unmarshalling subject sets", err)
	}

	scs, err := h.CreateSubjectConditionSet(ss, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Error creating subject condition set", err)
	}

	subjectSetsJSON, err := marshalSubjectSetsProto(scs.GetSubjectSets())
	if err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.GetId()},
		{"SubjectSets", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(scs.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, scs.GetId(), t, scs)
}

func policy_getSubjectConditionSet(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	scs, err := h.GetSubjectConditionSet(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
	}
	subjectSetsJSON, err := marshalSubjectSetsProto(scs.GetSubjectSets())
	if err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.GetId()},
		{"SubjectSets", string(subjectSetsJSON)},
	}
	if mdRows := getMetadataRows(scs.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, scs.GetId(), t, scs)
}

func policy_listSubjectConditionSets(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	scsList, err := h.ListSubjectConditionSets()
	if err != nil {
		cli.ExitWithError("Error listing subject condition sets", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("subject_sets", "SubjectSets", cli.FlexColumnWidthFour),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, scs := range scsList {
		subjectSetsJSON, err := marshalSubjectSetsProto(scs.GetSubjectSets())
		if err != nil {
			cli.ExitWithError("Error marshalling subject condition set", err)
		}
		metadata := cli.ConstructMetadata(scs.GetMetadata())
		rows = append(rows, table.NewRow(table.RowData{
			"id":           scs.GetId(),
			"subject_sets": string(subjectSetsJSON),
			"labels":       metadata["Labels"],
			"created_at":   metadata["Created At"],
			"updated_at":   metadata["Updated At"],
		}))
	}
	t = t.WithRows(rows)
	HandleSuccess(cmd, "", t, scsList)
}

func policy_updateSubjectConditionSet(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})
	ssFlagJSON := c.Flags.GetOptionalString("subject-sets")
	ssFileJSON := c.Flags.GetOptionalString("subject-sets-file-json")

	var ssBytes []byte
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

		bytes, err := io.ReadAll(jsonFile)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", ssFileJSON), err)
		}
		ssBytes = bytes
	} else {
		ssBytes = []byte(ssFlagJSON)
	}

	ss, err := unmarshalSubjectSetsProto(ssBytes)
	if err != nil {
		cli.ExitWithError("Error unmarshalling subject sets", err)
	}

	_, err = h.UpdateSubjectConditionSet(id, ss, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Error updating subject condition set", err)
	}

	scs, err := h.GetSubjectConditionSet(id)
	if err != nil {
		cli.ExitWithError("Error getting subject condition set", err)
	}

	subjectSetsJSON, err := marshalSubjectSetsProto(scs.GetSubjectSets())
	if err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.GetId()},
		{"SubjectSets", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(scs.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, scs.GetId(), t, scs)
}

func policy_deleteSubjectConditionSet(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	scs, err := h.GetSubjectConditionSet(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
	}

	cli.ConfirmAction(cli.ActionDelete, "Subject Condition Set", id, false)

	if err := h.DeleteSubjectConditionSet(id); err != nil {
		cli.ExitWithError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
	}

	subjectSetsJSON, err := marshalSubjectSetsProto(scs.GetSubjectSets())
	if err != nil {
		cli.ExitWithError("Error marshalling subject condition set", err)
	}

	rows := [][]string{
		{"Id", scs.GetId()},
		{"SubjectSets", string(subjectSetsJSON)},
	}

	if mdRows := getMetadataRows(scs.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, scs.GetId(), t, scs)
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

	listDoc := man.Docs.GetCommand("policy/subject-condition-sets/list",
		man.WithRun(policy_listSubjectConditionSets),
	)

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
	updateDoc.Flags().StringP(
		createDoc.GetDocFlag("subject-sets-file-json").Name,
		createDoc.GetDocFlag("subject-sets-file-json").Shorthand,
		createDoc.GetDocFlag("subject-sets-file-json").Default,
		createDoc.GetDocFlag("subject-sets-file-json").Description,
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
			listDoc,
			updateDoc,
			deleteDoc,
		),
	)
	policy_subjectConditionSetsCmd = &doc.Command
	policyCmd.AddCommand(policy_subjectConditionSetsCmd)
}
