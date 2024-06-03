package cli

import (
	"fmt"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ExitWithError(errMsg string, err error) {
	ExitWithNotFoundError(errMsg, err)
	if err != nil {
		fmt.Println(ErrorMessage(errMsg, err))
		os.Exit(1)
	}
}

func ExitWithNotFoundError(errMsg string, err error) {
	if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
		fmt.Println(ErrorMessage(errMsg+": not found", nil))
		os.Exit(1)
	}
}

func ExitWithWarning(warnMsg string) {
	fmt.Println(WarningMessage(warnMsg))
	os.Exit(0)
}
