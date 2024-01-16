package cmd

import (
	"fmt"

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

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(designCmd)
}
