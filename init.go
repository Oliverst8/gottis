package main

import (
	"fmt"
	"os"
)

func Init(projectName string) {
	fmt.Println("Initialising project")
	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create directory")
		fmt.Println(err)
	}
	os.Chdir("./" + projectName)
	//os.
}
