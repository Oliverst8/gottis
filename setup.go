package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

	var user User
	if setupSubmit {
		user = getKattisInfo()
	}

	var config Config
	config.DefaultLang = langs[lang].Name
	config.User = user

	jsonData, err := json.MarshalIndent(config, "", " ")

	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile("config.json", jsonData, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}

func getKattisInfo() User {
	var user User
	fmt.Println("We need information from your kattis account, the information can be found on the following link (You need to login first):")
	fmt.Println("https://open.kattis.com/download/kattisrc")
	fmt.Println("Please input your Kattis username:")
	var answer string
	readInput(&answer)
	user.Username = strings.Trim(answer, " \n")
	fmt.Println("Please input your Kattis Token:")
	readInput(&answer)
	user.Token = strings.Trim(answer, " \n")
	return user
}

func getUserInputForConfig() (bool, int) {
	fmt.Println("Choose a default language for gottis:")
	printLanguages()
	lang, err := readInt()
	validOption := isValidAndInRange(lang, err, 0, len(langs)-1)
	for !validOption {
		fmt.Println("Not a valid option. Please input a number from the following list:")
		printLanguages()
		lang, err = readInt()
		validOption = isValidAndInRange(lang, err, 0, len(langs)-1)
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
		fmt.Printf("%d. %s\n", index, lang.Name)
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
