package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Access the arugments
	cmdArgs := os.Args

	// Check for correct number of args
	if len(cmdArgs) < 2 || len(cmdArgs) > 5 {
		fmt.Printf("Error: wrong number of parameters.\n"+
			"usage: %s arg1 [arg2 … arg4].\n", cmdArgs[0])

		os.Exit(1)
	}

	fmt.Printf("%d parameters\n", len(cmdArgs))

	// convert args to ints
	for _, value := range cmdArgs[1:] {
		i, err := strconv.Atoi(value)

		if err == nil {
			fmt.Printf("%d kg = %.6f TP\n", i, float32(i)/1325.0)
		}
	}

}
