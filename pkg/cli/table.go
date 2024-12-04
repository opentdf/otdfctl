package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/platform/protocol/go/policy"
)

const (
	FlexColumnWidthOne   = 1
	FlexColumnWidthTwo   = 2
	FlexColumnWidthThree = 3
	FlexColumnWidthFour  = 4
	FlexColumnWidthFive  = 5
)

func NewTable(cols ...table.Column) table.Model {
	return table.New(cols).
		BorderRounded().
		WithBaseStyle(styleTable).
		WithNoPagination().
		WithTargetWidth(TermWidth())
}

func NewUUIDColumn() table.Column {
	return table.NewFlexColumn("id", "ID", FlexColumnWidthFive)
}

// Adds the page information to the table footer
func WithListPaginationFooter(t table.Model, p *policy.PageResponse) table.Model {
	info := []string{
		fmt.Sprintf("Total: %d", p.GetTotal()),
		fmt.Sprintf("Current Offset: %d", p.GetCurrentOffset()),
	}
	if p.GetNextOffset() > 0 {
		info = append(info, fmt.Sprintf("Next Offset: %d", p.GetNextOffset()))
	}

	columns := []table.Column{}
	for _, c := range info {
		columns = append(columns, table.NewFlexColumn(c, c, FlexColumnWidthOne))
	}

	content := strings.Join(info, "  |  ")
	// content := table.New(columns).
	// 	WithBaseStyle(styleTable).
	// 	BorderRounded().
	// 	WithTargetWidth(len(info) * 20).
	// 	WithNoPagination().View()

	leftAligned := lipgloss.NewStyle().Align(lipgloss.Left)
	return t.WithStaticFooter(content).WithBaseStyle(leftAligned)
}

// func WithListPaginationFooter(t table.Model, p *policy.PageResponse) table.Model {
// 	info := []string{
// 		fmt.Sprintf("Total: %d", p.GetTotal()),
// 		fmt.Sprintf("Current Offset: %d", p.GetCurrentOffset()),
// 	}
// 	if p.GetNextOffset() > 0 {
// 		info = append(info, fmt.Sprintf("Next Offset: %d", p.GetNextOffset()))
// 	}

// 	columns := []table.Column{}
// 	for _, c := range info {
// 		columns = append(columns, table.NewFlexColumn(c, c, FlexColumnWidthOne))
// 	}

// 	// content := strings.Join(info, "  |  ")
// 	content := table.New(columns).
// 		WithBaseStyle(styleTable).
// 		BorderRounded().
// 		WithTargetWidth(len(info) * 20).
// 		WithNoPagination()

// 	// fmt.Println(lipgloss.JoinVertical(lipgloss.Left, t.View(), content.View()))
// 	return content
// }
