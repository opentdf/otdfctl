package cmd

import (
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/chat"
)

func init() {
	err := chat.LoadConfig("chat_config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading chat_config: %v\n", err)
		os.Exit(1)
	}
	chat.ConfigureChatCommand()
	RootCmd.AddCommand(chat.GetChatCommand())
}
