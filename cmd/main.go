package main

import (
	"fmt"
	"math/big"

	spellnumber "github.com/josecleiton/spellnumber"
)

func main() {
	lexer := spellnumber.NewLexer(nil)

	tokens := lexer.ParseLine("duzentos mil e trezentos e cinquenta e um")

	parser := spellnumber.NewParser()

	parser.Parse(tokens)

	fmt.Println(big.NewInt(0).MulRange(1, 100).String())

	fmt.Println(tokens)
}
