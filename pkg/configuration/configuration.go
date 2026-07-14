package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	ClientId         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	AccessToken      string `json:"access_token"`
	configFilePath   string
	appConfigDirPath string
	appConfigDirName string
}

func NewConfiguration() (*Configuration, error) {
	baseConfigDir, err := os.UserConfigDir()
	appConfigDirName := "thingiverse-cli"

	if err != nil {
		fmt.Println("Cannot find user config dir :", err)
		return nil, err
	}
	appConfigDirPath := filepath.Join(baseConfigDir, appConfigDirName)
	err = os.MkdirAll(appConfigDirPath, 0755)
	if err != nil {
		fmt.Println("Cannot create config dir :", err)
		return nil, err
	}
	configFile := filepath.Join(appConfigDirPath, "config.json")
	return &Configuration{
		configFilePath:   configFile,
		appConfigDirPath: appConfigDirPath,
	}, nil
}

func (c *Configuration) Load() error {
	data, err := os.ReadFile(c.configFilePath)
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
	return os.WriteFile(c.configFilePath, data, 0644)
}

func (c *Configuration) RemoveConfigFolder() error {
	if !strings.HasSuffix(c.appConfigDirPath, c.appConfigDirName) {
		return fmt.Errorf("App config dir seems being badly setted. Action stopped")
	}
	err := os.RemoveAll(c.appConfigDirPath)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configuration) GenerateConnectionUrl() string {
	redirectURI := "https%3A%2F%2Fakrobate.github.io%2Fthingiverse-cli%2Ftoken.html"
	url := fmt.Sprintf("https://www.thingiverse.com/login/oauth/authorize?client_id=%s&response_type=token&redirect_uri=%s",
		c.ClientId,
		redirectURI,
	)
	return url
}

func (c *Configuration) ConfigurationExists() bool {
	_, err := os.Stat(c.configFilePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
