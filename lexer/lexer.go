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
		r, ok := l.peek()
		if !ok {
			return tokens, nil
		}
		if isWhite(r) {
			l.next()
			l.reduce()
			continue
		}
		switch r {
		case '(':
			tokens = append(tokens, token.LParen())
			l.next()
			l.reduce()
		case ')':
			tokens = append(tokens, token.RParen())
			l.next()
			l.reduce()
		case '+':
			tokens = append(tokens, token.Plus())
			l.next()
			l.reduce()
		default:
			token, ok := l.parseInt()
			if ok {
				tokens = append(tokens, token)
				continue
			}

			for {
				r, ok = l.peek()
				if !ok || isWhite(r) {
					break
				}
				l.next()
			}
			err := NewErrInvalidInput(l.reduce())
			l.reset()
			return tokens, err
		}
	}
}

func (l *lexer) parseInt() (token.Token, bool) {
	for {
		// ok == true
		r, _ := l.peek()
		if isNumber(r) {
			l.next()
		} else if isSymbol(r) {
			return token.Token{}, false
		} else {
			if l.pos == 0 {
				return token.Token{}, false
			}
			return token.Int(l.reduce()), true
		}
	}
}

func (l *lexer) next() {
	l.pos++
}

func (l *lexer) peek() (rune, bool) {
	if l.pos >= len(l.s) {
		return 0, false
	}
	return l.s[l.pos], true
}

func (l *lexer) reduce() string {
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

func (l *lexer) redo() {
	l.pos = 0
}

func (l *lexer) reset() {
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
