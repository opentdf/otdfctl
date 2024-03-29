package cli

import (
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
)

func NewHandler(cmd *cobra.Command) handlers.Handler {
	h, err := handlers.New(cmd.Flag("host").Value.String())
	if err != nil {
		ExitWithError("Failed to connect to server", err)
	}
	return h
}
