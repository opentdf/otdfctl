package cli

import (
	// "github.com/evertras/bubble-table/table"
	"github.com/evertras/bubble-table/table"
)

func NewTable(cols ...table.Column) table.Model {
	return table.New(cols).
		BorderRounded().
		WithBaseStyle(styleTable).
		WithNoPagination().
		WithTargetWidth(TermWidth())
}

func NewUUIDColumn() table.Column {
	return table.NewFlexColumn("id", "ID", 5)
}
