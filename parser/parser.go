package parser

import (
	"strconv"

	"github.com/hiroyaonoe/go-lisp/node"
	"github.com/hiroyaonoe/go-lisp/token"
)

type parser struct {
	tokens []token.Token
	pos    int
}

func NewParser() *parser {
	return &parser{
		tokens: []token.Token{},
		pos:    0,
	}
}

func (p *parser) Parse(tokens []token.Token) ([]*node.Node, error) {
	p.tokens = append(p.tokens, tokens...)
	p.pos = 0
	nodes := []*node.Node{}
	for {
		if _, ok := p.peek(); !ok {
			break
		}
		n, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
		p.next()
	}
	p.reset()
	return nodes, nil
}

func (p *parser) parseNode() (*node.Node, error) {
	t, ok := p.peek()
	if !ok {
		return nil, ErrNeedNextTokens
	}
	switch t.Type {
	case token.TokenLParen:
		p.next()
		return p.parseParen()
	case token.TokenInt:
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			return nil, NewErrInvalidToken(t)
		}
		return node.Int(i), nil
	case token.TokenSymbol:
		return node.Symbol(t.Value), nil
	default:
		return nil, NewErrInvalidToken(t)
	}
}

func (p *parser) parseParen() (*node.Node, error) {
	t, ok := p.peek()
	if !ok {
		return nil, ErrNeedNextTokens
	}
	if t.Type == token.TokenRParen {
		return node.Nil(), nil
	} else {
		car, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		p.next()
		cdr, err := p.parseParen()
		if err != nil {
			return nil, err
		}
		return node.Cons(car, cdr), nil
	}
}

func (p *parser) next() {
	p.pos++
}

func (p *parser) peek() (token.Token, bool) {
	if p.pos >= len(p.tokens) {
		return token.Token{}, false
	}
	return p.tokens[p.pos], true
}

func (p *parser) reset() {
	p.tokens = []token.Token{}
	p.pos = 0
}
