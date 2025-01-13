package spellnumber

import (
	"math/big"
	"testing"
)

func TestLexerParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Simple number",
			input: "cem",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "100"},
			},
		},
		{
			name:  "Multiple numbers",
			input: "cento e vinte",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "120"},
			},
		},
		{
			name:  "Fatorial",
			input: "fatorial de três",
			expected: []Token{
				{Type: TOKEN_FACTORIAL, Value: "!"},
				{Type: TOKEN_NUMBER_PARSED, Value: "3"},
			},
		},
		{
			name:  "Operador mais",
			input: "trinta e seis mais dois",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "36"},
				{Type: TOKEN_PLUS, Value: "+"},
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
			},
		},
		{
			name:  "Operador menos",
			input: "trinta e seis menos dois",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "36"},
				{Type: TOKEN_MINUS, Value: "-"},
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
			},
		},
		{
			name:  "Operador vezes",
			input: "trinta e seis vezes dois",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "36"},
				{Type: TOKEN_TIMES, Value: "*"},
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
			},
		},
		{
			name:  "Operador dividido",
			input: "trinta e seis dividido por dois",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "36"},
				{Type: TOKEN_DIVIDE, Value: "/"},
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
			},
		},
		{
			name:  "Operador mod",
			input: "trinta e seis mod dois",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "36"},
				{Type: TOKEN_MOD, Value: "%"},
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
			},
		},
		{
			name:  "centena composta",
			input: "duzentos e tres",
			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "203"},
			},
		},
		{
			name:  "Invalid input",
			input: "abc",
			expected: []Token{
				{Type: TOKEN_ERROR, Value: "abc", Spell: "Invalid input"},
			},
		},
		{
			name:  "Invalid 'cento'",
			input: "cento",
			expected: []Token{
				{Type: TOKEN_ERROR},
				{Type: TOKEN_NUMBER_PARSED, Value: "100"},
			},
		},
		{
			name:  "Parentese",
			input: "abre parentese dois fecha parentese",
			expected: []Token{
				{Type: TOKEN_LEFT_BRACKET, Value: "("},
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
				{Type: TOKEN_RIGHT_BRACKET, Value: ")"},
			},
		},
		{
			name:  "4 Decilhões",
			input: "quatro decilhoes",

			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: big.NewInt(1).Mul(big.NewInt(4), big.NewInt(1).Exp(big.NewInt(10), big.NewInt(33), nil)).String()},
			},
		},
		{
			name:  "Elevado",
			input: "dois elevado por quatro",

			expected: []Token{
				{Type: TOKEN_NUMBER_PARSED, Value: "2"},
				{Type: TOKEN_POWER, Value: "^"},
				{Type: TOKEN_NUMBER_PARSED, Value: "4"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(nil)
			tokens, err := lexer.ParseLine(test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(tokens) != len(test.expected) {
				t.Errorf("expected %d tokens, got %d", len(test.expected), len(tokens))
			}

			for i, token := range tokens {
				if token.Type != test.expected[i].Type {
					t.Errorf("expected token type %v, got %v", test.expected[i].Type, token.Type)

					if token.Type == TOKEN_ERROR {
						continue
					}
				}
				if token.Value != test.expected[i].Value {
					t.Errorf("expected token value %v, got %v", test.expected[i].Value, token.Value)
				}
			}
		})
	}
}
