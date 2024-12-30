package main

import (
	"fmt"

	spellnumber "github.com/josecleiton/spellnumber"
)

func main() {
	lexer := spellnumber.NewLexer(nil)

	lexer.NextLine()

	fmt.Println(lexer.Tokens)
}
