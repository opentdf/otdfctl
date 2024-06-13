package cli

import (

	// "github.com/evertras/bubble-table/table"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/term"
)

var defaultTableWidth int

func NewTable(cols ...table.Column) table.Model {
	return table.New(cols).
		BorderRounded().
		WithBaseStyle(styleTable).
		WithNoPagination().
		WithTargetWidth(defaultTableWidth)
}

func NewUUIDColumn() table.Column {
	return table.NewColumn("id", "ID", 37)
}

func init() {
	// dynamically set the default table width based on terminal size breakpoints
	w, _, err := term.GetSize(0)
	if err != nil {
		w = 80
	}
	defaultTableWidth = w
}
