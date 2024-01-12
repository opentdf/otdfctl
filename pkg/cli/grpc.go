package cli

import (
	"fmt"
	"os"

	"github.com/opentdf/tructl/pkg/grpc"
	"github.com/spf13/cobra"
)

func GrpcConnect(cmd *cobra.Command) func() {
	if err := grpc.Connect(cmd.Flag("host").Value.String()); err != nil {
		fmt.Println(ErrorMessage("Could not connect to server", err))
		os.Exit(1)
	}
	return func() {
		grpc.Conn.Close()
	}
}
