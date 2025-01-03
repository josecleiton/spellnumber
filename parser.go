package spellnumber

import (
	"io"
	"log"
	"math/big"
	"os"
)

type Parser struct {
	index   int
	tokens  []Token
	verbose bool
}

func NewParser(tokens []Token, verbose bool) *Parser {
	return &Parser{tokens: tokens, verbose: verbose}
}

func (p *Parser) Parse() *big.Int {
	if !p.verbose {
		log.SetOutput(io.Discard)

		defer log.SetOutput(os.Stdout)
	}

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
	sym := p.sym()
	if sym == TOKEN_PLUS || sym == TOKEN_MINUS {
		p.nextSym()
	}

	first := p.term()

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

		second := p.term()

		if op == TOKEN_MINUS {
			second = second.Mul(big.NewInt(-1), second)
		}

		first = big.NewInt(1).Add(first, second)
	}

	return first
}

func (p *Parser) term() *big.Int {
	first := p.factorial()

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

	defer p.nextSym()

	return p.value()
}

func (p *Parser) value() *big.Int {
	return p.tokens[p.index].Number
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
