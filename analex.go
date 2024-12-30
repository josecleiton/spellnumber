package spellnumber

import (
	"bufio"
	"os"
	"strings"
)

type TokenType int

type Token struct {
	Type  TokenType
	Value string
	Spell string
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

	TOKEN_NUMBER
)

type Lexer struct {
	Tokens       []Token
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
		Tokens:       make([]Token, 0, 1024),
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
			"e":               {state: 8, value: "0"},
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
		},
	}
}

func (l *Lexer) NextLine() {
	line, _ := l.scannerStdIn.ReadString('\n')

	words := strings.Split(line, " ")

	index := 0

	state := 0
	for {
		if index >= len(words) {
			break
		}

		w := words[index]

		if strings.Contains(w, "\n") {
			w = strings.ReplaceAll(w, "\n", "")
		}

		if state == 0 {
			if w == "mais" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_PLUS, Value: w})
			} else if w == "menos" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_MINUS, Value: w})
			} else if w == "vezes" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_TIMES, Value: w})
			} else if w == "elevado" {
				state = 1
			} else if w == "abre" {
				state = 2
			} else if w == "fecha" {
				state = 3
			} else if w == "fatorial" {
				state = 4
			} else if w == "dividido" {
				state = 5
			} else {
				if val, ok := l.numberDict[w]; ok {
					l.Tokens = append(l.Tokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: w})

					state = val.state
				} else {
					l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
					return
				}
			}
		} else if state == 1 {
			if w != "por" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
				return
			}

			l.Tokens = append(l.Tokens, Token{Type: TOKEN_POWER, Value: w})

			state = 0
		} else if state == 2 {
			if w != "parentese" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
				return
			}

			l.Tokens = append(l.Tokens, Token{Type: TOKEN_LEFT_BRACKET, Value: w})

			state = 0
		} else if state == 3 {
			if w != "parentese" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
				return
			}

			l.Tokens = append(l.Tokens, Token{Type: TOKEN_RIGHT_BRACKET, Value: w})

			state = 0
		} else if state == 4 {
			if w != "de" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
				return
			}

			l.Tokens = append(l.Tokens, Token{Type: TOKEN_FACTORIAL, Value: w})

			state = 0
		} else if state == 5 {
			if w != "por" {
				l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
				return
			}

			l.Tokens = append(l.Tokens, Token{Type: TOKEN_DIVIDE, Value: w})

			state = 0
		} else if state == 6 {
			if val, ok := l.numberDict[w]; ok {
				if val.state != 11 {
					l.Tokens = append(l.Tokens, Token{Type: TOKEN_ERROR, Value: w})
					return
				}

				l.Tokens = append(l.Tokens, Token{Type: TOKEN_NUMBER, Value: val.value, Spell: w})
			} else {
				index--
			}

			state = 0
		}

		index++
	}
}
