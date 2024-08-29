package main

import "fmt"

func Test() {
	projectConfig, err := GetProjectConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	mainFile := projectConfig.MainFile

	language, err := GetLanguage(projectConfig.Language)
	if err != nil {
		fmt.Println(err)
		return
	}

	Run(mainFile, language)

}
