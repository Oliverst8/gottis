package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Testcase struct {
	TestNumber int
	Actual     string
	Expected   string
	Passed     bool
}

func getAllFileNames(extensions []string) []string {

	files, err := os.ReadDir(".")
	if err != nil {
		HandleError("error reading current directory", err)
	}

	var foundFiles []string

	for _, file := range files {
		for _, extension := range extensions {
			if filepath.Ext(file.Name()) == extension {
				foundFiles = append(foundFiles, file.Name())
			}
		}
	}
	return foundFiles
}

func isCorrect(actual string, expected string) bool {
	return actual == expected
}

func updateProgressbar(testCases []Testcase, totalTests int) {
	whiteColor := "\033[0m"
	redColor := "\033[0;31m"
	greenColor := "\033[0;32m"
	wrongCharacter := "X"
	rightCharacter := "#"

	dashAmount := 15

	totalPassed := 0

	fmt.Printf("\033[F")
	fmt.Printf("\r%s [", strings.Repeat("-", dashAmount))
	for _, test := range testCases {
		var color string
		var character string
		if test.Passed {
			color = greenColor
			character = rightCharacter
			totalPassed++
		} else {
			color = redColor
			character = wrongCharacter
		}
		fmt.Printf("%s%s", color, character)
	}

	for _ = range totalTests - len(testCases) {
		fmt.Printf("%s-", whiteColor)
	}

	fmt.Printf("%s] %d/%d testcases passed ", whiteColor, totalPassed, totalTests)
	fmt.Printf("%s\n", strings.Repeat("-", dashAmount))
	if len(testCases) < totalTests-1 {
		fmt.Printf("Running test %d/%d...", len(testCases)+1, totalTests)
	} else {
		fmt.Printf("\033[K")
	}

}

func Test() {
	projectConfig, err := GetProjectConfig()
	if err != nil {
		HandleError("error getting project config", err)
	}

	mainFile := projectConfig.MainFile

	language, err := GetLanguage(projectConfig.Language)
	if err != nil {
		HandleError("error converting projetconfig language to type language", err)
	}

	if language.CompilerName != "" {
		Compile(language)
	}

	inFiles := getAllFileNames([]string{".in"})

	var tests []Testcase

	updateProgressbar(tests, len(inFiles))

	for index, file := range inFiles {
		cmd := exec.Command(language.RunName, mainFile)
		inputPipe, err := cmd.StdinPipe()

		fileContent, err := os.ReadFile(file)

		if err != nil {
			HandleError("error trying to read file", err)
		}

		_, err = inputPipe.Write(fileContent)

		if err != nil {
			HandleError("error trying to pipe input intro program", err)
		}

		outputPipe, err := cmd.StdoutPipe()

		if err != nil {
			HandleError("error setting stdoutpipe", err)
		}

		err = cmd.Start()

		if err != nil {
			HandleError("error trying to start program", err)
		}

		content, err := io.ReadAll(outputPipe)

		if err != nil {
			HandleError("error trying to read output pipe", err)
		}

		actualContent := string(content[:])

		actualContent = strings.Trim(actualContent, " \n")

		fileName := strings.Split(file, ".")[0]

		answerFileBytes, err := os.ReadFile(fileName + ".ans")

		expectedContent := strings.Trim(string(answerFileBytes[:]), " \n")

		tests = append(tests, Testcase{TestNumber: index, Actual: actualContent, Expected: expectedContent, Passed: isCorrect(actualContent, expectedContent)})

		updateProgressbar(tests, len(inFiles))

	}

	for _, testcase := range tests {
		if !testcase.Passed {
			fmt.Printf("---------- Testcase %d ----------\n", testcase.TestNumber)
			fmt.Printf("Expected output:\n%s\n", testcase.Expected)
			fmt.Printf("Actual output:\n%s\n", testcase.Actual)
			fmt.Printf("---------------------------------\n\n")
		}
	}
	if language.CompilerName != "" {
		DeleteFilesWithExtension(".class")
	}
}
