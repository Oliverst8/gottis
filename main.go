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
	fmt.Println("Hello World!")
	choice := Sanitize(os.Args[1])

	switch choice {
	case "i":
	case "init":

	}

}
