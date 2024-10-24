package main

import (
	"log"
	"os"
	"runtime"
	"strings"
)

func Sanitize(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

func HandleError(message string, err error) {
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	panic(message + "\n" + err.Error() + "\n" + string(buf))
}

func main() {
	_, err := GetConfig()

	// Setup gottis if it hasn't been setup before
	if err != nil {
		Setup()
		return
	}

	if len(os.Args) < 2 {
		log.Fatal("Please supply an argument when using gottis.\n\"gottis <argument>\" see \"gottis help\" for more info")
	}
	choice := Sanitize(os.Args[1])

	switch {
	case choice == "i" || choice == "init":
		if len(os.Args) != 3 {
			log.Fatal("Please supply a name for the Kattis excercise when initializing. See \"gottis help\" for more info")
		}
		Init(os.Args[2], "java")
	case choice == "t" || choice == "test":
		// Debug statement
		os.Chdir("twosum")
		Test()
	case choice == "s" || choice == "submit":
		// Debug statement
		os.Chdir("twosum")

		Submit()
	case choice == "h" || choice == "help":
		Help()
	case choice == "o" || choice == "open":
		Open()
	default:
		panic("Not a recognized command please see \"gottis help\"")
	}
}
