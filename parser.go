package spellnumber

import (
	"errors"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
)

type Parser struct {
	index   int
	tokens  []Token
	verbose bool
}

func NewParser(tokens []Token, verbose bool) *Parser {
	return &Parser{tokens: tokens, verbose: verbose}
}

func (p *Parser) Parse() (*big.Int, error) {
	if !p.verbose {
		log.SetOutput(io.Discard)

		defer log.SetOutput(os.Stdout)
	}

	if len(p.tokens) == 0 {
		return big.NewInt(0), nil
	}

	errorTokens := make([]string, 0, len(p.tokens))

	for _, token := range p.tokens {
		if token.Type == TOKEN_ERROR {
			errorTokens = append(errorTokens, token.Spell)
		}
	}

	if val := strings.Join(errorTokens, "; "); val != "" {
		return nil, errors.New(val)
	}

	return p.expression()
}

func (p *Parser) expression() (*big.Int, error) {
	sym := p.sym()

	if sym == TOKEN_PLUS || sym == TOKEN_MINUS {
		p.nextSym()
	}

	first, err := p.term()

	if err != nil {
		return nil, err
	}

	if sym == TOKEN_MINUS {
		first = first.Mul(big.NewInt(-1), first)
	}

	acceptedSymbols := map[TokenType]bool{
		TOKEN_PLUS:  true,
		TOKEN_MINUS: true,
	}

	for {
		if val, ok := acceptedSymbols[p.sym()]; !(ok && val) {
			break
		}

		op := p.sym()

		p.nextSym()

		second, err := p.term()

		if err != nil {
			return nil, err
		}

		if op == TOKEN_MINUS {
			first = big.NewInt(0).Sub(first, second)

			continue
		}

		if op == TOKEN_PLUS {
			first = big.NewInt(0).Add(first, second)

			continue

		}

		return nil, errors.New("Esperado um operador: mais ou menos")
	}

	return first, nil
}

func (p *Parser) term() (*big.Int, error) {
	first, err := p.factorial()

	if err != nil {
		return nil, err
	}

	acceptedSymbols := map[TokenType]bool{
		TOKEN_TIMES:  true,
		TOKEN_DIVIDE: true,
		TOKEN_POWER:  true,
		TOKEN_MOD:    true,
	}

	for {
		if val, ok := acceptedSymbols[p.sym()]; !(ok && val) {
			break
		}

		op := p.sym()

		p.nextSym()

		second, err := p.factorial()

		if err != nil {
			return nil, err
		}

		if op == TOKEN_TIMES {
			first = big.NewInt(1).Mul(first, second)
			continue
		}

		if op == TOKEN_DIVIDE {
			first = big.NewInt(1).Div(first, second)
			continue
		}

		if op == TOKEN_POWER {
			first = big.NewInt(1).Exp(first, second, nil)
			continue
		}

		if op == TOKEN_MOD {
			first = big.NewInt(1).Mod(first, second)
			continue
		}

		return nil, errors.New("Esperado um dos operadores: vezes, dividido por, elevado por, mod")
	}

	return first, nil
}

func (p *Parser) factorial() (*big.Int, error) {
	if p.sym() == TOKEN_FACTORIAL {
		p.nextSym()

		result, err := p.factorial()

		if err != nil {
			return nil, err
		}

		return big.NewInt(1).MulRange(1, result.Int64()), nil
	}

	return p.parenthesis()
}

func (p *Parser) parenthesis() (*big.Int, error) {
	sym := p.sym()

	if sym == TOKEN_LEFT_BRACKET {
		p.nextSym()

		exp, err := p.expression()

		if err != nil {
			return nil, err
		}

		if p.sym() != TOKEN_RIGHT_BRACKET {
			return nil, errors.New("Esperado fecha parentese(s)")
		}

		p.nextSym()

		return exp, nil
	}

	defer p.nextSym()

	return p.value()
}

func (p *Parser) value() (*big.Int, error) {
	if p.sym() != TOKEN_NUMBER_PARSED {
		return nil, errors.New("Esperado um número")
	}

	return p.tokens[p.index].Number, nil
}
func (p *Parser) sym() TokenType {
	if len(p.tokens) == 0 || p.index >= len(p.tokens) {
		return TOKEN_EOF
	}

	return p.tokens[p.index].Type
}

func (p *Parser) nextSym() {
	p.index++
}
