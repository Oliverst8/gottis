package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type projectConfig struct {
	Language string `json:"language"`
	MainFile string `json:"mainFile"`
}

func CreateProjectConfigFile(mainFile string) {
	var projectConfig projectConfig

	projectConfig.MainFile = mainFile
	config, err := GetConfig()

	if err != nil {
		fmt.Println(err)
	}

	projectConfig.Language = config.DefaultLang

	jsonData, err := json.MarshalIndent(projectConfig, "", " ")

	err = os.WriteFile(".gottis", jsonData, os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
}

func GetProjectConfig() (projectConfig, error) {
	jsonFile, err := os.Open(".gottis")

	if err != nil {
		fmt.Println(err)
		return projectConfig{}, err
	}

	defer jsonFile.Close()

	var currentProjectConfig projectConfig

	err = json.NewDecoder(jsonFile).Decode(&currentProjectConfig)

	if err != nil {
		fmt.Println(err)
		return projectConfig{}, err
	}

	return currentProjectConfig, nil

}
