package spellnumber

import (
	"fmt"
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
	hundred   string
	hundreds  string
	verbose   bool
}

func NewSpeller() *Speller {
	return &Speller{
		and:      "e",
		negative: "menos",
		hundred:  "cem",
		hundreds: "cento",
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
			1:  {"mil", "mil"},
			2:  {"milhao", "milhoes"},
			3:  {"bilhao", "bilhoes"},
			4:  {"trilhao", "trilhoes"},
			5:  {"quatrilhao", "quatrilhoes"},
			6:  {"quintilhao", "quintilhoes"},
			7:  {"sextilhao", "sextilhoes"},
			8:  {"setilhao", "setilhoes"},
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

func (s Speller) formatNumberStr(numberStr string) string {
	length := len(numberStr)
	rest := length % 3

	if rest == 0 {
		return numberStr
	}

	formatStr := fmt.Sprintf("%s%d%s", "%0", length+3-rest, "s")

	return fmt.Sprintf(formatStr, numberStr)
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

	formattedNumber := s.formatNumberStr(numberStr)

	builder := strings.Builder{}

	builder.WriteString(negativeSign)

	addAnd := func() {
		builder.WriteString(" ")
		builder.WriteString(s.and)
		builder.WriteString(" ")
	}

	formattedNumberLen := len(formattedNumber)

	for i := 0; i < formattedNumberLen; i += 3 {
		nStr := formattedNumber[i : i+3]

		if nStr == "000" {
			continue
		}

		order := (formattedNumberLen - i - 1) / 3

		pluralIdx := 0

		if i > 0 {
			addAnd()
		}

		// mil
		if order == 1 && nStr == "001" {
			builder.WriteString(s.thousands[order][pluralIdx])
			continue
		}

		hadNumber := false
		for j := i; j < i+3; j++ {
			if formattedNumber[j] == '0' {
				continue
			}

			if j != i && hadNumber {
				addAnd()
			}

			n := int(formattedNumber[j]-'0') * int(math.Pow10(2-(j-i)))

			hadNumber = true

			// dez até dezenove
			if n == 10 {
				n = n + int(formattedNumber[j+1]-'0')

				j++
			}

			if n > 1 {
				pluralIdx = 1
			}

			if n == 100 {
				if strings.HasSuffix(nStr, "00") {
					builder.WriteString(s.hundred)
					break
				}

				builder.WriteString(s.hundreds)
				continue
			}

			builder.WriteString(s.numbers[n])
		}

		if order == 0 {
			continue
		}

		builder.WriteString(" ")
		builder.WriteString(s.thousands[order][pluralIdx])
	}

	return builder.String()
}
