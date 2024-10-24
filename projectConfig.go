package main

import (
	"encoding/json"
	"os"
)

type ProjectConfig struct {
	Language string `json:"language"`
	MainFile string `json:"mainFile"`
	Problem  string `json:"problem"`
}

func CreateProjectConfigFile(mainFile string, language Language, problem string) {
	var projectConfig ProjectConfig

	projectConfig.MainFile = mainFile
	projectConfig.Language = language.Name
	projectConfig.Problem = problem

	jsonData, err := json.MarshalIndent(projectConfig, "", " ")

	if err != nil {
		HandleError("could not marshal to json", err)
	}

	err = os.WriteFile(".gottis", jsonData, os.ModePerm)

	if err != nil {
		HandleError("could not write json to the file", err)
	}
}

func GetProjectConfig() (ProjectConfig, error) {
	jsonFile, err := os.Open(".gottis")

	if err != nil {
		return ProjectConfig{}, err
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			HandleError("could not close the json file", err)
		}
	}(jsonFile)

	var currentProjectConfig ProjectConfig

	err = json.NewDecoder(jsonFile).Decode(&currentProjectConfig)

	if err != nil {
		return ProjectConfig{}, err
	}

	return currentProjectConfig, nil
}
