package cli

import (
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/spf13/cobra"
)

func NewHandler(cmd *cobra.Command) handlers.Handler {
	platformEndpoint := cmd.Flag("host").Value.String()
	h, err := handlers.New(platformEndpoint)
	if err != nil {
		ExitWithError("Failed to connect to server", err)
	}
	return h
}
