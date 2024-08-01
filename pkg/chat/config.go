package chat

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Model     string `json:"model"`
	Verbosity string `json:"verbosity"`
	ApiURL    string `json:"apiURL"`
}

var chatConfig Config

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&chatConfig)
	if err != nil {
		return fmt.Errorf("could not decode config JSON: %v", err)
	}

	return nil
}

func ConfigureChatCommand() {
	// TODO: Make more configurable without losing dynamic selection, keeping it accessible via command line flag.
	chatCmd.PersistentFlags().StringVar(&chatConfig.Model, "model", chatConfig.Model, "Model name for Ollama")
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with a LLM helper aid",
	Long:  `This command starts an interactive chat session with a local LLM to help with setup, debugging, or generic troubleshooting`,
	Run:   runChatSession,
}

// TODO: add a 'one-off' --ask flag to allow for a single question to be asked and answered, DRYing existing chat code

func runChatSession(cmd *cobra.Command, args []string) {
	fmt.Println("Starting chat session. Type 'exit' to end.")
	userInputLoop()
}

func GetChatCommand() *cobra.Command {
	return chatCmd
}

// func init() {
// 	err := LoadConfig("chat_config.json")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error loading chat_config: %v\n", err)
// 		os.Exit(1)
// 	}
// 	configureChatCommand()
// }
