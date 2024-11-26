package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"toterich/golox/interp"
	"toterich/golox/parse"
	"toterich/golox/util"
)

var scanner parse.Scanner
var parser parse.Parser

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
	tokens, errs := scanner.ScanTokens(data)
	//	fmt.Println(tokens)
	if errs != nil {
		util.LogErrors(errs...)
		return fmt.Errorf("errors in Scanner")
	}

	expr, errs := parser.Parse(tokens)
	//	fmt.Print(expr.PrettyPrint())
	if errs != nil {
		util.LogErrors(errs...)
		return fmt.Errorf("errors in Parser")
	}

	value, err := interp.Evaluate(expr)
	if err != nil {
		util.LogErrors(err)
		return fmt.Errorf("errors in Interpreter")
	}

	fmt.Println(value)

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
