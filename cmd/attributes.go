/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

// attributesCmd represents the attributes command
var attributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: "Attributes CRUD operations",
	Long: `Manage your configured attributes [Create, Get one, List all, Update, Delete]
	through use of this CLI.`,
}

var (
	attrValues           []string
	groupBy              []string
	resourceDependencies []string
)

var attributeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an attribute",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cli.ExitWithError("Invalid ID", err)
		}

		close := cli.GrpcConnect(cmd)
		defer close()

		attr, err := handlers.GetAttribute(id)
		if err != nil {
			errMsg := fmt.Sprintf("Could not find attribute (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		fmt.Println(cli.SuccessMessage("Attribute found"))
		fmt.Println(
			cli.NewTabular().
				Rows([][]string{
					{"Id", strconv.Itoa(int(attr.Id))},
					{"Name", attr.Name},
					{"Rule", attr.Rule},
					{"Values", cli.CommaSeparated(attr.Values)},
					{"Namespace", attr.Namespace},
					{"Description", attr.Description},
				}...).Render(),
		)
	},
}

// List attributes
var attributesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attributes",
	Run: func(cmd *cobra.Command, args []string) {
		close := cli.GrpcConnect(cmd)
		defer close()

		attrs, err := handlers.ListAttributes()
		if err != nil {
			cli.ExitWithError("Could not get attributes", err)
		}

		t := cli.NewTable()
		t.Headers("Id", "Namespace", "Name", "Rule", "Values")
		for _, attr := range attrs {
			t.Row(
				strconv.Itoa(int(attr.Id)),
				attr.Namespace,
				attr.Name,
				attr.Rule,
				cli.CommaSeparated(attr.Values),
			)
		}
		fmt.Println(t.Render())
	},
}

// Create an attribute
var attributesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an attribute",
	Run: func(cmd *cobra.Command, args []string) {
		close := cli.GrpcConnect(cmd)
		defer close()

		flagHelper := cli.NewFlagHelper(cmd)
		name := flagHelper.GetRequiredString("name")
		rule := flagHelper.GetRequiredString("rule")
		values := flagHelper.GetStringSlice("values", attrValues, cli.FlagHelperStringSliceOptions{
			Min: 1,
		})
		namespace := flagHelper.GetRequiredString("namespace")
		description := flagHelper.GetRequiredString("description")

		if _, err := handlers.CreateAttribute(name, rule, values, namespace, description); err != nil {
			cli.ExitWithError("Could not create attribute", err)
		}

		fmt.Println(cli.SuccessMessage("Attribute created"))
		fmt.Println(
			cli.NewTabular().Rows([][]string{
				{"Name", name},
				{"Rule", rule},
				{"Values", cli.CommaSeparated(values)},
				{"Namespace", namespace},
				{"Description", description},
			}...).Render(),
		)
	},
}

var attributesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an attribute",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(cli.ErrorMessage("Invalid ID", err))
			os.Exit(1)
		}

		close := cli.GrpcConnect(cmd)
		defer close()

		attr, err := handlers.GetAttribute(id)
		if err != nil {
			errMsg := fmt.Sprintf("Could not find attribute (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		cli.ConfirmDelete("attribute", attr.Fqn)

		if err := handlers.DeleteAttribute(id); err != nil {
			errMsg := fmt.Sprintf("Could not delete attribute (%d)", id)
			cli.ExitIfNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		fmt.Println(cli.SuccessMessage("Attribute deleted"))
		fmt.Println(
			cli.NewTabular().
				Rows([][]string{
					{"Name", attr.Name},
					{"Rule", attr.Rule},
					{"Values", cli.CommaSeparated(attr.Values)},
					{"Namespace", attr.Namespace},
					{"Description", attr.Description},
				}...).Render(),
		)
	},
}

// Update one attribute
var attributeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an attribute",
	Run: func(cmd *cobra.Command, args []string) {
		close := cli.GrpcConnect(cmd)
		defer close()

		flagHelper := cli.NewFlagHelper(cmd)

		id := flagHelper.GetRequiredInt32("id")
		name := flagHelper.GetRequiredString("name")
		rule := flagHelper.GetRequiredString("rule")

		values := flagHelper.GetStringSlice("values", attrValues, cli.FlagHelperStringSliceOptions{
			Min: 1,
		})

		groupBy := flagHelper.GetStringSlice("group-by", groupBy, cli.FlagHelperStringSliceOptions{
			Min: 0,
		})

		resourceDependencies := flagHelper.GetStringSlice("resource-dependencies", resourceDependencies, cli.FlagHelperStringSliceOptions{
			Min: 0,
		})

		resourceId := flagHelper.GetRequiredInt32("resource-id")
		resourceVersion := flagHelper.GetRequiredInt32("resource-version")
		resourceName := flagHelper.GetRequiredString("resource-name")
		resourceNamespace := flagHelper.GetRequiredString("resource-namespace")
		resourceDescription := flagHelper.GetRequiredString("resource-description")

		if _, err := handlers.UpdateAttribute(
			id,
			name,
			rule,
			values,
			groupBy,
			resourceId,
			resourceVersion,
			resourceName,
			resourceNamespace,
			resourceDescription,
			resourceDependencies,
		); err != nil {
			cli.ExitWithError("Could not update attribute", err)
			return
		} else {
			fmt.Println(cli.SuccessMessage(fmt.Sprintf("Attribute id: %d updated.", id)))
		}
	},
}

func init() {
	rootCmd.AddCommand(attributesCmd)

	attributesCmd.AddCommand(attributeGetCmd)

	attributesCmd.AddCommand(attributesListCmd)

	attributesCmd.AddCommand(attributesCreateCmd)
	attributesCreateCmd.Flags().StringP("name", "n", "", "Name of the attribute")
	attributesCreateCmd.Flags().StringP("rule", "r", "", "Rule of the attribute")
	attributesCreateCmd.Flags().StringSliceVarP(&attrValues, "values", "v", []string{}, "Values of the attribute")
	attributesCreateCmd.Flags().StringP("namespace", "s", "", "Namespace of the attribute")
	attributesCreateCmd.Flags().StringP("description", "d", "", "Description of the attribute")

	attributesCmd.AddCommand(attributeUpdateCmd)
	attributeUpdateCmd.Flags().Int32P("id", "i", 0, "Id of the attribute")
	attributeUpdateCmd.Flags().StringP("name", "n", "", "Name of the attribute")
	attributeUpdateCmd.Flags().StringP("rule", "r", "", "Rule of the attribute")
	attributeUpdateCmd.Flags().StringSliceVarP(&attrValues, "values", "v", []string{}, "Values of the attribute")
	attributeUpdateCmd.Flags().StringSliceVarP(&groupBy, "group-by", "g", []string{}, "GroupBy of the attribute")
	attributeUpdateCmd.Flags().StringSliceVarP(&resourceDependencies, "resource-dependencies", "d", []string{}, "ResourceDependencies of the attribute definition descriptor")
	attributeUpdateCmd.Flags().Int32P("resource-id", "I", 0, "ResourceId of the attribute definition descriptor")
	attributeUpdateCmd.Flags().Int32P("resource-version", "V", 0, "ResourceVersion of the attribute definition descriptor")
	attributeUpdateCmd.Flags().StringP("resource-name", "N", "", "ResourceName of the attribute definition descriptor")
	attributeUpdateCmd.Flags().StringP("resource-namespace", "S", "", "ResourceNamespace of the attribute definition descriptor")
	attributeUpdateCmd.Flags().StringP("resource-description", "D", "", "ResourceDescription of the attribute definition descriptor")

	attributesCmd.AddCommand(attributesDeleteCmd)
}
