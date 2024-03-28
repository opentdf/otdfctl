package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

var (
	policy_subject_mappingsCmds = []string{
		policy_subject_mappingCreateCmd.Use,
		policy_subject_mappingGetCmd.Use,
		policy_subject_mappingsListCmd.Use,
		policy_subject_mappingUpdateCmd.Use,
		policy_subject_mappingDeleteCmd.Use,
	}

	standardActions []string
	customActions   []string

	policy_subject_mappingsCmd = &cobra.Command{
		Use:   "subject-mappings",
		Short: "Manage subject mappings [" + strings.Join(policy_subject_mappingsCmds, ", ") + "]",
		Long: `
Subject Mappings - relations between Attribute Values and Subject Condition Sets that define the allowed Actions.

If a User's properties match a Subject Condition Set, the corresponding Subject Mapping provides them a set of allowed Actions
on any Resource (data) containing the mapped Attribute Value. 

	Attribute Value  <------  Subject Mapping ------->  Subject Condition Set

	Subject Mapping: 
		- Attribute Value: associated Attribute Value that the Subject Mapping Actions are relevant to
		- Actions: permitted Actions a Subject can take on Resources containing the Attribute Value
		- Subject Condition Set: associated logical structure of external fields and values to match a Subject

Platform consumption flow:
Subject/User -> IdP/LDAP's External Fields & Values -> SubjectConditionSet -> SubjectMapping w/ Actions -> AttributeValue

Note: SubjectConditionSets are reusable among SubjectMappings and are available under separate 'policy' commands.
`,
	}

	policy_subject_mappingGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a subject mapping by id",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	policy_subject_mappingsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List subject mappings",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	policy_subject_mappingCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new subject mapping",
		Run: func(cmd *cobra.Command, args []string) {
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
						cli.ExitWithError(fmt.Sprintf("Invalid Standard Action: '%s'. Must be one of [DECRYPT, TRANSMIT].", a), nil)
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
		},
	}

	policy_subject_mappingDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a subject mapping by id",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	policy_subject_mappingUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a subject mapping",
		Long: `
Update a Subject Mapping by id.
'Actions' are updated in place, destructively replacing the current set. If you want to add or remove actions, you must provide the
full set of actions on update. `,
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}
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

func init() {
	policyCmd.AddCommand(policy_subject_mappingsCmd)

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingGetCmd)
	policy_subject_mappingGetCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingsListCmd)

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingCreateCmd)
	policy_subject_mappingCreateCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the mapped Attribute Value")
	policy_subject_mappingCreateCmd.Flags().StringSliceVarP(&standardActions, "action-standard", "s", []string{}, "Standard Action: [DECRYPT, TRANSMIT]")
	policy_subject_mappingCreateCmd.Flags().StringSliceVarP(&customActions, "action-custom", "c", []string{}, "Custom Action")
	policy_subject_mappingCreateCmd.Flags().String("subject-condition-set-id", "", "Pre-existing Subject Condition Set Id")
	// TODO: do we need to support creating a SM & SCS simultaneously? If so, it gets more complex.
	// policy_subject_mappingCreateCmd.Flags().StringP("new-subject-condition-set", "scs", "", "New Subject Condition Set (optional)")
	injectLabelFlags(policy_subject_mappingCreateCmd, false)

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingUpdateCmd)
	policy_subject_mappingUpdateCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
	policy_subject_mappingUpdateCmd.Flags().StringSliceVarP(&standardActions, "action-standard", "s", []string{}, "Standard Action: [DECRYPT, TRANSMIT]. Note: destructively replaces existing Actions.")
	policy_subject_mappingUpdateCmd.Flags().StringSliceVarP(&customActions, "action-custom", "c", []string{}, "Custom Action. Note: destructively replaces existing Actions.")
	policy_subject_mappingUpdateCmd.Flags().String("subject-condition-set-id", "", "Updated Subject Condition Set Id")
	injectLabelFlags(policy_subject_mappingUpdateCmd, true)

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingDeleteCmd)
	policy_subject_mappingDeleteCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
}
