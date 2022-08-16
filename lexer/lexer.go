package lexer

import (
	"fmt"

	"github.com/hiroyaonoe/go-lisp/token"
)

type lexer struct {
	s   []rune
	pos int
}

func NewLexer() *lexer {
	return &lexer{
		s:   []rune{},
		pos: 0,
	}
}

func (l *lexer) ReadString(s string) []token.Token {
	l.s = append(l.s, []rune(s)...)
	return l.parseTokens()
}

func (l *lexer) parseTokens() []token.Token {
	tokens := make([]token.Token, 0, len(l.s))
	for {
		r, ok := l.Next()
		fmt.Println(string(r), ok, tokens)
		if !ok {
			return tokens
		}
		if isWhite(r) {
			l.Reduce()
		}
		switch r {
		case '(':
			tokens = append(tokens, token.LParen())
			l.Reduce()
		case ')':
			tokens = append(tokens, token.RParen())
			l.Reduce()
		case '+':
			tokens = append(tokens, token.Plus())
			l.Reduce()
		default:
			token, ok := l.parseInt()
			if ok {
				tokens = append(tokens, token)
				continue
			}
		}
	}
}

func (l *lexer) parseInt() (token.Token, bool) {
	for {
		r, ok := l.Peek()
		if !ok {
			return token.Token{}, false
		}
		if isNumber(r) {
			l.Next()
		} else if isSymbol(r) {
			l.Reset()
			return token.Token{}, false
		} else {
			if l.pos == 0 {
				return token.Token{}, false
			}
			return token.Int(l.Reduce()), true
		}
	}
}

func (l *lexer) Next() (rune, bool) {
	if l.pos >= len(l.s) {
		return 0, false
	}
	l.pos++
	return l.s[l.pos-1], true
}

func (l *lexer) Peek() (rune, bool) {
	if l.pos >= len(l.s) {
		return 0, false
	}
	return l.s[l.pos], true
}

func (l *lexer) Reduce() string {
	if l.pos >= len(l.s) {
		return ""
	}
	v := string(l.s[:l.pos])
	l.s = l.s[l.pos:]
	l.pos = 0
	return v
}

func (l *lexer) Reset() {
	l.pos = 0
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isWhite(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r'
}

func isSymbol(r rune) bool {
	return (r >= 'a' && r <= 'z') || r == '-' || isNumber(r)
}
