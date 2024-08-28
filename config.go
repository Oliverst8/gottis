package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	DefaultLang string `json:"defaultLang"`
	User        User   `json:"user,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func GetConfigPath() (string, error) {
	var configDir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	appName := "gottis"

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), appName)
	case "darwin":
		configDir = filepath.Join(homeDir, "Library", "Application Support", appName)
	case "linux":
		configDir = filepath.Join(homeDir, ".config", appName)
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	return filepath.Join(configDir, "config.json"), nil
}

func GetConfig() (map[string]string, error) {
	configPath, err := GetConfigPath()

	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	var config Config

	err = json.NewDecoder(jsonFile).Decode(&config)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
