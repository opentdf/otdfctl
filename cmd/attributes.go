package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	attrValues []string

	attributeCommands = []string{
		attributesCreateCmd.Use,
		attributeGetCmd.Use,
		attributesListCmd.Use,
		attributeUpdateCmd.Use,
		attributesDeleteCmd.Use,
	}

	attributesCmd = &cobra.Command{
		Use:   "attributes",
		Short: "Manage attributes [" + strings.Join(attributeCommands, ", ") + "]",
		Long: `
Attributes - commands to manage attributes within the platform.

Attributes are used to to define the properties of a piece of data. These attributes will then be
used to define the access controls based on subject encodings and entity entitlements.
`,
	}

	// Create an attribute
	attributesCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			name := flagHelper.GetRequiredString("name")
			rule := flagHelper.GetRequiredString("rule")
			values := flagHelper.GetStringSlice("values", attrValues, cli.FlagHelperStringSliceOptions{})
			namespace := flagHelper.GetRequiredString("namespace")

			if _, err := h.CreateAttribute(name, rule, values, namespace); err != nil {
				cli.ExitWithError("Could not create attribute", err)
			}

			fmt.Println(cli.SuccessMessage("Attribute created"))
			fmt.Println(
				cli.NewTabular().Rows([][]string{
					{"Name", name},
					{"Rule", rule},
					{"Values", cli.CommaSeparated(values)},
					{"Namespace", namespace},
				}...).Render(),
			)
		},
	}

	// Get an attribute
	attributeGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			h := cli.NewHandler(cmd)
			defer h.Close()

			attr, err := h.GetAttribute(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find attribute (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			a := cli.GetSimpleAttribute(attr)
			fmt.Println(cli.SuccessMessage("Attribute found"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Id", a.Id},
						{"Name", a.Name},
						{"Rule", a.Rule},
						{"Values", cli.CommaSeparated(a.Values)},
						{"Namespace", a.Namespace},
					}...).Render(),
			)
		},
	}

	// List attributes
	attributesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List attributes",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			attrs, err := h.ListAttributes()
			if err != nil {
				cli.ExitWithError("Could not get attributes", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Namespace", "Name", "Rule", "Values")
			for _, attr := range attrs {
				a := cli.GetSimpleAttribute(attr)
				t.Row(
					a.Id,
					a.Namespace,
					a.Name,
					a.Rule,
					cli.CommaSeparated(a.Values),
				)
			}
			fmt.Println(t.Render())
		},
	}

	attributesDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			h := cli.NewHandler(cmd)
			defer h.Close()

			attr, err := h.GetAttribute(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not find attribute (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmDelete("attribute", attr.Name)

			attr, err = h.DeleteAttribute(id)
			if err != nil {
				errMsg := fmt.Sprintf("Could not delete attribute (%s)", id)
				cli.ExitWithNotFoundError(errMsg, err)
				cli.ExitWithError(errMsg, err)
			}

			a := cli.GetSimpleAttribute(attr)
			fmt.Println(cli.SuccessMessage("Attribute deleted"))
			fmt.Println(
				cli.NewTabular().
					Rows([][]string{
						{"Name", a.Name},
						{"Rule", a.Rule},
						{"Values", cli.CommaSeparated(a.Values)},
						{"Namespace", a.Namespace},
					}...).Render(),
			)
		},
	}

	// Update one attribute
	attributeUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			if _, err := h.UpdateAttribute(id); err != nil {
				cli.ExitWithError("Could not update attribute", err)
			} else {
				fmt.Println(cli.SuccessMessage(fmt.Sprintf("Attribute id: %s updated.", id)))
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(attributesCmd)

	// Create an attribute
	attributesCmd.AddCommand(attributesCreateCmd)
	attributesCreateCmd.Flags().StringP("name", "n", "", "Name of the attribute")
	attributesCreateCmd.Flags().StringP("rule", "r", "", "Rule of the attribute")
	attributesCreateCmd.Flags().StringSliceVarP(&attrValues, "values", "v", []string{}, "Values of the attribute")
	attributesCreateCmd.Flags().StringP("namespace", "s", "", "Namespace of the attribute")
	attributesCreateCmd.Flags().StringP("description", "d", "", "Description of the attribute")

	// Get an attribute
	attributesCmd.AddCommand(attributeGetCmd)
	attributeGetCmd.Flags().StringP("id", "i", "", "Id of the attribute")

	// List attributes
	attributesCmd.AddCommand(attributesListCmd)

	// Update an attribute
	attributesCmd.AddCommand(attributeUpdateCmd)
	attributeUpdateCmd.Flags().StringP("id", "i", "", "Id of the attribute")

	// Delete an attribute
	attributesCmd.AddCommand(attributesDeleteCmd)
	attributesDeleteCmd.Flags().StringP("id", "i", "", "Id of the attribute")
}
