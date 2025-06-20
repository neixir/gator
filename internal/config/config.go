package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DbUrl           string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

/*
Export a Read function that reads the JSON file found at ~/.gatorconfig.json and returns a Config struct.
It should read the file from the HOME directory, then decode the JSON string into a new Config struct.
I used os.UserHomeDir to get the location of HOME.
*/
func Read() (Config, error) {
	config := Config{}

	configFile, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("Error getting HOME directory: %v", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("Error rading file: %v", err)
	}

	var currentConfig Config
	err = json.Unmarshal(data, &currentConfig)
	if err != nil {
		return config, fmt.Errorf("Error parsing JSON data: %v", err)
	}

	return currentConfig, nil
}

// Export a SetUser method on the Config struct that writes the config struct to the JSON file after setting the current_user_name field.
func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	err := write(*c)

	return err
}

// I also wrote a few non-exported helper functions and added a constant to hold the filename.
// But you can implement the internals of the package however you like.
func getConfigFilePath() (string, error) {
	configFile := ""

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	} else {
		configFile = filepath.Join(homeDir, configFileName)
	}

	return configFile, nil
}

func write(cfg Config) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	newConfig, _ := json.Marshal(cfg)
	err = os.WriteFile(configFile, []byte(newConfig), 0644)
	
	return nil
}
