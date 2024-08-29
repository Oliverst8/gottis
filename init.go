package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Capatalize(word string) string {
	return strings.ToUpper(word[:1]) + word[1:]
}

func createProjectFile(projectName string, lang string) {

	fileExtention := ".java"

	var fileContent string

	capatalizedProjectName := Capatalize(projectName)

	switch lang {
	case "java":
		fileContent = "import java.util.Scanner\n" + "\n" + "public class " + capatalizedProjectName + " {\n" + "   public static void main(String[] args) {\n" + "       Scanner sc = new Scanner();\n" + "       sc.close();\n" + "   }\n" + "}"
	default:
		fileContent = ""
	}

	err := os.WriteFile(capatalizedProjectName+fileExtention, []byte(fileContent), os.ModePerm)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
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
	return nil
}

func getTestFiles(problemName string) {

	err := downloadTestFiles(problemName)

	if err != nil {
		fmt.Println(err)
	}

	err = unzipFile(problemName + ".zip")

}

func Init(projectName string, lang string) {

	fmt.Println("Initialising project")
	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create directory")
		fmt.Println(err)
	}
	os.Chdir("./" + projectName)

	createProjectFile(projectName, lang)

	getTestFiles(projectName)

}
