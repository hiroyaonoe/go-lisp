package lexer

import (
	"github.com/hiroyaonoe/go-lisp/token"
)

type lexer struct {
	s string
	pos int
}

func NewLexer() *lexer {
	return &lexer{
		s: "",
		pos: 0,
	}
}

func (l *lexer) ReadString(s string) []token.Token {
	l.s += s
	return []token.Token{}
}
