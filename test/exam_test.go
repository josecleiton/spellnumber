package test

import (
	"fmt"
	"testing"

	"github.com/josecleiton/spellnumber"
)

func TestExam(t *testing.T) {
	expressions := []struct {
		input    string
		expected string
	}{
		{
			input:    "setenta e quatro mais abre parentese dez menos abre parentese cinco menos abre parentese abre parentese seis menos quatro fecha parentese mais um fecha parentese fecha parentese fecha parentese",
			expected: "oitenta e dois",
		},
		{
			input:    "dez",
			expected: "dez",
		},
		{
			input:    "abre parentese treze mais cinco vezes abre parentese cinco menos abre parentese um mais sete fecha parentese fecha parentese vezes quatro fecha parentese mais um decilhao vezes trinta e um mais um sextilhao vezes fatorial de cinco",
			expected: "trinta e um decilhoes cento e dezenove sextilhoes novecentos e noventa e nove quintilhoes novecentos e noventa e nove quatrilhoes novecentos e noventa e nove trilhoes novecentos e noventa e nove bilhoes novecentos e noventa e nove milhoes novecentos e noventa e nove mil e novecentos e cinquenta e tres",
		},
		{
			input:    "abre parentese trezentos setilhoes mais quatro trilhoes fecha parentese vezes oito",
			expected: "dois octilhoes quatrocentos setilhoes e trinta e dois trilhoes",
		},
		{
			input:    "fatorial de trinta",
			expected: "duzentos e sessenta e cinco nonilhoes duzentos e cinquenta e dois octilhoes oitocentos e cinquenta e nove setilhoes oitocentos e doze sextilhoes cento e noventa e um quintilhoes cinquenta e oito quatrilhoes seiscentos e trinta e seis trilhoes trezentos e oito bilhoes e quatrocentos e oitenta milhoes",
		},
		{
			input:    "um milhao trezentos e cinquenta e sete mil novecentos e sessenta e tres dividido por cinco mil setecentos e oitenta e nove",
			expected: "duzentos e trinta e quatro",
		},
		{
			input:    "fatorial de trinta vezes abre parentese fatorial de quarenta vezes abre parentese fatorial de oito mais quatro decilhoes fecha parentese fecha parentese",
			expected: "865695448983915109733797736765501176117285968845982834642736397897194603293761850826838779310899200000000000000000",
		},
		{
			input:    "duzentos elevado por dez",
			expected: "cento e dois sextilhoes e quatrocentos quintilhoes",
		},
		{
			input:    "trezentos e cinquenta e quatro mil setecentos e oitenta e nove dividido por trezentos e cinquenta e sete",
			expected: "novecentos e noventa e tres",
		},
		{
			input:    "vinte mil",
			expected: "vinte mil",
		},
	}

	for i, exp := range expressions {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {

			lexer := spellnumber.NewLexer(nil)

			tokens, err := lexer.ParseLine(exp.input)

			if err != nil {
				t.Errorf("Lexer Error: %v\n", err)
				return
			}

			parser := spellnumber.NewParser(tokens)

			result, err := parser.Parse()

			if err != nil {
				t.Errorf("Parser Error: %v\n", err)
				return
			}

			speller := spellnumber.NewSpeller()

			spelled := speller.Spell(result)

			if spelled != exp.expected {
				t.Errorf("Expected %s, got %s", exp.expected, spelled)
			}

		})
	}
}
