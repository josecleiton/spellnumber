package main

import (
	"fmt"

	spellnumber "github.com/josecleiton/spellnumber"
)

func main() {
	lexer := spellnumber.NewLexer(nil)

	tokens := lexer.ParseLine("duzentos milhoes e trezentos e cinquenta e um mais fatorial de cinco")

	fmt.Printf("Tokens: %v\n", tokens)

	parser := spellnumber.NewParser(tokens)

	result := parser.Parse()

	fmt.Printf("Result: %v\n", result)

	speller := spellnumber.NewSpeller()

	fmt.Printf("Spell: %v\n", speller.Spell(result))
}
