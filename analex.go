package spellnumber

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
)

type TokenType int

type Token struct {
	Type   TokenType
	Value  string
	Spell  string
	Number *big.Int
}

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_LEFT_BRACKET
	TOKEN_RIGHT_BRACKET

	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_DIVIDE
	TOKEN_TIMES
	TOKEN_POWER
	TOKEN_FACTORIAL
	TOKEN_MOD

	TOKEN_NUMBER
	TOKEN_NUMBER_PARSED
)

type Lexer struct {
	scannerStdIn *bufio.Reader
	numberDict   map[string]numberState
}

type numberState struct {
	state int
	value string
}

func NewLexer(inputFile *os.File) *Lexer {
	file := inputFile

	if file == nil {
		file = os.Stdin
	}

	return &Lexer{
		scannerStdIn: bufio.NewReader(file),
		numberDict: map[string]numberState{
			"um":              {state: 6, value: "1"},
			"dois":            {state: 6, value: "2"},
			"tres":            {state: 6, value: "3"},
			"quatro":          {state: 6, value: "4"},
			"cinco":           {state: 6, value: "5"},
			"seis":            {state: 6, value: "6"},
			"sete":            {state: 6, value: "7"},
			"oito":            {state: 6, value: "8"},
			"nove":            {state: 6, value: "9"},
			"dez":             {state: 6, value: "10"},
			"onze":            {state: 6, value: "11"},
			"doze":            {state: 6, value: "12"},
			"treze":           {state: 6, value: "13"},
			"quatorze":        {state: 6, value: "14"},
			"quinze":          {state: 6, value: "15"},
			"dezesseis":       {state: 6, value: "16"},
			"dezessete":       {state: 6, value: "17"},
			"dezoito":         {state: 6, value: "18"},
			"dezenove":        {state: 6, value: "19"},
			"vinte":           {state: 7, value: "20"},
			"trinta":          {state: 7, value: "30"},
			"quarenta":        {state: 7, value: "40"},
			"cinquenta":       {state: 7, value: "50"},
			"sessenta":        {state: 7, value: "60"},
			"setenta":         {state: 7, value: "70"},
			"oitenta":         {state: 7, value: "80"},
			"noventa":         {state: 7, value: "90"},
			"cem":             {state: 9, value: "100"},
			"cento":           {state: 10, value: "100"},
			"duzentos":        {state: 10, value: "200"},
			"trezentos":       {state: 10, value: "300"},
			"quatrocentos":    {state: 10, value: "400"},
			"quinhentos":      {state: 10, value: "500"},
			"seiscentos":      {state: 10, value: "600"},
			"setecentos":      {state: 10, value: "700"},
			"oitocentos":      {state: 10, value: "800"},
			"novecentos":      {state: 10, value: "900"},
			"mil":             {state: 11, value: "1000"},
			"milhao":          {state: 11, value: "1000000"},
			"milhoes":         {state: 11, value: "1000000"},
			"bilhao":          {state: 11, value: "1000000000"},
			"bilhoes":         {state: 11, value: "1000000000"},
			"trilhao":         {state: 11, value: "1000000000000"},
			"trilhoes":        {state: 11, value: "1000000000000"},
			"quadrilhao":      {state: 11, value: "1000000000000000"},
			"quadrilhoes":     {state: 11, value: "1000000000000000"},
			"quintilhao":      {state: 11, value: "1000000000000000000"},
			"quintilhoes":     {state: 11, value: "1000000000000000000"},
			"sextilhao":       {state: 11, value: "1000000000000000000000"},
			"sextilhoes":      {state: 11, value: "1000000000000000000000"},
			"septilhao":       {state: 11, value: "1000000000000000000000000"},
			"septilhoes":      {state: 11, value: "1000000000000000000000000"},
			"octilhao":        {state: 11, value: "1000000000000000000000000000"},
			"octilhoes":       {state: 11, value: "1000000000000000000000000000"},
			"nonilhao":        {state: 11, value: "1000000000000000000000000000000"},
			"nonilhoes":       {state: 11, value: "1000000000000000000000000000000"},
			"decilhao":        {state: 11, value: "1000000000000000000000000000000000"},
			"decilhoes":       {state: 11, value: "1000000000000000000000000000000000"},
			"undecilhao":      {state: 11, value: "1000000000000000000000000000000000000"},
			"undecilhoes":     {state: 11, value: "1000000000000000000000000000000000000"},
			"duodecilhao":     {state: 11, value: "1000000000000000000000000000000000000000"},
			"duodecilhoes":    {state: 11, value: "1000000000000000000000000000000000000000"},
			"tridecilhao":     {state: 11, value: "1000000000000000000000000000000000000000000"},
			"tridecilhoes":    {state: 11, value: "1000000000000000000000000000000000000000000"},
			"quatradecilhao":  {state: 11, value: "1000000000000000000000000000000000000000000000"},
			"quatradecilhoes": {state: 11, value: "1000000000000000000000000000000000000000000000"},
			"e":               {state: 8, value: "0"},
		},
	}
}

func (l *Lexer) NextLine() {
	line, _ := l.scannerStdIn.ReadString('\n')

	l.ParseLine(line)
}

func (l *Lexer) ParseLine(line string) []Token {
	tokens := make([]Token, 0, 64)

	words := strings.Split(line, " ")

	index := 0

	state := 0

	numberTokens := make([]Token, 0)
	for {
		if index >= len(words) {
			if len(numberTokens) > 0 {
				tokens = append(tokens, getNumberTokenFromList(numberTokens))
			}
			break
		}

		lexeme := words[index]

		fmt.Printf("state: %d | lexeme: %s\n", state, lexeme)

		if strings.Contains(lexeme, "\n") {
			lexeme = strings.ReplaceAll(lexeme, "\n", "")
		}

		if state == 0 {
			if lexeme == "mais" {
				tokens = append(tokens, Token{Type: TOKEN_PLUS, Value: lexeme})
			} else if lexeme == "menos" {
				tokens = append(tokens, Token{Type: TOKEN_MINUS, Value: lexeme})
			} else if lexeme == "vezes" {
				tokens = append(tokens, Token{Type: TOKEN_TIMES, Value: lexeme})
			} else if lexeme == "elevado" {
				state = 1
			} else if lexeme == "abre" {
				state = 2
			} else if lexeme == "fecha" {
				state = 3
			} else if lexeme == "fatorial" {
				state = 4
			} else if lexeme == "dividido" {
				state = 5
			} else {
				if val, ok := l.numberDict[lexeme]; ok {
					numberTokens = append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme})

					state = val.state
				} else {
					tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
					break
				}
			}
		} else if state == 1 {
			if lexeme != "por" {
				tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
				break
			}

			tokens = append(tokens, Token{Type: TOKEN_POWER, Value: lexeme})

			state = 0
		} else if state == 2 {
			if lexeme != "parentese" {
				tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
				break
			}

			tokens = append(tokens, Token{Type: TOKEN_LEFT_BRACKET, Value: lexeme})

			state = 0
		} else if state == 3 {
			if lexeme != "parentese" {
				tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
				break
			}

			tokens = append(tokens, Token{Type: TOKEN_RIGHT_BRACKET, Value: lexeme})

			state = 0
		} else if state == 4 {
			if lexeme != "de" {
				tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
				break
			}

			tokens = append(tokens, Token{Type: TOKEN_FACTORIAL, Value: lexeme})

			state = 0
		} else if state == 5 {
			if lexeme != "por" {
				tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
				break
			}

			tokens = append(tokens, Token{Type: TOKEN_DIVIDE, Value: lexeme})

			state = 0
		} else if state == 6 {
			if val, ok := l.numberDict[lexeme]; ok {
				if val.state != 11 {
					tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
					break
				}

				numberTokens = append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme})

				state = 8
			} else {
				index--

				tokens = append(tokens, getNumberTokenFromList(numberTokens))

				state = 0
			}

		} else if state == 7 {
			if _, ok := l.numberDict[lexeme]; lexeme != "e" && ok {
				tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
				break
			}

			if lexeme == "e" {
				state = 6
			} else {
				state = 0
				index--

				tokens = append(tokens, getNumberTokenFromList(numberTokens))
			}
		} else if state == 8 {
			if val, ok := l.numberDict[lexeme]; ok {
				if val.state > 10 {
					tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
					break
				}

				state = val.state

				if val.state != 8 {

					numberTokens = append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme})
				}
			} else {
				index--

				tokens = append(tokens, getNumberTokenFromList(numberTokens))

				state = 0
			}
		} else if state == 9 {
			if val, ok := l.numberDict[lexeme]; ok {
				if val.state != 11 {
					tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
					break
				}

				numberTokens = append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme})
			} else {
				index--

				tokens = append(tokens, getNumberTokenFromList(numberTokens))
			}

			state = 0
		} else if state == 9 {
			if val, ok := l.numberDict[lexeme]; ok {
				if val.state == 6 || val.state == 7 {
					state = val.state
				} else if val.state == 11 {
					state = 0

				} else {
					tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
					break
				}

				numberTokens = append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme})
			}
		} else if state == 10 {
			if val, ok := l.numberDict[lexeme]; ok {
				if val.state == 6 || val.state == 7 {
					state = val.state
				} else if val.state == 11 {
					state = 8
				} else {
					tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
					break
				}

				numberTokens = append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme})
			} else if lexeme == "e" {
				state = 9
			} else {
				state = 0
				index--

				tokens = append(tokens, getNumberTokenFromList(numberTokens))
			}
		}

		index++
	}

	return tokens
}

func getNumberTokenFromList(numberTokens []Token) Token {
	if len(numberTokens) == 0 {
		return Token{Type: TOKEN_ERROR, Value: "0"}
	}

	fmt.Println(numberTokens)

	order := 1
	orderMilhar := len("1000")

	number := big.NewInt(0)

	for i := len(numberTokens) - 1; i >= 0; i-- {
		token := numberTokens[i]

		fmt.Println(token)

		tokenOrder := len(token.Value)

		if tokenOrder >= orderMilhar {
			if order > orderMilhar && tokenOrder <= order {
				return Token{Type: TOKEN_ERROR, Value: "0"}
			}

			order = tokenOrder

			continue
		}

		currentUnit := big.NewInt(0)

		currentUnit, ok := currentUnit.SetString(token.Value, 10)

		if !ok {
			return Token{Type: TOKEN_ERROR, Value: "0"}
		}

		orderNumber := big.NewInt(int64(math.Pow(10, float64(order-1))))

		currentNumber := currentUnit.Mul(currentUnit, orderNumber)

		fmt.Printf("%s * %s = %s\n", token.Value, orderNumber.String(), currentNumber.String())

		number = number.Add(number, currentNumber)
	}

	return Token{Type: TOKEN_NUMBER_PARSED, Value: number.String(), Number: number}
}
