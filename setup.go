package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var langs = [5]string{"java", "go", "c", "python", "c#"}

func Setup() {
	path, err := GetConfigDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		e := os.Mkdir(path, os.ModePerm)
		if e != nil {
			fmt.Println("Cannot create directory")
			fmt.Println(e)
			return
		}
	} else if err != nil {
		fmt.Println(err)
		return
	}
	err = os.Chdir(path)
	if err != nil {
		return
	}
	setupSubmit, lang := getUserInputForConfig()

}

func getUserInputForConfig() (bool, int) {
	fmt.Println("Choose a default language for gottis:")
	printLanguages()
	lang, err := readInt()
	validOption := isValidAndInRange(lang, err, 1, len(langs))
	for !validOption {
		fmt.Println("Not a valid option. Please input a number from the following list:")
		printLanguages()
		lang, err = readInt()
		validOption = isValidAndInRange(lang, err, 1, len(langs))
	}

	fmt.Println("Do you want to be able to submit to kattis through gottis? [y/N]")
	setupSubmit := readBool()

	return setupSubmit, lang
}

func isValidAndInRange(num int, err error, a int, b int) bool {
	return err == nil && num >= a && num <= b
}

func printLanguages() {
	for index, lang := range langs {
		fmt.Printf("%d. %s\n", index+1, lang)
	}
}

func readBool() bool {
	var input string
	readInput(&input)
	input = strings.ToLower(input)
	if input == "y" || input == "yes" {
		return true
	}
	return false
}

func readInput(answer *string) {
	_, err := fmt.Scanf("%s \n", answer)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readInt() (int, error) {
	var input string
	readInput(&input)
	readAnswer, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return readAnswer, nil
}
