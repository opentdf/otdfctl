package cmd

import (
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/chat"
)

func init() {
	// Load in configs from YAML file - change to `otdfctl.yaml` for production, `otdfctl-example.yaml` for development
	err := chat.LoadConfig("otdfctl-example.yaml")
	// err := chat.LoadConfig("otdfctl.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading chat config: %v\n", err)
		os.Exit(1)
	}
	chat.ConfigureChatCommand()
	RootCmd.AddCommand(chat.GetChatCommand())
}
