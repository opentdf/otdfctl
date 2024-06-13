package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
)

func NewTabular(rows ...[]string) table.Model {
	t := NewTable(
		table.NewColumn("property", "Property", 15),
		table.NewColumn("value", "Value", 15),
	)

	tr := []table.Row{}
	if len(rows) == 0 {
		tr = append(tr, table.NewRow(table.RowData{
			"property": "No properties found",
			"value":    "",
		}))
	}
	for _, r := range rows {
		p := r[0]
		v := ""
		if len(r) > 1 {
			v = r[1]
		}
		tr = append(tr, table.NewRow(table.RowData{
			"property": p,
			"value":    v,
		}))
	}

	return t
}

func getJsonHelper(command string) string {
	return fmt.Sprintf("Use '%s --json' to see all properties", command)
}

func PrintSuccessTable(cmd *cobra.Command, id string, t table.Model) {
	resource := cmd.Parent().Use

	var msg struct {
		verb   string
		helper string
	}
	switch cmd.Use {
	case ActionGet:
		msg.verb = fmt.Sprintf("Found %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case ActionCreate:
		msg.verb = fmt.Sprintf("Created %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case ActionUpdate:
		msg.verb = fmt.Sprintf("Updated %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case ActionDelete:
		msg.verb = fmt.Sprintf("Deleted %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " list")
	case ActionDeactivate:
		msg.verb = fmt.Sprintf("Deactivated %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " list") // TODO: make sure the filters are provided here to get ACTIVE/INACTIVE/ANY
	case ActionList:
		msg.verb = fmt.Sprintf("Found %s list", resource)
		msg.helper = getJsonHelper(resource + " get --id=<id>")
	default:
		msg.verb = ""
		msg.helper = ""
	}

	successMessage := SuccessMessage(msg.verb)
	jsonDirections := FooterMessage(msg.helper)

	ts := t.View()
	if ts == "" {
		fmt.Println(lipgloss.JoinVertical(lipgloss.Top, successMessage, jsonDirections))
		return
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Top, successMessage, ts, jsonDirections))
}
