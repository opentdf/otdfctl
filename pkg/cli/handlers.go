package cli

import (
	"errors"
	"fmt"

	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/spf13/cobra"
)

func NewHandler(cmd *cobra.Command) handlers.Handler {
	platformEndpoint := cmd.Flag("host").Value.String()
	h, err := handlers.New(platformEndpoint)
	if err != nil {
		if errors.Is(err, handlers.ErrUnauthenticated) {
			ExitWithError(fmt.Sprintf("Not logged in. Please authenticate via CLI auth flow(s) before using command (%s)", cmd.Use), err)
		}
		ExitWithError("Failed to connect to server", err)
	}
	return h
}
