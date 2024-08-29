package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type projectConfig struct {
	Language string `json:"language"`
	MainFile string `json:"mainFile"`
}

func Capatalize(word string) string {
	return strings.ToUpper(word[:1]) + word[1:]
}

func createProjectFile(projectName string, lang string) string {

	fileExtention := ".java"

	var fileContent string

	capatalizedProjectName := Capatalize(projectName)

	switch lang {
	case "java":
		fileContent = "import java.util.Scanner\n" + "\n" + "public class " + capatalizedProjectName + " {\n" + "   public static void main(String[] args) {\n" + "       Scanner sc = new Scanner();\n" + "       sc.close();\n" + "   }\n" + "}"
	default:
		fileContent = ""
	}

	fileName := capatalizedProjectName + fileExtention

	err := os.WriteFile(fileName, []byte(fileContent), os.ModePerm)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
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
		fmt.Println(err)
	}

	err = unzipFile(problemName + ".zip")

	if err != nil {
		fmt.Println(err)
	}

	err = deleteFile(problemName + ".zip")

}

func createProjectConfigFile(mainFile string) {
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

func Init(projectName string, lang string) {

	fmt.Println("Initialising project")
	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create directory")
		fmt.Println(err)
	}
	os.Chdir("./" + projectName)

	mainFile := createProjectFile(projectName, lang)

	getTestFiles(projectName)

	createProjectConfigFile(mainFile)

}
