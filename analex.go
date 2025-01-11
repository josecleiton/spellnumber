package spellnumber

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
	verbose      bool
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
			"cem":             {state: 8, value: "100"},
			"cento":           {state: 9, value: "100"},
			"duzentos":        {state: 10, value: "200"},
			"trezentos":       {state: 10, value: "300"},
			"quatrocentos":    {state: 10, value: "400"},
			"quinhentos":      {state: 10, value: "500"},
			"seiscentos":      {state: 10, value: "600"},
			"setecentos":      {state: 10, value: "700"},
			"oitocentos":      {state: 10, value: "800"},
			"novecentos":      {state: 10, value: "900"},
			"mil":             {state: 13, value: "1000"},
			"milhao":          {state: 13, value: "1000000"},
			"milhoes":         {state: 13, value: "1000000"},
			"bilhao":          {state: 13, value: "1000000000"},
			"bilhoes":         {state: 13, value: "1000000000"},
			"trilhao":         {state: 13, value: "1000000000000"},
			"trilhoes":        {state: 13, value: "1000000000000"},
			"quadrilhao":      {state: 13, value: "1000000000000000"},
			"quadrilhoes":     {state: 13, value: "1000000000000000"},
			"quintilhao":      {state: 13, value: "1000000000000000000"},
			"quintilhoes":     {state: 13, value: "1000000000000000000"},
			"sextilhao":       {state: 13, value: "1000000000000000000000"},
			"sextilhoes":      {state: 13, value: "1000000000000000000000"},
			"septilhao":       {state: 13, value: "1000000000000000000000000"},
			"septilhoes":      {state: 13, value: "1000000000000000000000000"},
			"setilhao":        {state: 13, value: "1000000000000000000000000"},
			"setilhoes":       {state: 13, value: "1000000000000000000000000"},
			"octilhao":        {state: 13, value: "1000000000000000000000000000"},
			"octilhoes":       {state: 13, value: "1000000000000000000000000000"},
			"nonilhao":        {state: 13, value: "1000000000000000000000000000000"},
			"nonilhoes":       {state: 13, value: "1000000000000000000000000000000"},
			"decilhao":        {state: 13, value: "1000000000000000000000000000000000"},
			"decilhoes":       {state: 13, value: "1000000000000000000000000000000000"},
			"undecilhao":      {state: 13, value: "1000000000000000000000000000000000000"},
			"undecilhoes":     {state: 13, value: "1000000000000000000000000000000000000"},
			"duodecilhao":     {state: 13, value: "1000000000000000000000000000000000000000"},
			"duodecilhoes":    {state: 13, value: "1000000000000000000000000000000000000000"},
			"tridecilhao":     {state: 13, value: "1000000000000000000000000000000000000000000"},
			"tridecilhoes":    {state: 13, value: "1000000000000000000000000000000000000000000"},
			"quatradecilhao":  {state: 13, value: "1000000000000000000000000000000000000000000000"},
			"quatradecilhoes": {state: 13, value: "1000000000000000000000000000000000000000000000"},
			"zero":            {state: 15, value: "0"},
			"e":               {state: 200, value: "0"},
		},
	}
}

func (l *Lexer) SetVerbose(verbose bool) {
	l.verbose = verbose
}

func (l *Lexer) NextLine() ([]Token, error) {
	line, err := l.scannerStdIn.ReadString('\n')

	if err != nil {
		return []Token{}, err
	}

	line = strings.TrimSuffix(line, "\n")

	if line == "q" || line == "" {
		return []Token{}, nil
	}

	return l.ParseLine(line)
}

func (l *Lexer) ParseLine(rawLine string) ([]Token, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	line, _, err := transform.String(t, rawLine)

	if err != nil {
		return []Token{}, err
	}

	if !l.verbose {
		log.SetOutput(io.Discard)

		defer log.SetOutput(os.Stdout)
	}

	line = strings.Join(strings.Fields(strings.ToLower(line)), " ")

	tokens := make([]Token, 0, 64)

	words := strings.Split(line, " ")

	index := 0

	state := 0

	numberTokens := make([]Token, 0)
	for {
		lexeme := ""

		if index < len(words) {
			lexeme = words[index]
		}

		if lexeme == "" && state == 0 {
			break
		}

		if l.verbose {
			log.Printf("state: %d | lexeme: %s\n", state, lexeme)
		}

		if state == 0 {
			if len(numberTokens) > 0 {
				tokens = append(tokens, l.getNumberTokenFromList(numberTokens))

				numberTokens = make([]Token, 0, len(numberTokens)+1)
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

			if state == 0 {
				index--
			}
		} else if state == 10 {
			state, numberTokens, tokens = l.q10(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 11 {
			state, numberTokens, tokens = l.q11(lexeme, numberTokens, tokens)
		} else if state == 12 {
			state, numberTokens, tokens = l.q12(lexeme, numberTokens, tokens)
		} else if state == 13 {
			state, numberTokens, tokens = l.q13(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 14 {
			state, numberTokens, tokens = l.q14(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else if state == 15 {
			state, numberTokens, tokens = l.q15(lexeme, numberTokens, tokens)

			if state == 0 {
				index--
			}
		} else {
			tokens = append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: fmt.Sprintf("Lexema '%s' não reconhecido", lexeme)})
		}

		if len(tokens) > 0 && tokens[len(tokens)-1].Type == TOKEN_ERROR {
			break
		}

		index++
	}

	if len(numberTokens) > 0 {
		tokens = append(tokens, l.getNumberTokenFromList(numberTokens))

		numberTokens = []Token{}
	}

	return tokens, nil
}

func (l Lexer) q0(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "mais" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_PLUS, Value: "+"})
	}

	if lexeme == "menos" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_MINUS, Value: "-"})
	}

	if lexeme == "vezes" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_TIMES, Value: "*"})
	}

	if lexeme == "mod" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_MOD, Value: "%"})
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

	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{6, 7, 8, 9, 10, 15}); ok || val.value == "1000" {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: fmt.Sprintf("Lexema '%s' não reconhecido", lexeme)})
}

func (l Lexer) q1(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme != "por" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Expected 'por' after 'elevado'"})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_POWER, Value: "^"})
}

func (l Lexer) checkParenthesis(lexeme string) bool {
	regex := regexp.MustCompile(`^parenteses?$`)
	return regex.MatchString(lexeme)
}

func (l Lexer) q2(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if l.checkParenthesis(lexeme) {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_LEFT_BRACKET, Value: "("})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'parentese(s)' após 'abre'"})
}

func (l Lexer) q3(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if l.checkParenthesis(lexeme) {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_RIGHT_BRACKET, Value: ")"})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'parentese(s)' após 'fecha'"})
}

func (l Lexer) q4(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme != "de" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'de' após 'fatorial'"})

	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_FACTORIAL, Value: "!"})
}

func (l Lexer) q5(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme != "por" {
		return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'por' após 'dividido'"})
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_DIVIDE, Value: "/"})
}

func (l Lexer) q6(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.numberDict[lexeme]; ok {
		if val.state != 13 {
			return 6, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Não é esperado um número após '{unidade}'"})
		}

		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q7(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 12, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok {
		if val.state != 13 {
			return 7, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Não é esperado um número após '{dezena}'"})
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
	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{13}); ok {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}

		return 8, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Não é esperado U/D/C após 'cem'"})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q9(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 11, numberTokens, tokens
	}

	return 0, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'e' após 'cento'"})
}

func (l Lexer) q10(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 11, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{13}); ok {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}

		return 8, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'e' ou milhar após '{centena}'"})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q11(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.isOneState(lexeme, []int{6, 7}); ok {
		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 10, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado dezena ou unidade após '{centena} e'"})
}

func (l Lexer) q12(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.isOneState(lexeme, []int{6}); ok {
		return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
	}

	return 10, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado unidade após '{dezena} e'"})
}

func (l Lexer) q13(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if lexeme == "e" {
		return 14, numberTokens, tokens
	}

	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{6, 7, 8, 9, 10}); ok || val.value == "1000" {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}

		return 13, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado 'e' ou U/C/D depois de '{milhar}'"})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q14(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if val, ok := l.numberDict[lexeme]; ok {
		if _, ok := l.isOneState(lexeme, []int{6, 7, 8, 9, 10}); ok || val.value == "1000" {
			return val.state, append(numberTokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: lexeme}), tokens
		}
		return 14, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Esperado U/C/D depois de '{milhar} e'"})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) q15(lexeme string, numberTokens []Token, tokens []Token) (int, []Token, []Token) {
	if _, ok := l.numberDict[lexeme]; ok {
		return 15, numberTokens, append(tokens, Token{Type: TOKEN_ERROR, Value: lexeme, Spell: "Não esperado número após 'zero'"})
	}

	return 0, numberTokens, tokens
}

func (l Lexer) getNumberTokenFromList(numberTokens []Token) Token {
	if len(numberTokens) == 0 {
		return Token{Type: TOKEN_ERROR, Value: "0"}
	}

	log.Println(numberTokens)

	order := 1
	orderMilhar := len("1000")

	number := big.NewInt(0)

	for i := len(numberTokens) - 1; i >= 0; i-- {
		token := numberTokens[i]

		log.Println(token)

		tokenOrder := len(token.Value)

		if tokenOrder >= orderMilhar {
			if order > orderMilhar && tokenOrder <= order {
				return Token{Type: TOKEN_ERROR, Value: "0"}
			}

			if len(number.String()) < order {

				number = big.NewInt(0).Add(number, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(order-1)), nil))
			}

			order = tokenOrder

			continue
		}

		currentUnit := big.NewInt(0)

		currentUnit, ok := currentUnit.SetString(token.Value, 10)

		if !ok {
			return Token{Type: TOKEN_ERROR, Value: "0"}
		}

		exponent := big.NewInt(int64(order - 1))

		orderNumber := big.NewInt(1).Exp(big.NewInt(10), exponent, nil)

		currentNumber := currentUnit.Mul(currentUnit, orderNumber)

		number = number.Add(number, currentNumber)
	}

	// Adjust in case of 'mil' not prefixed by {unidade} | {dezena} | {centena}
	if order == 4 {
		milhar := big.NewInt(1000)

		if number.Cmp(milhar) == -1 {
			number = big.NewInt(0).Add(number, milhar)
		}
	}

	return Token{Type: TOKEN_NUMBER_PARSED, Value: number.String(), Number: number}
}
