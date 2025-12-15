//nolint:forbidigo // print statements require flexibility
package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

var ErrPrinterExpectsCommand = fmt.Errorf("printer expects a command")

type Printer struct {
	enabled bool
	json    bool
	debug   bool
}

func newPrinter(cli *Cli) *Printer {
	p := &Printer{
		enabled: true,
		json:    false,
		debug:   false,
	}

	// if json output is enabled, disable the printer
	printJSON := cli.Flags.GetOptionalBool("json")
	p.setJSON(printJSON)

	return p
}

func (p *Printer) setJSON(json bool) {
	p.json = json
	p.enabled = !json
}

// PrintJSON prints the given value as json
// ignores the printer enabled flag
func (c *Cli) printJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		ExitWithError("failed to marshal json", err)
	}
}

func (c *Cli) println(args ...interface{}) {
	if c.printer.enabled {
		fmt.Fprintln(os.Stdout, args...)
	}
}

func (c *Cli) SetJSONOutput(enabled bool) {
	if c.printer == nil {
		return
	}
	c.printer.setJSON(enabled)
}
