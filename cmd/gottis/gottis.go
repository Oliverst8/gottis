package main

import (
	"github.com/oliverst8/gottis/internal"
	"log"
	"os"
	"strings"
)

func Sanitize(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

func parse() (map[string]string, []string) {
	flags := make(map[string]string)
	var args []string
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			flags[strings.TrimPrefix(arg, "-")] = os.Args[i+1]
		} else if i != 0 && !strings.HasPrefix(os.Args[i-1], "-") {
			args = append(args, arg)
		}
	}
	return flags, args
}

func main() {
	config, err := internal.GetConfig()

	// Setup gottis if it hasn't been setup before
	if err != nil {
		internal.Setup()
		return
	}

	flags, args := parse()

	if args[0] == "" {
		log.Fatal("Please supply an argument when using gottis.\n\"gottis <argument>\" see \"gottis help\" for more info")
	}

	choice := Sanitize(args[0])
	var language string
	if flags["lang"] == "" {
		language = config.DefaultLang
	} else {
		language = flags["lang"]
	}

	switch {
	case choice == "i" || choice == "init":
		if len(args) != 2 {
			log.Fatal("Please supply a name for the Kattis excercise when initializing. See \"gottis help\" for more info")
		}

		internal.Init(args[1], language)
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
