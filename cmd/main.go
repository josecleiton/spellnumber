package main

import (
	"flag"
	"fmt"
	"log"

	spellnumber "github.com/josecleiton/spellnumber"
)

var verboseFlag bool

func init() {
	flag.BoolVar(&verboseFlag, "v", false, "verbose output")

	flag.Parse()
}

func main() {
	lexer := spellnumber.NewLexer(nil, verboseFlag)

	tokens := lexer.ParseLine("zero")

	log.Printf("Tokens: %v\n", tokens)

	parser := spellnumber.NewParser(tokens, verboseFlag)

	result, err := parser.Parse()

	if err != nil {
		log.Fatalf("Parser Error: %v\n", err)
	}

	fmt.Printf("Result: %v\n", result)

	speller := spellnumber.NewSpeller(verboseFlag)

	fmt.Printf("Spell: %v\n", speller.Spell(result))
}
