package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Configuration struct {
	ClientId              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	AccessToken           string `json:"access_token"`
	configurationFilePath string
}

func NewConfiguration() (*Configuration, error) {
	baseConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Cannot find user config dir :", err)
		return nil, err
	}
	appConfigDir := filepath.Join(baseConfigDir, "thingiverse-cli")
	err = os.MkdirAll(appConfigDir, 0755)
	if err != nil {
		fmt.Println("Cannot create config dir :", err)
		return nil, err
	}
	configFile := filepath.Join(appConfigDir, "config.json")
	return &Configuration{
		configurationFilePath: configFile,
	}, nil
}

func (c *Configuration) Load() error {
	data, err := os.ReadFile(c.configurationFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c *Configuration) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.configurationFilePath, data, 0644)
}
