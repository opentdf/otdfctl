package cmd

import (
	"errors"
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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			id := args[0]
			if id == "" {
				cli.ExitWithError("Invalid ID", errors.New(id))
			}

			mapping, err := h.GetSubjectMapping(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find subject mapping (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			rows := [][]string{
				{"Id", mapping.Id},
				{"Subject Attribute", mapping.SubjectAttribute},
				{"Subject Values", strings.Join(mapping.SubjectValues, ", ")},
				{"Operator", handlers.GetSubjectMappingOperatorChoiceFromEnum(mapping.Operator)},
				// TODO: render attribute here somehow
				// {"Attribute Value", mapping.AttributeValue.Value},
			}

			if mdRows := getMetadataRows(mapping.Metadata); mdRows != nil {
				rows = append(rows, mdRows...)
			}

			fmt.Println(cli.SuccessMessage("Subject mapping found"))
			fmt.Println(
				cli.NewTabular().
					Rows(rows...).
					Render(),
			)
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

			t := cli.NewTable()
			t.Headers("Id", "Subject Attribute", "Subject Values", "Operator" /* "Attribute Value",*/, "Metadata")
			for _, ns := range list {
				t.Row(
					ns.Id,
					ns.SubjectAttribute,
					strings.Join(ns.SubjectValues, ", "),
					handlers.GetSubjectMappingOperatorChoiceFromEnum(ns.Operator),
					// TODO: attribute values
					// TODO: metadata rows
				)
			}
			fmt.Println(t.Render())
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
			subjectValues := flagHelper.GetStringSlice("subject-values", subjectValues, cli.FlagHelperStringSliceOptions{ Min: 1 })
			operator := flagHelper.GetRequiredString("operator")
			// TODO: metadata

			mapping, err := h.CreateNewSubjectMapping(attrValueId, subjectAttribute, subjectValues, operator)
			if err != nil {
				cli.ExitWithError("Could not create subject mapping", err)
			}

			rows := [][]string{
				{"Id", mapping.Id},
				{"Subject Attribute", mapping.SubjectAttribute},
				{"Subject Values", strings.Join(mapping.SubjectValues, ", ")},
				{"Operator", handlers.GetSubjectMappingOperatorChoiceFromEnum(mapping.Operator)},
				// TODO: render attribute here somehow
				// {"Attribute Value", mapping.AttributeValue.Value},
			}

			if mdRows := getMetadataRows(mapping.Metadata); mdRows != nil {
				rows = append(rows, mdRows...)
			}

			fmt.Println(cli.SuccessMessage("Subject mapping found"))
			fmt.Println(
				cli.NewTabular().
					Rows(rows...).
					Render(),
			)
		},
	}

	subjectMappingDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a subject mapping by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			// id := args[0]
			// if id == "" {
			// 	fmt.Println(cli.ErrorMessage("Invalid ID", errors.New(id)))
			// 	os.Exit(1)
			// }
			// ns, err := h.GetSubjectMapping(id)
			// if err != nil {
			// 	errMsg := fmt.Sprintf("Could not find subject mapping (%s)", id)
			// 	cli.ExitWithNotFoundError(errMsg, err)
			// 	cli.ExitWithError(errMsg, err)
			// }

			// cli.ConfirmDelete("subject mapping", ns.Name)

			// if err := h.DeleteSubjectMapping(id); err != nil {
			// 	errMsg := fmt.Sprintf("Could not delete subject mapping (%s)", id)
			// 	cli.ExitWithNotFoundError(errMsg, err)
			// 	cli.ExitWithError(errMsg, err)
			// }

			// fmt.Println(cli.SuccessMessage("Subject mapping deleted"))
			// fmt.Println(
			// 	cli.NewTabular().
			// 		Rows([][]string{
			// 			{"Id", ns.Id},
			// 			{"Name", ns.Name},
			// 		}...).Render(),
			// )
		},
	}

	subjectMappingUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a subject mapping",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			// flagHelper := cli.NewFlagHelper(cmd)

			// id := flagHelper.GetRequiredString("id")
			// name := flagHelper.GetRequiredString("name")

			// if _, err := h.UpdateSubjectMapping(
			// 	id,
			// 	name,
			// ); err != nil {
			// 	cli.ExitWithError("Could not update subject mapping", err)
			// 	return
			// } else {
			// 	fmt.Println(cli.SuccessMessage(fmt.Sprintf("Subject mapping id: (%s) updated. Name set to (%s).", id, name)))
			// }
		},
	}
)

func init() {
	rootCmd.AddCommand(subjectMappingsCmd)

	subjectMappingsCmd.AddCommand(subjectMappingGetCmd)

	subjectMappingsCmd.AddCommand(subjectMappingsListCmd)

	subjectMappingsCmd.AddCommand(subjectMappingCreateCmd)
	subjectMappingCreateCmd.Flags().StringP("attribute-value-id", "a", "", "Id of the attribute value")
	subjectMappingCreateCmd.Flags().StringP("subject-attribute", "s", "", "Subject attribute")
	subjectMappingCreateCmd.Flags().StringSliceVarP(&subjectValues, "subject-values", "v", []string{}, "Subject values")
	subjectMappingCreateCmd.Flags().StringP("operator", "o", "", "Operator")


	subjectMappingsCmd.AddCommand(subjectMappingUpdateCmd)
	subjectMappingUpdateCmd.Flags().StringP("id", "i", "", "Id of the subject mapping")
	subjectMappingUpdateCmd.Flags().StringP("name", "n", "", "Name value of the subject mapping")

	subjectMappingsCmd.AddCommand(subjectMappingDeleteCmd)
}
