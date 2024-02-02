package cli

import (
	"fmt"
	"strings"

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
	width := defaultTableWidth
	resource := cmd.Parent().Use

	var msg struct {
		verb   string
		helper string
	}
	switch cmd.Use {
	case "get":
		msg.verb = fmt.Sprintf("Found %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case "create":
		msg.verb = fmt.Sprintf("Created %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case "update":
		msg.verb = fmt.Sprintf("Updated %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case "delete":
		msg.verb = fmt.Sprintf("Deleted %s: %s", resource, id)
		msg.helper = getJsonHelper(resource + " list")
	case "list":
		msg.verb = fmt.Sprintf("Found %s list", resource)
		msg.helper = getJsonHelper(resource + " get --id=<id>")
	default:
		msg.verb = ""
		msg.helper = ""
	}

	successMessage := SuccessMessage(msg.verb) + "\n"

	padding := (width - len(msg.helper) - 6)
	jsonDirections := "\n" + strings.Repeat("/", width) + "\n" +
		"///" + strings.Repeat(" ", (padding/2)) + msg.helper + strings.Repeat(" ", (padding/2)+(padding%2)) + "///" + "\n" +
		strings.Repeat("/", width)

	fmt.Println(successMessage + "\n" + t.Render() + "\n" + jsonDirections + "\n")
}
