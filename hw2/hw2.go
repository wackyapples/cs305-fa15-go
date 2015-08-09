package main

import (
	"fmt"
	"os"
	"strings"
)

/*
 *
 */
func main() {
	// Command line args
	args := os.Args

	// Check for required parameters
	if len(args) < 2 {
		fmt.Printf("%s: Incorrect number of parameters.\n", args[0])

		os.Exit(1)
	}

	for i, str := range args[1:] {
		// print the word
		fmt.Printf("%s", capsString(str))

		// If it's the second or last word, make a newline,
		// otherwise, print a space.
		if i%2 == 1 || i == len(args[1:])-1 {
			fmt.Printf("\n")
		} else {
			fmt.Printf(" ")
		}
	}

	// fmt.Printf("%v\n", args)
}

func capsString(s string) string {
	if len(s)%2 == 1 {
		return strings.ToUpper(s)
	}

	return s
}
