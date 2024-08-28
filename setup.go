package main

import (
	"errors"
	"fmt"
	"os"
)

func Setup() {
	path, err := GetConfigPath()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		e := os.Mkdir("gottis", os.ModePerm)
		if e != nil {
			fmt.Println("Cannot create directory")
			fmt.Println(e)
			return
		}
	} else if err != nil {
		fmt.Println(err)
		return
	}
	os.Chdir(path)
	getUserInputForConfig()
}

func getUserInputForConfig() {
	fmt.Println("Do yu want...")
}
