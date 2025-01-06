package spellnumber

import (
	"errors"
	"math/big"
	"testing"
)

func TestParserParse(t *testing.T) {
	tests := []struct {
		name          string
		input         []Token
		expected      *big.Int
		expectedError error
	}{
		{
			name:     "Simple number",
			input:    []Token{{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(100)}},
			expected: big.NewInt(100),
		},
		{
			name: "Multiple numbers",
			input: []Token{
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(100)},
				{Type: TOKEN_PLUS, Value: "+"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(20)},
			},
			expected: big.NewInt(120),
		},
		{
			name: "Fatorial",
			input: []Token{
				{Type: TOKEN_FACTORIAL, Value: "!"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(4)},
			},
			expected: big.NewInt(24),
		},
		{
			name: "100 * (20 + 10) = 3000",
			input: []Token{
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(100)},
				{Type: TOKEN_TIMES, Value: "*"},
				{Type: TOKEN_LEFT_BRACKET, Value: "("},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(20)},
				{Type: TOKEN_PLUS, Value: "+"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(10)},
				{Type: TOKEN_RIGHT_BRACKET, Value: ")"},
			},
			expected: big.NewInt(3000),
		},
		{
			name: "100 * (20 + 10) / 2 = 1500",
			input: []Token{
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(100)},
				{Type: TOKEN_TIMES, Value: "*"},
				{Type: TOKEN_LEFT_BRACKET, Value: "("},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(20)},
				{Type: TOKEN_PLUS, Value: "+"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(10)},
				{Type: TOKEN_RIGHT_BRACKET, Value: ")"},
				{Type: TOKEN_DIVIDE, Value: "/"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(2)},
			},
			expected: big.NewInt(1500),
		},
		{
			name: "3 ^ 9 = 19683",
			input: []Token{
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(3)},
				{Type: TOKEN_POWER, Value: "^"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(9)},
			},
			expected: big.NewInt(19683),
		},
		{
			name: "2 dividido por 2 = 1",
			input: []Token{
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(2)},
				{Type: TOKEN_DIVIDE, Value: "/"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(2)},
			},
			expected: big.NewInt(1),
		},
		{
			name: "2 mod 2 = 0",
			input: []Token{
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(2)},
				{Type: TOKEN_MOD, Value: "%"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(2)},
			},
			expected: big.NewInt(0),
		},
		{
			name: "parentese sem fechar",
			input: []Token{
				{Type: TOKEN_LEFT_BRACKET, Value: "("},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(2)},
			},
			expectedError: errors.New("Esperado fecha parentese(s)"),
		},
		{
			name: "operador sem numero",
			input: []Token{
				{Type: TOKEN_PLUS, Value: "+"},
			},
			expectedError: errors.New("Esperado um número"),
		},
		{
			name: "fatorial sem numero",
			input: []Token{
				{Type: TOKEN_FACTORIAL, Value: "!"},
			},
			expectedError: errors.New("Esperado um número"),
		},
		{
			name: "número negativo",
			input: []Token{
				{Type: TOKEN_MINUS, Value: "-"},
				{Type: TOKEN_NUMBER_PARSED, Number: big.NewInt(10)},
			},
			expected: big.NewInt(-10),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser(test.input)

			result, err := parser.Parse()

			if !(err == nil && test.expectedError == nil) && err.Error() != test.expectedError.Error() {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}

			if test.expected != nil && (result == nil || result.Cmp(test.expected) != 0) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
