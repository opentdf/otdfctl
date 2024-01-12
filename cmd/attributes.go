/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

var attributeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an attribute",
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
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				fmt.Println(cli.ErrorMessage(errMsg+" not found", nil))
				os.Exit(1)
			}
			fmt.Println(cli.ErrorMessage(errMsg, err))
			os.Exit(1)
		}

		fmt.Println(cli.SuccessMessage("Attribute found"))
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

// List attributes
var attributesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attributes",
	Run: func(cmd *cobra.Command, args []string) {
		close := cli.GrpcConnect(cmd)
		defer close()

		attrs, err := handlers.ListAttributes()
		if err != nil {
			fmt.Println(cli.ErrorMessage("Could not get attributes", err))
			return
		}

		t := cli.NewTable()
		t.Headers("Namespace", "Name", "Rule", "Values")
		for _, attr := range attrs {
			t.Row(
				attr.Namespace,
				attr.Name,
				attr.Rule,
				cli.CommaSeparated(attr.Values),
			)
		}
		fmt.Print(t.Render())
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
		values := flagHelper.GetRequiredStringSlice("values", attrValues, cli.FlagHelperStringSliceOptions{
			Min: 1,
		})
		namespace := flagHelper.GetRequiredString("namespace")
		description := flagHelper.GetRequiredString("description")

		if _, err := handlers.CreateAttribute(name, rule, values, namespace, description); err != nil {
			fmt.Println(cli.ErrorMessage("Could not create attribute", err))
			os.Exit(1)
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
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				fmt.Println(cli.ErrorMessage(errMsg+" not found", nil))
				os.Exit(1)
			}
			fmt.Println(cli.ErrorMessage(errMsg, err))
			os.Exit(1)
		}

		// prompt for confirmation
		var confirm bool
		err = huh.NewConfirm().
			Title(fmt.Sprintf("Are you sure you want to delete attribute:\n\n\t%s", attr.Fqn)).
			Affirmative("yes").
			Negative("no").
			Value(&confirm).
			Run()
		if err != nil {
			fmt.Println(cli.ErrorMessage("Confirmation prompt failed", err))
			os.Exit(1)
		}

		if !confirm {
			fmt.Println(cli.ErrorMessage("Aborted", nil))
			os.Exit(1)
		}

		if err := handlers.DeleteAttribute(id); err != nil {
			errMsg := fmt.Sprintf("Could not delete attribute (%d)", id)
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				fmt.Println(cli.ErrorMessage(errMsg+" not found", nil))
				os.Exit(1)
			}
			fmt.Println(cli.ErrorMessage(errMsg, err))
			os.Exit(1)
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

	attributesCmd.AddCommand(attributesDeleteCmd)
}
