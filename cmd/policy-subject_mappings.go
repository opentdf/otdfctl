package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

var (
	subjectMappingsCmds = []string{
		subjectMappingCreateCmd.Use,
		subjectMappingGetCmd.Use,
		subjectMappingsListCmd.Use,
		subjectMappingUpdateCmd.Use,
		subjectMappingDeleteCmd.Use,
	}

	subjectValues []string

	subjectMappingsCmd = &cobra.Command{
		Use:   "subject-mappings",
		Short: "Manage subject mappings [" + strings.Join(subjectMappingsCmds, ", ") + "]",
		Long: `
Subject Mappings - commands to manage relationships between subjects (PEs, NPEs, etc) and attributes.
		
For example: a subject mapping could be created such that the AcmeCorp engineering
team member named "Alice" is "IN" the value "Engineering" for attribute "Teams" in
namespace "acmecorp.com", but is not mapped to the attribute value "Sales" within the
same attribute and namespace. 
`,
	}

	subjectMappingGetCmd = &cobra.Command{
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

			if !jsonOutput {
				cli.PrintSuccessTable(cmd, id, cli.NewTabular().Rows(rows...))
			} else {
				if output, err := json.MarshalIndent(mapping, "", "  "); err != nil {
					cli.ExitWithError("Error marshalling subject mapping", err)
				} else {
					fmt.Println(string(output))
				}
			}
		},
	}

	subjectMappingsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List subject mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			list, err := h.ListSubjectMappings()
			if err != nil {
				cli.ExitWithError("Could not get subject mappings", err)
			}

			if jsonOutput {
				if output, err := json.MarshalIndent(list, "", "  "); err != nil {
					cli.ExitWithError("Error marshalling subject mappings", err)
				} else {
					fmt.Println(string(output))
				}
				return
			}

			t := cli.NewTable().Width(180)
			t.Headers("Id", "Subject Attribute", "Operator", "Subject Values", "Attribute Value ID")
			for _, sm := range list {
				rowCells := []string{
					sm.Id,
					sm.SubjectAttribute,
					handlers.GetSubjectMappingOperatorChoiceFromEnum(sm.Operator),
					strings.Join(sm.SubjectValues, ", "),
					sm.AttributeValue.Id,
				}
				t.Row(rowCells...)
			}
			cli.PrintSuccessTable(cmd, "", t)
		},
	}

	subjectMappingCreateCmd = &cobra.Command{
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

			if jsonOutput {
				if output, err := json.MarshalIndent(mapping, "", "  "); err != nil {
					cli.ExitWithError("Error marshalling subject mapping", err)
				} else {
					fmt.Println(string(output))
				}
				return
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

			cli.PrintSuccessTable(cmd, mapping.Id,
				cli.NewTabular().
					Rows(rows...))
		},
	}

	subjectMappingDeleteCmd = &cobra.Command{
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

	subjectMappingUpdateCmd = &cobra.Command{
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
	policyCmd.AddCommand(subjectMappingsCmd)

	subjectMappingsCmd.AddCommand(subjectMappingGetCmd)
	subjectMappingGetCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")

	subjectMappingsCmd.AddCommand(subjectMappingsListCmd)

	subjectMappingsCmd.AddCommand(subjectMappingCreateCmd)
	subjectMappingCreateCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the attribute value")
	subjectMappingCreateCmd.Flags().StringP("subject-attribute", "s", "", "Subject attribute")
	subjectMappingCreateCmd.Flags().StringSliceVarP(&subjectValues, "subject-values", "v", []string{}, "Subject values")
	subjectMappingCreateCmd.Flags().StringP("operator", "o", "", "Operator")
	subjectMappingCreateCmd.Flags().StringP("metadata", "m", "", "Metadata (optional): labels and description")

	subjectMappingsCmd.AddCommand(subjectMappingUpdateCmd)
	subjectMappingUpdateCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
	subjectMappingUpdateCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the attribute value")
	subjectMappingUpdateCmd.Flags().StringP("subject-attribute", "s", "", "Subject attribute")
	subjectMappingUpdateCmd.Flags().StringSliceVarP(&subjectValues, "subject-values", "v", []string{}, "Subject values")
	subjectMappingUpdateCmd.Flags().StringP("operator", "o", "", "Operator: [IN, NOT_IN]")
	subjectMappingUpdateCmd.Flags().StringP("metadata", "m", "", "Metadata (optional): labels and description")

	subjectMappingsCmd.AddCommand(subjectMappingDeleteCmd)
	subjectMappingDeleteCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
}
