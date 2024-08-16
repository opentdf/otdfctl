package cli

import "fmt"

type Printer struct {
	enabled bool
}

func NewPrinter(enabled bool) *Printer {
	return &Printer{
		enabled: enabled,
	}
}

func (p *Printer) Print(args ...interface{}) {
	if p.enabled {
		fmt.Print(args...)
	}
}

func (p *Printer) Printf(format string, args ...interface{}) {
	if p.enabled {
		fmt.Printf(format, args...)
	}
}

func (p *Printer) Println(args ...interface{}) {
	if p.enabled {
		fmt.Println(args...)
	}
}
