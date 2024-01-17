package cmd

import (
	"fmt"
	"strconv"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

// acse represents the Access Control Subject Encodings command
var subjectEncodingsCmd = &cobra.Command{
	Use:   "subject encodings",
	Short: "Access Control Subject Encodings CRUD operations",
	Long: `Manage your configured Subject Encoding Mappings [Create, Get one, List all, Update, Delete]
	through use of this CLI.`,
}

// Get one Access Control Subject Encoding
var subjectEncodingsMappingGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an Access Control Subject Encoding",
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
			cli.ExitWithNotFoundError(errMsg, err)
			cli.ExitWithError(errMsg, err)
		}

		fmt.Println(cli.SuccessMessage("Access Control Subject Encoding Mapping found"))
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

// List all access control subject encodings
var subjectEncodingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list Access Control Subject Encodings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		close := cli.GrpcConnect(cmd)
		defer close()

		// TODO: selector?

		mappings, err := handlers.ListSubjectMappings()
		if err != nil {
			cli.ExitWithError("Could not list subject encodings", err)
		}

		t := cli.NewTable()
		// TODO: accurate table columns
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

func init() {
	rootCmd.AddCommand(subjectEncodingsCmd)

	subjectEncodingsCmd.AddCommand(subjectEncodingsMappingGetCmd)
	subjectEncodingsCmd.AddCommand(subjectEncodingsListCmd)
}
