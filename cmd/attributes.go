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

var attrValues []string

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
		if rule == "" {
			fmt.Println("Rule is required")
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

// TODO: Update an attribute

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
}
