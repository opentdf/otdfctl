package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

func NewTabular() *table.Table {
	t := NewTable()
	t.Headers("Property", "Value")
	return t
}

func getJsonHelper(command string) string {
	return fmt.Sprintf("Use '%s --json' to see all properties", command)
}

func PrintSuccessTable(cmd *cobra.Command, id string, t *table.Table) {
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

	if t == nil {
		fmt.Println(lipgloss.JoinVertical(lipgloss.Top, successMessage, jsonDirections))
		return
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Top, successMessage, t.Render(), jsonDirections))
}
