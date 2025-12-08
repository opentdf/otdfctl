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

	isDebug := cli.Flags.GetOptionalBool("debug")
	p.setDebug(isDebug)

	return p
}

func (p *Printer) setJSON(json bool) {
	p.json = json
	p.enabled = !json
}

func (p *Printer) setDebug(debug bool) {
	p.debug = debug
}

const debugPrefix = "DEBUG: "

func (c *Cli) Debug(args ...interface{}) {
	if c.printer.debug {
		args = append([]interface{}{debugPrefix}, args...)
	}
}

func (c *Cli) Debugf(format string, args ...interface{}) {
	format = "DEBUG: " + format
	if c.printer.debug {
	}
}

func (c *Cli) Debugln(args ...interface{}) {
	if c.printer.debug {
		args = append([]interface{}{"DEBUG: "}, args...)
	}
}

// PrintJSON prints the given value as json
// ignores the printer enabled flag
func (c *Cli) printJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		ExitWithError("failed to marshal json", err)
	}
	fmt.Fprintln(os.Stdout, string(b))
}

func (c *Cli) println(args ...interface{}) {
	if c.printer.enabled {
		fmt.Fprintln(os.Stdout, args...)
	}
}
