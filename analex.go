package spellnumber

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"regexp"
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
			"um":        {state: 6, value: "1"},
			"dois":      {state: 6, value: "2"},
			"tres":      {state: 6, value: "3"},
			"quatro":    {state: 6, value: "4"},
			"cinco":     {state: 6, value: "5"},
			"seis":      {state: 6, value: "6"},
			"sete":      {state: 6, value: "7"},
			"oito":      {state: 6, value: "8"},
			"nove":      {state: 6, value: "9"},
			"dez":       {state: 6, value: "10"},
			"onze":      {state: 6, value: "11"},
			"doze":      {state: 6, value: "12"},
			"treze":     {state: 6, value: "13"},
			"quatorze":  {state: 6, value: "14"},
			"quinze":    {state: 6, value: "15"},
			"dezesseis": {state: 6, value: "16"},
			"dezessete": {state: 6, value: "17"},
			"dezoito":   {state: 6, value: "18"},
			"dezenove":  {state: 6, value: "19"},
			"vinte":     {state: 7, value: "20"},
			"trinta":    {state: 7, value: "30"},
			"quarenta":  {state: 7, value: "40"},
			"cinquenta": {state: 7, value: "50"},
			"sessenta":  {state: 7, value: "60"},
			"setenta":   {state: 7, value: "70"},
			"oitenta":   {state: 7, value: "80"},
			"noventa":   {state: 7, value: "90"},
			// TODO: support "cem"
			"cem":             {state: 8, value: "100"},
			"cento":           {state: 8, value: "100"},
			"duzentos":        {state: 8, value: "200"},
			"trezentos":       {state: 8, value: "300"},
			"quatrocentos":    {state: 8, value: "400"},
			"quinhentos":      {state: 8, value: "500"},
			"seiscentos":      {state: 8, value: "600"},
			"setecentos":      {state: 8, value: "700"},
			"oitocentos":      {state: 8, value: "800"},
			"novecentos":      {state: 8, value: "900"},
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
			"e":               {state: 200, value: "0"},
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
			break
		}

		lexeme := words[index]

		fmt.Printf("state: %d | lexeme: %s\n", state, lexeme)

		if strings.Contains(lexeme, "\n") {
			lexeme = strings.ReplaceAll(lexeme, "\n", "")
		}

		if state == 0 {
			if len(numberTokens) > 0 {
				tokens = append(tokens, l.getNumberTokenFromList(numberTokens))

				clear(numberTokens)
			}

			state, numberTokens, tokens = l.q0(lexeme, numberTokens, tokens)
		} else if state == 1 {
			state, numberTokens, tokens = l.q1(lexeme, numberTokens, tokens)
		} else if state == 2 {
			state, numberTokens, tokens = l.q2(lexeme, numberTokens, tokens)
		} else if state == 3 {
			state, numberTokens, tokens = l.q3(lexeme, numberTokens, tokens)
		} else if state == 4 {
			state, numberTokens, tokens = l.q4(lexeme, numberTokens, tokens)
		} else if state == 5 {
			state, numberTokens, tokens = l.q5(lexeme, numberTokens, tokens)
		} else if state == 6 {
			state, numberTokens, tokens = l.q6(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 7 {
			state, numberTokens, tokens = l.q7(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 8 {
			state, numberTokens, tokens = l.q8(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 9 {
			state, numberTokens, tokens = l.q9(lexeme, numberTokens, tokens)
		} else if state == 10 {
			state, numberTokens, tokens = l.q10(lexeme, numberTokens, tokens)
		} else if state == 11 {
			state, numberTokens, tokens = l.q11(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 12 {
			state, numberTokens, tokens = l.q12(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else {
			tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
		}

		if len(tokens) > 0 && tokens[len(tokens)-1].Type == TOKEN_ERROR {
			return tokens
		}

		index++
	}

	if len(numberTokens) > 0 {
		tokens = append(tokens, l.getNumberTokenFromList(numberTokens))

		clear(numberTokens)
	}

	return tokens
}

func (l Lexer) q0(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "mais" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_PLUS, Value: lexeme})
	}

	if lexeme == "menos" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_MINUS, Value: lexeme})
	}

	if lexeme == "vezes" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_TIMES, Value: lexeme})
	}

	if lexeme == "elevado" {
		return 1, numberTokens, tokens
	}

	if lexeme == "abre" {
		return 2, numberTokens, tokens
	}

	if lexeme == "fecha" {
		return 3, numberTokens, tokens
	}

	if lexeme == "fatorial" {
		return 4, numberTokens, tokens
	}

	if lexeme == "dividido" {
		return 5, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok && val.state >= 6 && val.state <= 10 {
		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
}

func (l Lexer) q1(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme != "por" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_POWER, Value: lexeme})
}

func (l Lexer) checkParenthesis(lexeme string) bool {
	regex := regexp.MustCompile(`^parenteses?$`)
	return regex.MatchString(lexeme)
}

func (l Lexer) q2(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if l.checkParenthesis(lexeme) {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_LEFT_BRACKET, Value: lexeme})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
}

func (l Lexer) q3(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if l.checkParenthesis(lexeme) {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_RIGHT_BRACKET, Value: lexeme})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
}

func (l Lexer) q4(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme != "de" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})

	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_FACTORIAL, Value: lexeme})
}

func (l Lexer) q5(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme != "por" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_DIVIDE, Value: lexeme})
}

func (l Lexer) q6(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.numberDict[lexeme]; ok {
		if val.state != 11 {
			return 6, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
		}

		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q7(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 10, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok {
		if val.state != 11 {
			return 7, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
		}

		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 0, numberTokens, tokens

}

func (l Lexer) isOneState(lexeme string, states []int) (numberState, bool) {
	for _, state := range states {

		if l.numberDict[lexeme].state == state {
			return l.numberDict[lexeme], true
		}
	}

	return numberState{}, false
}

func (l Lexer) q8(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 9, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{6, 7, 11}); ok {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}

		return 8, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q9(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.isOneState(lexeme, []int{6, 7}); ok {
		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 9, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
}

func (l Lexer) q10(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.isOneState(lexeme, []int{6}); ok {
		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 10, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
}

func (l Lexer) q11(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 12, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{6, 7, 8}); ok {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}

		return 11, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q12(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{6, 7, 8}); ok {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}
		return 12, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) getNumberTokenFromList(numberTokens []Token) Token {
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
