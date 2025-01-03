package spellnumber

import (
	"io"
	"log"
	"math/big"
	"os"
)

type Speller struct {
	thousands map[int][]string
	verbose   bool
}

func NewSpeller(verbose bool) *Speller {
	return &Speller{
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

	return "not implemented"
}
