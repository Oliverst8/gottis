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
		Setup()
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Please supply an argument when using gottis.\n\"gottis <argument>\" see \"gottis help\" for more info")
		return
	}
	choice := Sanitize(os.Args[1])
	os.Chdir("twosum")
	choice = "t"

	switch {
	case choice == "i" || choice == "init":
		if len(os.Args) != 3 {
			fmt.Println("Please supply a name for the Kattis excercise when initializing. See \"gottis help\" for more info")
			return
		}
		Init(os.Args[2], "java")
	case choice == "t" || choice == "test":
		Test()
	case choice == "s" || choice == "submit":
		fmt.Println("Submitting to Kattis")
		Submit()
	case choice == "h" || choice == "help":
		Help()
	case choice == "o" || choice == "open":
		Open()
	default:
		fmt.Println("Not a regocnized command please see \"gottis help\"")
	}

}
