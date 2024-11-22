package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"toterich/golox/parser"
)

// Check for error and exit
// If exitCode is 0, only log error and don't exit
func check(e error, exitCode int) {
	if e != nil {
		log.Print(e)
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}
}

func run(data string) error {
	scanner := parser.Scanner{}
	tokens := scanner.ScanTokens(data)
	fmt.Printf("Tokens: %v \n", tokens)
	return nil
}

func runFile(file string) {
	data, err := os.ReadFile(file)
	check(err, 1)
	err = run(string(data))
	check(err, 2)
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		check(err, 0)
		err = run(line)
		check(err, 0)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: golox [script.lox]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
