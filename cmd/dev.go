package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/spf13/cobra"
)

// devCmd is the command for playground-style development
var devCmd = man.Docs.GetCommand("dev")

func dev_designSystem(cmd *cobra.Command, args []string) {
	fmt.Printf("Design system\n=============\n\n")

	printDSComponent("Table", renderDSTable())

	printDSComponent("Messages", renderDSMessages())
}

func printDSComponent(title string, component string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("-----\n\n")
	fmt.Printf("%s\n", component)
	fmt.Printf("\n\n")
}

func renderDSTable() string {
	tbl := cli.NewTable(
		table.NewFlexColumn("one", "One", 1),
		table.NewFlexColumn("two", "Two", 1),
		table.NewFlexColumn("three", "Three", 1),
	).WithRows([]table.Row{
		table.NewRow(table.RowData{
			"one":   "1",
			"two":   "2",
			"three": "3",
		}),
		table.NewRow(table.RowData{
			"one":   "4",
			"two":   "5",
			"three": "6",
		}),
	})
	return tbl.View()
}

func renderDSMessages() string {
	return cli.SuccessMessage("Success message") + "\n" + cli.ErrorMessage("Error message", nil)
}

func getMetadataRows(m *common.Metadata) [][]string {
	if m != nil {
		metadata := cli.ConstructMetadata(m)
		metadataRows := [][]string{
			{"Created At", metadata["Created At"]},
			{"Updated At", metadata["Updated At"]},
		}
		if m.Labels != nil {
			metadataRows = append(metadataRows, []string{"Labels", metadata["Labels"]})
		}
		return metadataRows
	}
	return nil
}

func unMarshalMetadata(m string) *common.MetadataMutable {
	if m != "" {
		metadata := &common.MetadataMutable{}
		if err := json.Unmarshal([]byte(m), metadata); err != nil {
			cli.ExitWithError("Failed to unmarshal metadata", err)
		}
		return metadata
	}
	return nil
}

func getMetadataMutable(labels []string) *common.MetadataMutable {
	metadata := common.MetadataMutable{}
	if len(labels) > 0 {
		metadata.Labels = map[string]string{}
		for _, label := range labels {
			kv := strings.Split(label, "=")
			if len(kv) != 2 {
				cli.ExitWithError("Invalid label format", nil)
			}
			metadata.Labels[kv[0]] = kv[1]
		}
		return &metadata
	}
	return nil
}

func getMetadataUpdateBehavior() common.MetadataUpdateEnum {
	if forceReplaceMetadataLabels {
		return common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_REPLACE
	}
	return common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_EXTEND
}

// HandleSuccess prints a success message according to the configured format (styled table or JSON)
func HandleSuccess(command *cobra.Command, id string, t table.Model, policyObject interface{}) {
	c := cli.New(command, []string{})
	if OtdfctlCfg.Output.Format == config.OutputJSON || configFlagOverrides.OutputFormatJSON {
		c.PrintJson(policyObject)
		return
	}
	cli.PrintSuccessTable(command, id, t)
}

// Adds reusable create/update label flags to a Policy command and the optional force-replace-labels flag for updates only
func injectLabelFlags(cmd *cobra.Command, isUpdate bool) {
	cmd.Flags().StringSliceVarP(&metadataLabels, "label", "l", []string{}, "Optional metadata 'labels' in the format: key=value")
	if isUpdate {
		cmd.Flags().BoolVar(&forceReplaceMetadataLabels, "force-replace-labels", false, "Destructively replace entire set of existing metadata 'labels' with any provided to this command.")
	}
}

// Read bytes from stdin without blocking by checking size first
func readPipedStdin() []byte {
	stat, err := os.Stdin.Stat()
	if err != nil {
		cli.ExitWithError("Failed to read stat from stdin", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			cli.ExitWithError("failed to scan bytes from stdin", err)
		}
		return buf
	}
	return nil
}

func readBytesFromFile(filePath string) []byte {
	fileToEncrypt, err := os.Open(filePath)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to open file at path: %s", filePath), err)
	}
	defer fileToEncrypt.Close()

	bytes, err := io.ReadAll(fileToEncrypt)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", filePath), err)
	}
	return bytes
}

func init() {
	designCmd := man.Docs.GetCommand("dev/design-system",
		man.WithRun(dev_designSystem),
	)
	devCmd.AddCommand(&designCmd.Command)
	RootCmd.AddCommand(&devCmd.Command)
}
