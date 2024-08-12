package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jeremywohl/flatten"
	"github.com/spf13/cobra"
)

func NewTabular(rows ...[]string) table.Model {
	columnKeyProperty := "Property"
	columnKeyValue := "Value"
	t := NewTable(
		table.NewFlexColumn(columnKeyProperty, columnKeyProperty, 1),
		table.NewFlexColumn(columnKeyValue, columnKeyValue, 2),
	)

	tr := []table.Row{}
	if len(rows) == 0 {
		tr = append(tr, table.NewRow(table.RowData{
			columnKeyProperty: "No properties found",
			columnKeyValue:    "",
		}))
	}
	for _, r := range rows {
		p := r[0]
		v := ""
		if len(r) > 1 {
			v = r[1]
		}
		tr = append(tr, table.NewRow(table.RowData{
			columnKeyProperty: p,
			columnKeyValue:    v,
		}))
	}

	t = t.WithTargetWidth(TermWidth())

	t = t.WithRows(tr)
	return t
}

func NewTabularFromStruct(s interface{}) (table.Model, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return table.Model{}, err
	}

	return NewTabularFromJson(b)
}

func NewTabularFromJson(jsonData []byte) (table.Model, error) {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return table.Model{}, err
	}

	v, err := flatten.Flatten(result, "", flatten.DotStyle)
	if err != nil {
		return table.Model{}, err
	}

	var rows [][]string
	for k, v := range v {
		rows = append(rows, []string{k, fmt.Sprintf("%v", v)})
	}

	// Sort rows by key
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	return NewTabular(rows...), nil
}

func getJsonHelper(command string) string {
	return fmt.Sprintf("Use '%s --json' to see all properties", command)
}

func PrintSuccessTable(cmd *cobra.Command, id string, t table.Model) {
	parent := cmd.Parent()
	resourceShort := parent.Use
	resource := parent.Use
	for parent.Parent() != nil {
		resource = parent.Parent().Use + " " + resource
		parent = parent.Parent()
	}

	var msg struct {
		verb   string
		helper string
	}
	switch cmd.Use {
	case ActionGet:
		msg.verb = fmt.Sprintf("Found %s: %s", resourceShort, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case ActionCreate:
		msg.verb = fmt.Sprintf("Created %s: %s", resourceShort, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case ActionUpdate:
		msg.verb = fmt.Sprintf("Updated %s: %s", resourceShort, id)
		msg.helper = getJsonHelper(resource + " get --id=" + id)
	case ActionDelete:
		msg.verb = fmt.Sprintf("Deleted %s: %s", resourceShort, id)
		// strip off unsafe subcommand if found to get proper path to the list command
		msg.helper = getJsonHelper(strings.ReplaceAll(resource, " unsafe", "") + " list")
	case ActionDeactivate:
		msg.verb = fmt.Sprintf("Deactivated %s: %s", resourceShort, id)
		msg.helper = getJsonHelper(resource + " list") // TODO: make sure the filters are provided here to get ACTIVE/INACTIVE/ANY
	case ActionList:
		msg.verb = fmt.Sprintf("Found %s list", resourceShort)
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
