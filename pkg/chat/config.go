package chat

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Configs are loaded from a otdfctl.yaml file in home directory where defaults are provided
type ChatConfig struct {
	Model      string `yaml:"model" default:"llama3"`
	ApiURL     string `yaml:"apiUrl" default:"http://localhost:11434/api/generate"`
	LogLength  int    `yaml:"logLength" default:"100"`
	Verbose    bool   `yaml:"verbose" default:"true"`
	TokenLimit int    `yaml:"tokenLimit" default:"1000"`
	NumGPU     int    `yaml:"numGPU" default:"1"`
}

type Output struct {
	Format string `yaml:"format"`
}

type Config struct {
	Output Output     `yaml:"output"`
	Chat   ChatConfig `yaml:"chat"`
}

var chatConfig Config
var ask string

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&chatConfig)
	if err != nil {
		return fmt.Errorf("could not decode config YAML: %v", err)
	}

	return nil
}

func ConfigureChatCommand() {
	chatCmd.PersistentFlags().StringVar(&chatConfig.Chat.Model, "model", chatConfig.Chat.Model, "Model name for Ollama")
	chatCmd.PersistentFlags().StringVar(&ask, "ask", "", "Ask a one-off question without entering the chat session") // --ask invocation: ./otdfctl chat --ask "[$question_here]"
}

func RunChatSession(cmd *cobra.Command, args []string) {
	// Call the Setup function before starting the chat session
	// Load in configs from YAML file - change to `otdfctl.yaml` for production, `otdfctl-example.yaml` for development
	configFile := "otdfctl-example.yaml"
	// configFile = "otdfctl.yaml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configFile = "otdfctl-example.yaml"
	}

	err := Setup(configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
		return
	}
	logger, err := NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		return
	}
	defer logger.Close()

	if ask != "" {
		HandleUserInput(ask, logger)
		return
	}

	fmt.Println("Starting chat session. Type 'exit' or 'quit' to end.")
	UserInputLoop(logger)
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with a LLM helper aid",
	Long:  `This command starts an interactive chat session with a local LLM to help with setup, debugging, or generic troubleshooting`,
	Run:   RunChatSession,
}

func GetChatCommand() *cobra.Command {
	return chatCmd
}
