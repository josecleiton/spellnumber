package spellnumber

import (
	"io"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
)

type Speller struct {
	thousands map[int][]string
	numbers   map[int]string
	and       string
	negative  string
	verbose   bool
}

func NewSpeller() *Speller {
	return &Speller{
		and:      "e",
		negative: "menos",
		numbers: map[int]string{
			-1:  "zero",
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
			0:  {"", ""},
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

func (s *Speller) SetVerbose(verbose bool) {
	s.verbose = verbose
}

func (s Speller) Spell(number *big.Int) string {
	if !s.verbose {
		log.SetOutput(io.Discard)

		defer log.SetOutput(os.Stdout)
	}

	negativeSign := ""

	if number.Sign() < 0 {
		negativeSign = s.negative + " "

		number = number.Abs(number)
	}

	numberStr := number.String()

	numberStrLen := len(numberStr)

	// Support until 10^49
	if numberStrLen > 49 {
		return numberStr
	}

	if numberStr == "0" {
		return s.numbers[-1]
	}

	builder := strings.Builder{}

	builder.WriteString(negativeSign)

	pluralIdx := 0

	addAnd := func(i int) {
		if i == 0 || i > len(numberStr)-1 {
			return
		}

		builder.WriteString(" ")
		builder.WriteString(s.and)
		builder.WriteString(" ")
	}

	for i := 0; i < len(numberStr); i++ {
		currentNumber := int(numberStr[i] - '0')

		numberPartIdx := (numberStrLen - i - 1) % 3

		order := (len(numberStr) - i - 1) / 3

		if numberPartIdx == 1 && currentNumber == 1 {
			currentNumber = 10 + int(numberStr[i+1]-'0')

			numberPartIdx = 0
			i++
		} else {
			currentNumber = currentNumber * int(math.Pow(10, float64(numberPartIdx)))
		}

		if currentNumber > 1 {
			pluralIdx = 1

			if numberPartIdx < 2 {
				addAnd(i)
			}
		}

		addedCurrentNumber := false

		if !(currentNumber == 1 && numberPartIdx == 0) {
			builder.WriteString(s.numbers[currentNumber])

			addedCurrentNumber = true
		}

		if numberPartIdx == 0 && order > 0 {

			if i < len(numberStr)-1 && addedCurrentNumber {
				builder.WriteString(" ")
			}

			builder.WriteString(s.thousands[order][pluralIdx])

			pluralIdx = 0

			addAnd(i)
		}
	}

	return builder.String()
}
