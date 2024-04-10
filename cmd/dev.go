package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/opentdf/otdfctl/internal/config"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/spf13/cobra"
)

// devCmd is the command for playground-style development
var devCmd = man.Docs.GetCommand("dev")

func dev_designSystem(cmd *cobra.Command, args []string) {
	fmt.Printf("Design system\n")
	fmt.Printf("=============\n\n")

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
func HandleSuccess(command *cobra.Command, id string, t *table.Table, policyObject interface{}) {
	if OtdfctlCfg.Output.Format == config.OutputJSON || configFlagOverrides.OutputFormatJSON {
		if output, err := json.MarshalIndent(policyObject, "", "  "); err != nil {
			cli.ExitWithError("Error marshalling policy object", err)
		} else {
			fmt.Println(string(output))
		}
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

func init() {
	designCmd := man.Docs.GetCommand("dev/design-system",
		man.WithRun(dev_designSystem),
	)
	devCmd.AddCommand(&designCmd.Command)
	rootCmd.AddCommand(&devCmd.Command)
}
