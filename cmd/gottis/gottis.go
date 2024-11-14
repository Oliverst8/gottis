package main

import (
	"github.com/oliverst8/gottis/internal"
	"github.com/oliverst8/gottis/internal/tui"
	"log"
	"os"
	"strings"
)

func Sanitize(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

func main() {
	_, err := internal.GetConfig()

	// Setup gottis if it hasn't been setup before
	if err != nil {
		internal.Setup()
		return
	}

	if len(os.Args) < 2 {
		tui.Draw()
		//log.Fatal("Please supply an argument when using gottis.\n\"gottis <argument>\" see \"gottis help\" for more info")
	}
	choice := Sanitize(os.Args[1])

	switch {
	case choice == "i" || choice == "init":
		if len(os.Args) != 3 {
			log.Fatal("Please supply a name for the Kattis excercise when initializing. See \"gottis help\" for more info")
		}
		internal.Init(os.Args[2], "java")
	case choice == "t" || choice == "test":

		internal.Test()
	case choice == "s" || choice == "submit":

		internal.Submit()
	case choice == "setup":
		internal.Setup()
	case choice == "h" || choice == "help":
		internal.Help()
	case choice == "o" || choice == "open":
		internal.Open()
	default:
		panic("Not a recognized command please see \"gottis help\"")
	}
}
