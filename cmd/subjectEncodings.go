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
		var attrRefLabels map[string]string
		h := cli.NewHandler(cmd)
		defer h.Close()

		flagHelper := cli.NewFlagHelper(cmd)

		name := flagHelper.GetRequiredString("name")
		description := flagHelper.GetOptionalString("description")
		resourceDeps := flagHelper.GetStringSlice("resource-dependencies", resourceDependencies, cli.FlagHelperListOptions{Min: 0})

		operator := flagHelper.GetRequiredString("operator")
		subjectAttr := flagHelper.GetRequiredString("subject-attribute")
		subjectValues := flagHelper.GetStringSlice("subject-values", subjectValues, cli.FlagHelperListOptions{Min: 1})

		attrRefName := flagHelper.GetOptionalString("attribute-ref-name")
		if attrRefName == "" {
			if len(attrValueLabels) == 0 { // optional, we'll ignore for now
				cli.ExitWithError("Either attribute-ref-name or attribute-ref-labels must be specified", nil)
			}
			attrRefLabels = flagHelper.GetKeyValuesMap("attribute-ref-labels", attrValueLabels, cli.FlagHelperListOptions{Min: 1})
		}

		if err := h.CreateSubjectMapping(
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

		h := cli.NewHandler(cmd)
		defer h.Close()

		mapping, err := h.GetSubjectMapping(id)
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
		h := cli.NewHandler(cmd)
		defer h.Close()

		var (
			selectorName   string
			selectorLabels map[string]string
		)

		flagHelper := cli.NewFlagHelper(cmd)

		selectorName = flagHelper.GetOptionalString("resource-selector-name")
		if selectorName == "" {
			if len(resourceSelectorLabels) == 0 {
				cli.ExitWithError("Either resource-selector-name or resource-selector-labels must be specified", nil)
			}
			selectorLabels = flagHelper.GetKeyValuesMap("resource-selector-labels", resourceSelectorLabels, cli.FlagHelperListOptions{Min: 1})
		}

		mappings, err := h.ListSubjectMappings(selectorName, selectorLabels)
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

		h := cli.NewHandler(cmd)
		defer h.Close()

		mapping, err := h.GetSubjectMapping(id)
		if err != nil {
			errMsg := fmt.Sprintf("Could not find subject mapping (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		cli.ConfirmDelete("subject mapping", mapping.Name)

		if err := h.DeleteSubjectMapping(id); err != nil {
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
