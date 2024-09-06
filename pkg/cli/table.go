package cli

import (
	"github.com/evertras/bubble-table/table"
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
