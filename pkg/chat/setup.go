package chat

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

// CheckConfigFile checks if the YAML configuration file is present and valid.
func CheckConfigFile(filename string) error {
	fmt.Print("Checking YAML configuration file... ")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("could not decode config YAML: %v", err)
	}

	fmt.Println("PASSED")
	return nil
}

// CheckModelRunning tests if the local model is running
func CheckModelRunning(apiURL string) error {
	fmt.Print("Testing if local model is running... ")
	//remove the /api/generate from the URL to check if the model is running
	apiURL = apiURL[:len(apiURL)-12]
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("FAILED")
		return fmt.Errorf("local model is not running or reachable at %s", apiURL)
	}
	fmt.Println("PASSED")
	return nil
}

// DisplayConfigSettings prints the current configuration settings in a table format.
func DisplayConfigSettings(config Config) {
	if chatConfig.Chat.Verbose {
		fmt.Println("Current configuration settings:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Setting\tValue")
		fmt.Fprintln(w, "-------\t-----")
		fmt.Fprintf(w, "Model\t%s\n", config.Chat.Model)
		fmt.Fprintf(w, "API URL\t%s\n", config.Chat.ApiURL)
		fmt.Fprintf(w, "Log Length\t%d\n", config.Chat.LogLength)
		fmt.Fprintf(w, "Verbose\t%t\n", config.Chat.Verbose)
		fmt.Fprintf(w, "Token Limit\t%d\n", config.Chat.TokenLimit)
		w.Flush()
	}
}

// Setup initializes the chat model by checking the configuration file and model status.
func Setup(configFile string) error {
	err := CheckConfigFile(configFile)
	if err != nil {
		return err
	}

	err = CheckModelRunning(chatConfig.Chat.ApiURL)
	if err != nil {
		return err
	}

	DisplayConfigSettings(chatConfig)
	return nil
}
