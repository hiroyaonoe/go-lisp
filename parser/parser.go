package parser

import (
	"errors"

	"github.com/hiroyaonoe/go-lisp/token"
	"github.com/hiroyaonoe/go-lisp/node"
)

var (
	ErrNeedNextTokens = errors.New("must read next tokens")
)

type parser struct {
	tokens []token.Token
	pos int
}

func NewParser() *parser {
	return &parser{
		tokens: []token.Token{},
		pos: 0, 
	}
}

func (p *parser) Parse(tokens []token.Token) (*node.Node, error) {
	p.tokens = append(p.tokens, tokens...)
	return &node.Node{}, nil
}
