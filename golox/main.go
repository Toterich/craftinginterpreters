package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func run(data string) {
	scanner := Scanner{source: data}
	tokens := scanner.scanTokens()
	fmt.Println("Tokens:", tokens)
}

func runFile(file string) {
	data, err := os.ReadFile(file)
	check(err)
	run(string(data))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		check(err)
		run(line)
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
