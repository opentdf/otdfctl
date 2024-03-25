package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/tructl/pkg/cli"
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

	subjectSets []*policy.SubjectSet

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
			var ss []*policy.SubjectSet

			flagHelper := cli.NewFlagHelper(cmd)
			// ss := flagHelper.GetStringSlice("subject-set", subjectSets, cli.FlagHelperStringSliceOptions{Min: 1})
			subjectSetFile := flagHelper.GetOptionalString("subject-set-from-file")
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			jsonFile, err := os.Open(subjectSetFile)
			if err != nil {
				cli.ExitWithError(fmt.Sprintf("Failed to open file %s", subjectSetFile), err)
			}
			defer jsonFile.Close()

			bytes, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				cli.ExitWithError(fmt.Sprintf("Failed to read file %s", subjectSetFile), err)
			}
			if err := json.Unmarshal(bytes, &ss); err != nil {
				cli.ExitWithError(fmt.Sprintf("Failed to unmarshal file contents %s", string(bytes)), err)
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
				cli.ExitWithNotFoundError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
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
				cli.ExitWithNotFoundError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
			}

			cli.ConfirmDelete("Subject Condition Set", id)

			if err := h.DeleteSubjectConditionSet(id); err != nil {
				cli.ExitWithNotFoundError(fmt.Sprintf("Subject Condition Set with id %s not found", id), err)
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
	// policy_subject_condition_setCreateCmd.Flags().StringSliceVarP(&subjectSets, "subject-set", "s", []string{}, "A subject set, containing a list of condition groups, each with one or more conditions.")
	policy_subject_condition_setCreateCmd.Flags().String("subject-set-from-file", "", "A JSON file with path from $HOME containing a subject set")

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setGetCmd)
	policy_subject_condition_setGetCmd.Flags().StringP("id", "i", "", "Id of the subject condition set")

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setListCmd)

	policy_subject_condition_setCmd.AddCommand(policy_subject_condition_setUpdateCmd)
	policy_subject_condition_setUpdateCmd.Flags().StringP("id", "i", "", "Id of the subject condition set")
	injectLabelFlags(policy_subject_condition_setUpdateCmd, true)

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

func getMarshaledSubjectSets(subjectSets []string) ([]*policy.SubjectSet, error) {
	var ss []*policy.SubjectSet
	for _, subjectSet := range subjectSets {

		var s policy.SubjectSet
		if err := json.Unmarshal([]byte(subjectSet), &s); err != nil {
			return nil, err
		}
		// for _, cg := range s.ConditionGroups {
		// 	for _, c := range cg.Conditions {
		// 		op, err := getSubjectConditionSetOperatorFromChoice(c.Operator)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		c.Operator = op.String()
		// 	}
		// }
		ss = append(ss, &s)
	}
	return ss, nil
}
