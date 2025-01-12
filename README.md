# Spellnumber Project

## Overview

The Spellnumber project is a Go library that provides functionality for parsing and spelling out numbers in Portuguese. The library consists of several packages, including `spellnumber`, which contains the core functionality, and `cmd`, which contains a command-line interface for testing and demonstrating the library.

This project is a rewrite of a C project coded in a bachelor course [analisador-extenso](https://github.com/josecleiton/analisador-extenso)

## Packages

### spellnumber

The `spellnumber` package contains the core functionality of the project. It includes several sub-packages:

* `analex`: This package contains the lexical analyzer, which is responsible for breaking down input text into tokens.
* `parser`: This package contains the parser, which takes the tokens produced by the lexical analyzer and produces a `*big.Int`.
* `speller`: This package contains the speller, which takes the `*big.Int` and produces a string representation of the number.

### cmd

The `cmd` package contains a command-line interface for testing and demonstrating the `spellnumber` library. It includes a single command, `main`, which takes a string input and produces a spelled-out representation of the number.

## Functions

### spellnumber.Lexer.ParseLine

This function takes a string input and produces a slice of tokens.

### spellnumber.Parser.Parse

This function takes a slice of tokens and produces a `*big.Int` as result.

### spellnumber.Speller.Spell

This function takes a `*big.Int` and produces a string representation of the number.

## Usage

To use the `spellnumber` library, create a new instance of the `Lexer`, `Parser`, and `Speller` structs, and call the corresponding methods to parse and spell out a number.

```go
package main

import (
	"fmt"
	"github.com/josecleiton/spellnumber"
)

func main() {
	lexer := spellnumber.NewLexer(nil, false)
	tokens := lexer.ParseLine("cento")
	parser := spellnumber.NewParser(tokens, false)
	result := parser.Parse()
	speller := spellnumber.NewSpeller(false)
	fmt.Println(speller.Spell(result))
}
```


## Diagram of Lexer

Beatiful and complete lexer diagram

<img src="/docs/diagram.png" alt="Lexer Diagram" width="600" height="auto">
