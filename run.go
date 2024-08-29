package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func getAllFiles(extensions []string) []string {

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return []string{}
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

func Run(filename string, language Language) {
	if language.CompilerName != "" {
		compile(language)
	}
}

func compile(language Language) {
	filenames := getAllFiles(language.Extensions)

	cmd := exec.Command(language.CompilerName, filenames...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not compile")
		fmt.Println("Tried to run: " + cmd.String())
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Standard Output:\n%s\n", stdout.String())
		fmt.Printf("Standard Error:\n%s\n", stderr.String())
	} else {
		// Print output if the command succeeds
		fmt.Printf("Compilation succeeded\n%s\n", stdout.String())
	}

}
