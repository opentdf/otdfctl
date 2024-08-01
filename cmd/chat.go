package cmd

import (
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/chat"
)

func init() {
	// Load in configs from YAML file - change to `otdfctl.yaml` for production
	err := chat.LoadConfig("otdfctl-example.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading chat_config: %v\n", err)
		os.Exit(1)
	}
	chat.ConfigureChatCommand()
	RootCmd.AddCommand(chat.GetChatCommand())
}
