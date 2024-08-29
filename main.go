package main

import (
	"fmt"
	"os"
	"strings"
)

func Sanitize(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

func main() {

	_, err := GetConfig()

	if err != nil {
		//Setup()
		//return
	}

	fmt.Println("Hello World!")
	if len(os.Args) < 2 {
		fmt.Println("Please supply an argument when using gottis.\n\"gottis <argument>\" see \"gottis help\" for more info")
		return
	}
	choice := Sanitize(os.Args[1])

	fmt.Println("Entered switch case")
	switch choice {
	case "i":
	case "init":
		if len(os.Args) != 3 {
			fmt.Println("Please supply a name for the Kattis excercise when initializing. See \"gottis help\" for more info")
			break
		}
		Init(os.Args[2], "java")
		break
	case "t":
	case "test":
		Test()
		break
	case "s":
	case "submit":
		fmt.Println("Submitting to Kattis")
		Submit()
		break
	case "h":
	case "help":
		Help()
		break
	case "o":
	case "open":
		Open()
		break
	default:
		fmt.Println("Not a regocnized command please see \"gottis help\"")
		break
	}
}
