package cli

import (
	"errors"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ExitWithError(errMsg string, err error) {
	// This is temporary until we can refactor the code to use the Cli struct
	(&Cli{printer: &Printer{enabled: true}}).ExitWithError(errMsg, err)
}

func ExitWithNotFoundError(errMsg string, err error) {
	// This is temporary until we can refactor the code to use the Cli struct
	(&Cli{printer: &Printer{enabled: true}}).ExitWithNotFoundError(errMsg, err)
}

func ExitWithWarning(warnMsg string) {
	// This is temporary until we can refactor the code to use the Cli struct
	(&Cli{printer: &Printer{enabled: true}}).ExitWithWarning(warnMsg)
}

// ExitWithError prints an error message and exits with a non-zero status code.
func (c *Cli) ExitWithError(errMsg string, err error) {
	if err == nil {
		// ensure message and exit are satisfied even if nil error passed
		err = errors.New("")
	}
	c.ExitWithNotFoundError(errMsg, err)
	if err != nil {
		c.Println(ErrorMessage(errMsg, err))
		os.Exit(1)
	}
}

// ExitWithNotFoundError prints an error message and exits with a non-zero status code if the error is a NotFound error.
func (c *Cli) ExitWithNotFoundError(errMsg string, err error) {
	if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
		c.Println(ErrorMessage(errMsg+": not found", nil))
		os.Exit(1)
	}
}

func (c *Cli) ExitWithMessage(msg string, code int) {
	c.Println(msg)
	os.Exit(code)
}

func (c *Cli) ExitWithWarning(warnMsg string) {
	c.ExitWithMessage(WarningMessage(warnMsg), 0)
}

func (c *Cli) ExitWithSuccess(msg string) {
	c.ExitWithMessage(SuccessMessage(msg), 0)
}

func (c *Cli) ExitWithJson(v interface{}) {
	if c.printer.json {
		c.PrintJson(v)
		os.Exit(0)
	}
}
