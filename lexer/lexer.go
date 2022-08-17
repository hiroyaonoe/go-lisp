package lexer

import (
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

func (l *lexer) ReadString(s string) ([]token.Token, error) {
	l.s = append(l.s, []rune(s)...)
	l.s = append(l.s, ' ')
	return l.parseTokens()
}

func (l *lexer) parseTokens() ([]token.Token, error) {
	tokens := make([]token.Token, 0, len(l.s))
	for {
		r, ok := l.Peek()
		if !ok {
			return tokens, nil
		}
		if isWhite(r) {
			l.Next()
			l.Reduce()
			continue
		}
		switch r {
		case '(':
			tokens = append(tokens, token.LParen())
			l.Next()
			l.Reduce()
		case ')':
			tokens = append(tokens, token.RParen())
			l.Next()
			l.Reduce()
		case '+':
			tokens = append(tokens, token.Plus())
			l.Next()
			l.Reduce()
		default:
			token, ok := l.parseInt()
			if ok {
				tokens = append(tokens, token)
				continue
			}

			for {
				r, ok = l.Peek()
				if !ok || isWhite(r) {
					break
				}
				l.Next()
			}
			err := NewErrInvalidInput(l.Reduce())
			l.Reset()
			return tokens, err
		}
	}
}

func (l *lexer) parseInt() (token.Token, bool) {
	for {
		// ok == true
		r, _ := l.Peek()
		if isNumber(r) {
			l.Next()
		} else if isSymbol(r) {
			return token.Token{}, false
		} else {
			if l.pos == 0 {
				return token.Token{}, false
			}
			return token.Int(l.Reduce()), true
		}
	}
}

func (l *lexer) Next() {
	l.pos++
}

func (l *lexer) Peek() (rune, bool) {
	if l.pos >= len(l.s) {
		return 0, false
	}
	return l.s[l.pos], true
}

func (l *lexer) Reduce() string {
	if l.pos >= len(l.s) {
		l.s = []rune{}
		l.pos = 0
		return string(l.s)
	}
	v := string(l.s[:l.pos])
	l.s = l.s[l.pos:]
	l.pos = 0
	return v
}

func (l *lexer) Redo() {
	l.pos = 0
}

func (l *lexer) Reset() {
	l.s = []rune{}
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
