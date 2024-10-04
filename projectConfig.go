package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type projectConfig struct {
	Language string `json:"language"`
	MainFile string `json:"mainFile"`
	Problem  string `json:"problem"`
}

func CreateProjectConfigFile(mainFile string, language Language, problem string) {
	var projectConfig projectConfig

	projectConfig.MainFile = mainFile
	projectConfig.Language = language.Name
	projectConfig.Problem = problem

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
