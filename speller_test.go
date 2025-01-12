package spellnumber

import (
	"testing"

	"math/big"
)

func TestSpellerSpell(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{
			name:     "Zero",
			input:    big.NewInt(0),
			expected: "zero",
		},
		{
			name:     "One",
			input:    big.NewInt(1),
			expected: "um",
		},
		{
			name:     "Ten",
			input:    big.NewInt(10),
			expected: "dez",
		},
		{
			name:     "Hundred",
			input:    big.NewInt(100),
			expected: "cem",
		},
		{
			name:     "Hundred and ten",
			input:    big.NewInt(110),
			expected: "cento e dez",
		},
		{
			name:     "Thousand",
			input:    big.NewInt(1000),
			expected: "mil",
		},
		{
			name:     "Negative number",
			input:    big.NewInt(-16),
			expected: "menos dezesseis",
		},
		{
			name:     "Negative zero",
			input:    big.NewInt(1).Mul(big.NewInt(-1), big.NewInt(0)),
			expected: "zero",
		},
		{
			name:     "Complex number",
			input:    big.NewInt(32),
			expected: "trinta e dois",
		},
		{
			name:     "Large number",
			input:    big.NewInt(123456789),
			expected: "cento e vinte e tres milhoes quatrocentos e cinquenta e seis mil e setecentos e oitenta e nove",
		},
		{
			name:     "Stupendous number",
			input:    big.NewInt(1).MulRange(1, 100),
			expected: big.NewInt(1).MulRange(1, 100).String(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			speller := NewSpeller()
			result := speller.Spell(test.input)

			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
