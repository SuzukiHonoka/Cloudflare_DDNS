package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// check if any error to panic
func check(err error) {
	if err != nil {
		log.Fatal("Fatal Error: " + err.Error())
	}
}

// getInput gets the input from terminal prompt
func getInput(question string) string {
	// prompt question first
	fmt.Print(question + ": ")
	reader := bufio.NewReader(os.Stdin)
	rp, err := reader.ReadString('\n')
	check(err)
	// replace space from either windows or linux terminal
	return strings.Trim(strings.Trim(rp, "\r\n"), "\n")
}

// getInputBool gets the input and parse to bool
func getInputBool(question string) bool {
	return strings.ToLower(getInput(question+"? (y/n)")) == "y"
}
