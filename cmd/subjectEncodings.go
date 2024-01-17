package cmd

import (
	"fmt"
	"strconv"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

var (
	resourceSelectorLabels []string
	attrValueLabels        []string
	subjectValues          []string
)

// acse represents the Access Control Subject Mappings command
var subjectMappingsCmd = &cobra.Command{
	Use:   "subject mappings",
	Short: "Access Control Subject Mappings/Encodings CRUD operations",
	Long: `Manage your configured Subject Mappings/Encodings [Create, Get one, List all, Update, Delete]
	through use of this CLI.`,
}

var subjectMappingsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Access Control Subject Mapping",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			attrRefLabels map[string]string
		)
		close := cli.GrpcConnect(cmd)
		defer close()

		h := cli.NewFlagHelper(cmd)

		
		name := h.GetRequiredString("name")
		description := h.GetOptionalString("description")
		resourceDeps := h.GetStringSlice("resource-dependencies", resourceDependencies, cli.FlagHelperListOptions{Min: 0})

		operator := h.GetRequiredString("operator")
		subjectAttr := h.GetRequiredString("subject-attribute")
		subjectValues := h.GetStringSlice("subject-values", subjectValues, cli.FlagHelperListOptions{Min: 1})

		attrRefName := h.GetOptionalString("attribute-ref-name")
		if attrRefName == "" {
			if len(attrValueLabels) == 0 { // optional, we'll ignore for now
				cli.ExitWithError("Either attribute-ref-name or attribute-ref-labels must be specified", nil)
			}
			attrRefLabels = h.GetKeyValuesMap("attribute-ref-labels", attrValueLabels, cli.FlagHelperListOptions{Min: 1})
		}

		if err := handlers.CreateSubjectMapping(
			handlers.SubjectMapping{
				Name:          name,
				Operator:      operator,
				SubjectAttr:   subjectAttr,
				SubjectValues: subjectValues,
			},
			description,
			resourceDeps,
			attrRefName,
			attrRefLabels,
		); err != nil {
			cli.ExitWithError("Could not create subject mapping", err)
		}

		fmt.Println(cli.SuccessMessage("Access Control Subject Mapping created"))
		fmt.Println(
			cli.NewTabular().
				Rows([][]string{
					{"Name", mapping.Name},
					{"Subject Attribute", mapping.SubjectAttr},
					{"Operator", mapping.Operator},
					{"Subject Values", cli.CommaSeparated(mapping.SubjectValues)},
				}...).Render(),
		)
	},
}

// Get one Access Control Subject Mapping
var subjectMappingsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an Access Control Subject Mapping",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cli.ExitWithError("Invalid ID", err)
		}

		close := cli.GrpcConnect(cmd)
		defer close()

		mapping, err := handlers.GetSubjectMapping(id)
		if err != nil {
			errMsg := fmt.Sprintf("Could not find attribute (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		fmt.Println(cli.SuccessMessage("Access Control Subject Mapping found"))
		fmt.Println(
			cli.NewTabular().
				Rows([][]string{
					{"Name", mapping.Name},
					{"Subject Attribute", mapping.SubjectAttr},
					{"Operator", mapping.Operator},
					{"Subject Values", cli.CommaSeparated(mapping.SubjectValues)},
				}...).Render(),
		)
	},
}

// List all access control subject mappings
var subjectMappingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list Access Control Subject Mappings",
	Run: func(cmd *cobra.Command, args []string) {
		close := cli.GrpcConnect(cmd)
		defer close()

		var (
			selectorName   string
			selectorLabels map[string]string
		)

		h := cli.NewFlagHelper(cmd)

		selectorName = h.GetOptionalString("resource-selector-name")
		if selectorName == "" {
			if len(resourceSelectorLabels) == 0 {
				cli.ExitWithError("Either resource-selector-name or resource-selector-labels must be specified", nil)
			}
			selectorLabels = h.GetKeyValuesMap("resource-selector-labels", resourceSelectorLabels, cli.FlagHelperListOptions{Min: 1})
		}

		mappings, err := handlers.ListSubjectMappings(selectorName, selectorLabels)
		if err != nil {
			cli.ExitWithError("Could not list subject mappings", err)
		}

		t := cli.NewTable()
		t.Headers("Name", "Subject Attribute", "Operator", "Subject Values")
		for _, m := range mappings {
			t.Row(
				m.Name,
				m.SubjectAttr,
				m.Operator,
				cli.CommaSeparated(m.SubjectValues),
			)
		}
		fmt.Println(t.Render())
	},
}

// Delete one Access Control Subject Mapping
var subjectMappingsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an Access Control Subject Mapping",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cli.ExitWithError("Invalid ID", err)
		}

		close := cli.GrpcConnect(cmd)
		defer close()

		mapping, err := handlers.GetSubjectMapping(id)
		if err != nil {
			errMsg := fmt.Sprintf("Could not find subject mapping (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		cli.ConfirmDelete("subject mapping", mapping.Name)

		if err := handlers.DeleteSubjectMapping(id); err != nil {
			errMsg := fmt.Sprintf("Could not delete subject mapping (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		fmt.Println(cli.SuccessMessage("Access Control Subject Mapping deleted"))
		fmt.Println(
			cli.NewTabular().
				Rows([][]string{
					{"Name", mapping.Name},
					{"Subject Attribute", mapping.SubjectAttr},
					{"Operator", mapping.Operator},
					{"Subject Values", cli.CommaSeparated(mapping.SubjectValues)},
				}...).Render(),
		)
	},
}

func init() {
	rootCmd.AddCommand(subjectMappingsCmd)

	subjectMappingsCmd.AddCommand(subjectMappingsCreateCmd)

	subjectMappingsCmd.AddCommand(subjectMappingsGetCmd)

	subjectMappingsCmd.AddCommand(subjectMappingsListCmd)
	attributeUpdateCmd.Flags().StringP("resource-selector-name", "n", "", "Resource Selector Name")
	attributeUpdateCmd.Flags().StringSliceVarP(&resourceSelectorLabels, "resource-selector-labels", "l", []string{}, "Resource Selector Labels defined as <key>::<value>")

	subjectMappingsCmd.AddCommand(subjectMappingsDeleteCmd)
}
