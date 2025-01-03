package main

import (
	"flag"
	"fmt"

	spellnumber "github.com/josecleiton/spellnumber"
)

var verboseFlag bool

func init() {
	flag.BoolVar(&verboseFlag, "v", false, "verbose output")

	flag.Parse()
}

func main() {
	lexer := spellnumber.NewLexer(nil, verboseFlag)

	tokens := lexer.ParseLine("cento")

	fmt.Printf("Tokens: %v\n", tokens)

	parser := spellnumber.NewParser(tokens, verboseFlag)

	result := parser.Parse()

	fmt.Printf("Result: %v\n", result)

	speller := spellnumber.NewSpeller(verboseFlag)

	fmt.Printf("Spell: %v\n", speller.Spell(result))
}
