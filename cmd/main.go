package main

import (
	"fmt"

	spellnumber "github.com/josecleiton/spellnumber"
)

func main() {
	lexer := spellnumber.NewLexer(nil)

	tokens := lexer.ParseLine("duzentos mil e trezentos e cinquenta e um")

	fmt.Printf("Tokens: %v\n", tokens)

	parser := spellnumber.NewParser(tokens)

	result := parser.Parse()

	fmt.Printf("Result: %v\n", result)
}
