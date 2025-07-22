package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

func Read() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}

func write(config Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonConfig, os.FileMode(0644)) // owners read/write, group and others read
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homePath + "/" + configFileName, nil
}
