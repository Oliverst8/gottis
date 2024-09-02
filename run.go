package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Run() {

}

func Compile(language Language) {
	filenames := getAllFileNames(language.Extensions)

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
