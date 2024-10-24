package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Capitalize(word string) string {
	return strings.ToUpper(word[:1]) + word[1:]
}

func createProjectFile(projectName string, lang string) string {

	currentLang, err := GetLanguage(lang)

	if err != nil {
		HandleError("could not retrieve the language", err)
	}

	fileExtention := currentLang.Extensions[0]

	capatalizedProjectName := Capitalize(projectName)

	fileContent := currentLang.Boilerplate(projectName)

	fileName := capatalizedProjectName + fileExtention

	err = os.WriteFile(fileName, []byte(fileContent), os.ModePerm)

	if err != nil {
		HandleError("could not write to file", err)
	}

	return fileName
}

func downloadTestFiles(problemName string) error {

	sampleFilesUrl := "https://open.kattis.com/problems/" + problemName + "/file/statement/samples.zip"

	response, err := http.Get(sampleFilesUrl)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", response.StatusCode)
	}

	file, err := os.Create(problemName + ".zip")

	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil

}

func unzipFile(filename string) error {
	zipFile, err := zip.OpenReader(filename)

	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		rc, err := file.Open()

		if err != nil {
			return fmt.Errorf("failed to open file in zip archive: %w", err)
		}

		defer rc.Close()

		createdFile, err := os.Create(file.Name)

		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer createdFile.Close()

		if _, err := io.Copy(createdFile, rc); err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}

	}
	return nil
}

func deleteFile(fileName string) error {
	return os.Remove(fileName)
}

func getTestFiles(problemName string) {

	err := downloadTestFiles(problemName)

	if err != nil {
		HandleError("problem encountered during the downloading of test files", err)
	}

	err = unzipFile(problemName + ".zip")

	if err != nil {
		HandleError("problem encountered during the unzipping of the testfiles", err)
	}
	err = deleteFile(problemName + ".zip")
	if err != nil {
		HandleError("problem encountered during the deletion of the zipfile", err)
	}
}

func Init(projectName string, lang string) {

	fmt.Println("Initialising project")
	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		HandleError("could not create directory", err)
	}
	err = os.Chdir("./" + projectName)
	if err != nil {
		HandleError("could not change directory", err)
	}

	mainFile := createProjectFile(projectName, lang)

	getTestFiles(projectName)
	language, err := GetLanguage(lang)
	if err != nil {
		HandleError("could not retrieve the language", err)
	}
	CreateProjectConfigFile(mainFile, language, projectName)
}
