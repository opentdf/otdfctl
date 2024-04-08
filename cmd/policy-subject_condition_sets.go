package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

var (
	policy_subject_condition_setsCmds = []string{
		policy_subject_condition_setCreateCmd.Use,
		policy_subject_condition_setGetCmd.Use,
		policy_subject_condition_setListCmd.Use,
		policy_subject_condition_setUpdateCmd.Use,
		policy_subject_condition_setDeleteCmd.Use,
	}

	policy_subject_condition_setCmd = &cobra.Command{
		Use:   "subject-condition-sets",
		Short: "Manage subject condition sets" + strings.Join(policy_subject_condition_setsCmds, ", ") + "]",
		Long: `
Subject Condition Sets - fields and values known to an external user source that are utilized to relate a Subject (PE/NPE) to
a Subject Mapping and, by said mapping, an Attribute Value.`,
	}

	policy_subject_condition_setCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a subject condition set",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	policy_subject_condition_setGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a subject condition set by id",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	policy_subject_condition_setListCmd = &cobra.Command{
		Use:   "list",
		Short: "List subject condition sets",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			scsList, err := h.ListSubjectConditionSets()
			if err != nil {
				cli.ExitWithError("Error listing subject condition sets", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "SubjectSets")
			for _, scs := range scsList {
				var subjectSetsJSON []byte
				if subjectSetsJSON, err = json.Marshal(scs.SubjectSets); err != nil {
					cli.ExitWithError("Error marshalling subject condition set", err)
				}
				rowCells := []string{scs.Id, string(subjectSetsJSON)}
				t.Row(rowCells...)
			}

			HandleSuccess(cmd, "", t, scsList)
		},
	}

	policy_subject_condition_setUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a subject condition set",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	policy_subject_condition_setDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a subject condition set",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_subject_condition_setCmd)

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setCreateCmd)
	injectLabelFlags(policy_subject_condition_setCreateCmd, false)
	policy_subject_condition_setCreateCmd.Flags().StringP("subject-sets", "s", "", "A JSON array of subject sets, containing a list of condition groups, each with one or more conditions")
	policy_subject_condition_setCreateCmd.Flags().StringP("subject-sets-file-json", "j", "", "A JSON file with path from $HOME containing an array of subject sets")

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setGetCmd)
	policy_subject_condition_setGetCmd.Flags().StringP("id", "i", "", "Id of the subject condition set")

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setListCmd)

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setUpdateCmd)
	policy_subject_condition_setUpdateCmd.Flags().StringP("id", "i", "", "Id of the subject condition set")
	injectLabelFlags(policy_subject_condition_setUpdateCmd, true)
	policy_subject_condition_setUpdateCmd.Flags().StringP("subject-sets", "s", "", "A JSON array of subject sets, containing a list of condition groups, each with one or more conditions")

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setDeleteCmd)
	policy_subject_condition_setDeleteCmd.Flags().StringP("id", "i", "", "Id of the subject condition set")
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
