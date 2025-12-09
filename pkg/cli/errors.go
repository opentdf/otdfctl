package cli

import (
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ExitCodeSuccess = 0
	ExitCodeError   = 1
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
	c.ExitWithNotFoundError(errMsg, err)
	c.ExitWithMessage(ErrorMessage(errMsg, err), ExitCodeError)
}

// ExitWithNotFoundError prints an error message and exits with a non-zero status code if the error is a NotFound error.
func (c *Cli) ExitWithNotFoundError(errMsg string, err error) {
	if err != nil {
		if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
			c.ExitWithMessage(ErrorMessage(errMsg+": not found", nil), ExitCodeError)
		}
	}
}

func (c *Cli) ExitWithMessage(msg string, code int) {
	c.println(msg)
	os.Exit(code)
}

func (c *Cli) ExitWithWarning(warnMsg string) {
	c.ExitWithMessage(WarningMessage(warnMsg), ExitCodeError)
}

func (c *Cli) ExitWithSuccess(msg string) {
	c.ExitWithMessage(SuccessMessage(msg), ExitCodeSuccess)
}

func (c *Cli) ExitWithStyled(msg string) {
	if c.printer.enabled {
		c.println(msg)
		os.Exit(ExitCodeSuccess)
	}
}

func (c *Cli) ExitWithJSON(v interface{}) {
	if c.printer.json {
		c.printJSON(v)
		os.Exit(ExitCodeSuccess)
	}
}
