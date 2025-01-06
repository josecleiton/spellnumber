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
	lexer := spellnumber.NewLexer(nil)
	lexer.SetVerbose(verboseFlag)

	tokens, err := lexer.ParseLine("trezentos milhões e mil e duzentos e três")

	if err != nil {
		log.Fatalf("Lexer Error: %v\n", err)
	}

	log.Printf("Tokens: %v\n", tokens)

	parser := spellnumber.NewParser(tokens)
	parser.SetVerbose(verboseFlag)

	result, err := parser.Parse()

	if err != nil {
		log.Fatalf("Parser Error: %v\n", err)
	}

	fmt.Printf("Result: %v\n", result)

	speller := spellnumber.NewSpeller()
	speller.SetVerbose(verboseFlag)

	fmt.Printf("Spell: %v\n", speller.Spell(result))
}
