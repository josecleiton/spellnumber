package spellnumber

import (
	"log"
	"math/big"
)

type Parser struct {
	index  int
	tokens []Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() *big.Int {
	if len(p.tokens) == 0 {
		return big.NewInt(0)
	}

	errorTokens := make([]Token, 0, len(p.tokens))

	for _, token := range p.tokens {
		if token.Type == TOKEN_ERROR {
			errorTokens = append(errorTokens, token)
		}
	}

	if len(errorTokens) > 0 {
		// better this
		log.Fatalln(errorTokens)
	}

	return p.expression()
}

func (p *Parser) expression() *big.Int {
	first := p.term()

	for p.sym() == TOKEN_PLUS || p.sym() == TOKEN_MINUS {
		op := p.sym()
		p.nextSym()

		second := p.term()

		if op == TOKEN_PLUS {
			first = big.NewInt(0).Add(first, second)

			continue
		}

		if op == TOKEN_MINUS {
			first = big.NewInt(0).Sub(first, second)

			continue
		}
	}

	return first
}

func (p *Parser) term() *big.Int {
	first := p.factorial()

	sym := p.sym()

	for sym == TOKEN_TIMES || sym == TOKEN_DIVIDE || sym == TOKEN_POWER || sym == TOKEN_MOD {
		op := sym

		p.nextSym()

		second := p.factorial()

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
	}

	return first
}

func (p *Parser) factorial() *big.Int {
	if p.sym() == TOKEN_FACTORIAL {
		p.nextSym()

		result := p.factorial()

		return big.NewInt(1).MulRange(1, result.Int64())
	}

	return p.parenthesis()
}

func (p *Parser) parenthesis() *big.Int {
	sym := p.sym()

	if sym == TOKEN_LEFT_BRACKET {
		p.nextSym()

		term := p.term()

		if p.sym() != TOKEN_RIGHT_BRACKET {
			log.Fatalln("Missing right bracket")
		}

		p.nextSym()

		return term
	}

	return p.value()
}

func (p *Parser) value() *big.Int {
	return p.tokens[p.index].Number
}
func (p *Parser) sym() TokenType {
	return p.tokens[p.index].Type
}

func (p *Parser) nextSym() {
	p.index++
}
