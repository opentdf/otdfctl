/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/opentdf/tructl/pkg/grpc"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

// attributesCmd represents the attributes command
var attributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var (
	attrValues           []string
	groupBy              []string
	resourceDependencies []string
)

// List attributes
var attributesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attributes",
	Run: func(cmd *cobra.Command, args []string) {
		if err := grpc.Connect(cmd.Flag("host").Value.String()); err != nil {
			fmt.Println(err)
			return
		}
		defer grpc.Conn.Close()

		resp, err := handlers.ListAttributes()
		if err != nil {
			fmt.Println(err)
			return
		}

		columns := []table.Column{
			{Title: "Namespace", Width: 20},
			{Title: "Name", Width: 20},
			{Title: "Rule", Width: 20},
			{Title: "Values", Width: 20},
		}

		rows := []table.Row{}
		for _, attr := range resp.Definitions {
			values := ""
			for i, v := range attr.Values {
				if i != 0 {
					values += ", "
				}
				values += v.Value
			}

			rows = append(rows, table.Row{
				attr.Descriptor_.Namespace,
				attr.Name,
				handlers.GetAttributeRuleFromAttributeType(attr.Rule),
				values,
			})
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithHeight(7),
		)

		t.Update("")
		fmt.Print(t.View())
	},
}

// Create an attribute
var attributesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an attribute",
	Run: func(cmd *cobra.Command, args []string) {
		if err := grpc.Connect(cmd.Flag("host").Value.String()); err != nil {
			fmt.Println(err)
			return
		}
		defer grpc.Conn.Close()

		name := cmd.Flag("name").Value.String()
		if name == "" {
			fmt.Println("Name is required")
			return
		}

		rule := cmd.Flag("rule").Value.String()
		if rule == "" || handlers.GetAttributeRuleFromReadableString(rule) == 0 {
			fmt.Printf("Flag 'rule' is required and must be one of: %v", handlers.GetAttributeRuleOptions())
			return
		}

		if len(attrValues) == 0 {
			fmt.Println("Values is required")
			return
		}

		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			fmt.Println("Namespace is required")
			return
		}

		description := cmd.Flag("description").Value.String()
		if description == "" {
			fmt.Println("Description is required")
			return
		}

		if resp, err := handlers.CreateAttribute(name, rule, attrValues, namespace, description); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(resp)
		}
	},
}

// TODO: think about how to improve this. Passing a 12 flags/args in a CLI is very challenging...
// Update one attribute
var attributeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an attribute",
	Run: func(cmd *cobra.Command, args []string) {
		if err := grpc.Connect(cmd.Flag("host").Value.String()); err != nil {
			fmt.Println(err)
			return
		}
		defer grpc.Conn.Close()

		Id, e := cmd.Flags().GetInt32("id")
		if e != nil {
			fmt.Println("Flag 'id' is required")
			return
		}

		name := cmd.Flag("name").Value.String()
		if name == "" {
			fmt.Println("Flag 'name' is required")
			return
		}

		rule := cmd.Flag("rule").Value.String()
		if rule == "" || handlers.GetAttributeRuleFromReadableString(rule) == 0 {
			fmt.Printf("Flag 'rule' is required and must be one of: %v", handlers.GetAttributeRuleOptions())
			return
		}

		// Would this memory leak since the var is in scope to create & update both?
		// TODO: check the same for groupBy, dependencies
		if len(attrValues) == 0 {
			fmt.Println("Flag 'values' is required")
			return
		}

		if len(groupBy) == 0 {
			fmt.Println("Flag 'group-by' is required")
			return
		}

		if len(resourceDependencies) == 0 {
			fmt.Println("Flag 'resource-dependencies' is required")
			return
		}

		// TODO: are all of these required, or can some be defaulted / looked up?
		resourceId, e := cmd.Flags().GetInt32("resource-id")
		if e != nil {
			fmt.Println("Flag 'resource-id' is required")
			return
		}

		resourceVersion, e := cmd.Flags().GetInt32("res-version")
		if e != nil {
			fmt.Println("Flag 'resource-version' is required")
			return
		}

		resourceName := cmd.Flag("resource-name").Value.String()
		if resourceName == "" {
			fmt.Println("Flag 'resource-name' is required")
			return
		}

		resourceNamespace := cmd.Flag("resource-namespace").Value.String()
		if resourceNamespace == "" {
			fmt.Println("Flag 'resource-namespace' is required")
			return
		}

		resourceFqn := cmd.Flag("resource-fqn").Value.String()
		if resourceFqn == "" {
			fmt.Println("Flag 'resource-fqn' is required")
			return
		}

		resourceDescription := cmd.Flag("resource-description").Value.String()
		if resourceDescription == "" {
			fmt.Println("Flag 'resource-description' is required")
			return
		}

		if resp, err := handlers.UpdateAttribute(
			Id,
			name,
			rule,
			attrValues,
			groupBy,
			resourceId,
			resourceVersion,
			resourceName,
			resourceNamespace,
			resourceFqn,
			resourceDescription,
			resourceDependencies,
		); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(resp)
		}
	},
}

// TODO: Delete an attribute

func init() {
	rootCmd.AddCommand(attributesCmd)

	attributesCmd.AddCommand(attributesListCmd)

	attributesCmd.AddCommand(attributesCreateCmd)
	attributesCreateCmd.Flags().StringP("name", "n", "", "Name of the attribute")
	attributesCreateCmd.Flags().StringP("rule", "r", "", "Rule of the attribute")
	attributesCreateCmd.Flags().StringSliceVarP(&attrValues, "values", "v", []string{}, "Values of the attribute")
	attributesCreateCmd.Flags().StringP("namespace", "s", "", "Namespace of the attribute")
	attributesCreateCmd.Flags().StringP("description", "d", "", "Description of the attribute")

	attributesCmd.AddCommand(attributeUpdateCmd)
	// NOTE: I can't find the ID of created/listed attributes anywhere in grpc responses? Where is this located?
	attributeUpdateCmd.Flags().Int32P("id", "i", 0, "Id of the attribute")
	attributeUpdateCmd.Flags().StringP("name", "n", "", "Name of the attribute")
	attributeUpdateCmd.Flags().StringP("rule", "r", "", "Rule of the attribute")
	attributeUpdateCmd.Flags().StringSliceVarP(&attrValues, "values", "v", []string{}, "Values of the attribute")
	attributeUpdateCmd.Flags().StringSliceVarP(&groupBy, "group-by", "g", []string{}, "GroupBy of the attribute")
	// TODO: again, can any of these be defaulted/inferred via lookup?
	attributeUpdateCmd.Flags().StringSliceVarP(&resourceDependencies, "resource-dependencies", "d", []string{}, "ResourceDependencies of the attribute definition descriptor")
	attributeUpdateCmd.Flags().Int32P("resource-id", "I", 0, "ResourceId of the attribute definition descriptor")
	attributeUpdateCmd.Flags().Int32P("resource-version", "V", 0, "ResourceVersion of the attribute definition descriptor")
	attributeUpdateCmd.Flags().StringP("resource-name", "N", "", "ResourceName of the attribute definition descriptor")
	attributeUpdateCmd.Flags().StringP("resource-namespace", "S", "", "ResourceNamespace of the attribute definition descriptor")
	attributeUpdateCmd.Flags().StringP("resource-fqn", "F", "", "ResourceFqn of the attribute")
	attributeUpdateCmd.Flags().StringP("resource-description", "D", "", "ResourceDescription of the attribute definition descriptor")
}
