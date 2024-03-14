package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

var (
	policy_subject_mappingsCmds = []string{
		policy_subject_mappingCreateCmd.Use,
		policy_subject_mappingGetCmd.Use,
		policy_subject_mappingsListCmd.Use,
		policy_subject_mappingUpdateCmd.Use,
		policy_subject_mappingDeleteCmd.Use,
	}

	subjectValues []string

	policy_subject_mappingsCmd = &cobra.Command{
		Use:   "subject-mappings",
		Short: "Manage subject mappings [" + strings.Join(policy_subject_mappingsCmds, ", ") + "]",
		Long: `
Subject Mappings - commands to manage relationships between subjects (PEs, NPEs, etc) and attributes.

For example: a subject mapping could be created such that the AcmeCorp engineering
team member named "Alice" is "IN" the value "Engineering" for attribute "Teams" in
namespace "acmecorp.com", but is not mapped to the attribute value "Sales" within the
same attribute and namespace.
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
				cli.ExitWithError(errMsg, err)
			}

			rows := [][]string{
				{"Id", mapping.Id},
				{"Subject Attribute", mapping.SubjectAttribute},
				{"Operator", handlers.GetSubjectMappingOperatorChoiceFromEnum(mapping.Operator)},
				{"Subject Values", strings.Join(mapping.SubjectValues, ", ")},
			}

			if mdRows := getMetadataRows(mapping.Metadata); mdRows != nil {
				rows = append(rows, mdRows...)
			}

			cli.HandleSuccess(cmd, id, cli.NewTabular().Rows(rows...), mapping)
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
			t.Headers("Id", "Subject Attribute", "Operator", "Subject Values", "Attribute Value ID")
			for _, sm := range list {
				rowCells := []string{
					sm.Id,
					// sm.SubjectAttribute,
					// handlers.GetSubjectMappingOperatorChoiceFromEnum(sm.Operator),
					// strings.Join(sm.SubjectValues, ", "),
					sm.AttributeValue.Id,
				}
				t.Row(rowCells...)
			}
			cli.HandleSuccess(cmd, "", t, list)
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
			subjectAttribute := flagHelper.GetRequiredString("subject-attribute")
			subjectValues := flagHelper.GetStringSlice("subject-values", subjectValues, cli.FlagHelperStringSliceOptions{Min: 1})
			operator := flagHelper.GetRequiredString("operator")

			m := flagHelper.GetOptionalString("metadata")
			metadata := unMarshalMetadata(m)

			mapping, err := h.CreateNewSubjectMapping(attrValueId, subjectAttribute, subjectValues, operator, metadata)
			if err != nil {
				cli.ExitWithError("Could not create subject mapping", err)
			}

			rows := [][]string{
				{"Id", mapping.Id},
				{"Subject Attribute", mapping.SubjectAttribute},
				{"Operator", handlers.GetSubjectMappingOperatorChoiceFromEnum(mapping.Operator)},
				{"Subject Values", strings.Join(mapping.SubjectValues, ", ")},
				{"Attribute Value Id", mapping.AttributeValue.Id},
			}

			if mdRows := getMetadataRows(mapping.Metadata); mdRows != nil {
				rows = append(rows, mdRows...)
			}

			cli.HandleSuccess(cmd, mapping.Id,
				cli.NewTabular().
					Rows(rows...), mapping)
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
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmDelete("subject mapping", sm.Id)

			if err := h.DeleteSubjectMapping(id); err != nil {
				errMsg := fmt.Sprintf("Could not delete subject mapping (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			// TODO: handle json output once service sends back deleted subject mapping
			cli.PrintSuccessTable(cmd, id, nil)
		},
	}

	policy_subject_mappingUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a subject mapping",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			attrValueId := flagHelper.GetRequiredString("attribute-value-id")
			subjectAttribute := flagHelper.GetRequiredString("subject-attribute")
			subjectValues := flagHelper.GetStringSlice("subject-values", subjectValues, cli.FlagHelperStringSliceOptions{Min: 1})
			operator := flagHelper.GetRequiredString("operator")

			m := flagHelper.GetOptionalString("metadata")
			metadata := unMarshalMetadata(m)

			if _, err := h.UpdateSubjectMapping(
				id,
				attrValueId,
				subjectAttribute,
				subjectValues,
				operator,
				metadata,
			); err != nil {
				cli.ExitWithError("Could not update subject mapping", err)
			}

			// TODO: handle json output once service sends back updated subject mapping
			fmt.Println(cli.SuccessMessage(fmt.Sprintf("Subject mapping id: (%s) updated.", id)))
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_subject_mappingsCmd)

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingGetCmd)
	policy_subject_mappingGetCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingsListCmd)

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingCreateCmd)
	policy_subject_mappingCreateCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the attribute value")
	policy_subject_mappingCreateCmd.Flags().StringP("subject-attribute", "s", "", "Subject attribute")
	policy_subject_mappingCreateCmd.Flags().StringSliceVarP(&subjectValues, "subject-values", "v", []string{}, "Subject values")
	policy_subject_mappingCreateCmd.Flags().StringP("operator", "o", "", "Operator")
	policy_subject_mappingCreateCmd.Flags().StringP("metadata", "m", "", "Metadata (optional): labels and description")

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingUpdateCmd)
	policy_subject_mappingUpdateCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
	policy_subject_mappingUpdateCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the attribute value")
	policy_subject_mappingUpdateCmd.Flags().StringP("subject-attribute", "s", "", "Subject attribute")
	policy_subject_mappingUpdateCmd.Flags().StringSliceVarP(&subjectValues, "subject-values", "v", []string{}, "Subject values")
	policy_subject_mappingUpdateCmd.Flags().StringP("operator", "o", "", "Operator: [IN, NOT_IN]")
	policy_subject_mappingUpdateCmd.Flags().StringP("metadata", "m", "", "Metadata (optional): labels and description")

	policy_subject_mappingsCmd.AddCommand(policy_subject_mappingDeleteCmd)
	policy_subject_mappingDeleteCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
}
