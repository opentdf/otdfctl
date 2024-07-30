package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Model     string `json:"model"`
	Verbosity string `json:"verbosity"`
	ApiURL    string `json:"apiURL"`
}

var chat_config Config

func init() {
	err := LoadConfig("chat_config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}
}

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&chat_config)
	if err != nil {
		return fmt.Errorf("could not decode config JSON: %v", err)
	}

	return nil
}
