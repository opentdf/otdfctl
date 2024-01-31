package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opentdf/opentdf-v2-poc/sdk/common"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development tools",
}

var designCmd = &cobra.Command{
	Use:   "design-system",
	Short: "Show design system",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Design system\n")
		fmt.Printf("=============\n\n")

		printDSComponent("Table", renderDSTable())

		printDSComponent("Messages", renderDSMessages())
	},
}

func printDSComponent(title string, component string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("-----\n\n")
	fmt.Printf("%s\n", component)
	fmt.Printf("\n\n")
}

func renderDSTable() string {
	tbl := cli.NewTable()
	tbl.Headers("One", "Two", "Three")
	tbl.Row("1", "2", "3")
	tbl.Row("4", "5", "6")
	return tbl.Render()
}

func renderDSMessages() string {
	return cli.SuccessMessage("Success message") + "\n" + cli.ErrorMessage("Error message", nil)
}

func getMetadataRows(m *common.Metadata) [][]string {
	if m != nil {
		metadataRows := [][]string{
			{"Created At", m.CreatedAt.String()},
			{"Updated At", m.UpdatedAt.String()},
		}
		if m.Labels != nil {
			labelRows := []string{}
			for k, v := range m.Labels {
				labelRows = append(labelRows, fmt.Sprintf("%s: %s", k, v))
			}
			metadataRows = append(metadataRows, []string{"Labels", cli.CommaSeparated(labelRows)})
		}
		if m.Description != "" {
			metadataRows = append(metadataRows, []string{"Description", m.Description})
		}
		return metadataRows
	}
	return nil
}

func unMarshalMetadata(m string) *common.MetadataMutable {
	if m != "" {
		metadata := &common.MetadataMutable{}
		if err := json.Unmarshal([]byte(m), metadata); err != nil {
			cli.ExitWithError("Could not unmarshal metadata", err)
		}
		return metadata
	}
	return nil
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(designCmd)
}
