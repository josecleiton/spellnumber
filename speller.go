package spellnumber

import (
	"io"
	"log"
	"math/big"
	"os"
)

type Speller struct {
	thousands map[int][]string
	numbers   map[int]string
	and       string
	verbose   bool
}

func NewSpeller(verbose bool) *Speller {
	return &Speller{
		and: "e",
		numbers: map[int]string{
			0:   "zero",
			1:   "um",
			2:   "dois",
			3:   "tres",
			4:   "quatro",
			5:   "cinco",
			6:   "seis",
			7:   "sete",
			8:   "oito",
			9:   "nove",
			10:  "dez",
			11:  "onze",
			12:  "doze",
			13:  "treze",
			14:  "quatorze",
			15:  "quinze",
			16:  "dezesseis",
			17:  "dezessete",
			18:  "dezoito",
			19:  "dezenove",
			20:  "vinte",
			30:  "trinta",
			40:  "quarenta",
			50:  "cinquenta",
			60:  "sessenta",
			70:  "setenta",
			80:  "oitenta",
			90:  "noventa",
			100: "cem",
			200: "duzentos",
			300: "trezentos",
			400: "quatrocentos",
			500: "quinhentos",
			600: "seiscentos",
			700: "setecentos",
			800: "oitocentos",
			900: "novecentos",
		},
		thousands: map[int][]string{
			1:  {"mil"},
			2:  {"milhao", "milhoes"},
			3:  {"bilhao", "bilhoes"},
			4:  {"trilhao", "trilhoes"},
			5:  {"quadrilhao", "quadrilhoes"},
			6:  {"quintilhao", "quintilhoes"},
			7:  {"sextilhao", "sextilhoes"},
			8:  {"septilhao", "septilhoes"},
			9:  {"octilhao", "octilhoes"},
			10: {"nonilhao", "nonilhoes"},
			11: {"decilhao", "decilhoes"},
			12: {"undecilhao", "undecilhoes"},
			13: {"duodecilhao", "duodecilhoes"},
			14: {"tredecilhao", "tredecilhoes"},
			15: {"quatrodecilhao", "quatrodecilhoes"},
		},
	}

}

func (s Speller) Spell(number *big.Int) string {
	if !s.verbose {
		log.SetOutput(io.Discard)

		defer log.SetOutput(os.Stdout)
	}

	numberStr := number.String()

	// Support until 10^49
	if len(numberStr) > 49 {
		return numberStr
	}

	return "not implemented"
}
