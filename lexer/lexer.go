package lexer

import (
	"errors"
	"strings"

	"github.com/hiroyaonoe/go-lisp/token"
)

var (
	EOF = errors.New("EOF")
)

type lexer struct {
	tokens []token.Token
	s      []rune
	pos    int
}

func NewLexer() *lexer {
	return &lexer{
		tokens: []token.Token{},
		s:      []rune{},
		pos:    0,
	}
}

func (l *lexer) ReadString(s string) ([]token.Token, error) {
	l.s = append(l.s, []rune(s)...)
	l.s = append(l.s, '\n')
	return l.parseTokens()
}

func (l *lexer) parseTokens() ([]token.Token, error) {
	for {
		r, ok := l.peek()
		if !ok {
			return l.reset(), nil
		}
		if isWhite(r) {
			l.next()
			l.reduce()
			continue
		}
		switch r {
		case '(':
			l.append(token.LParen())
			l.next()
			l.reduce()
		case ')':
			l.append(token.RParen())
			l.next()
			l.reduce()
		case '"':
			t, ok := l.parseStr()
			if !ok {
				return nil, EOF
			}
			l.append(t)
		default:
			token, ok := l.parseInt()
			if ok {
				l.append(token)
				continue
			}
			token, ok = l.parseSymbol()
			if ok {
				l.append(token)
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
			return l.reset(), err
		}
	}
}

func (l *lexer) parseStr() (token.Token, bool) {
	for {
		l.next()
		r, ok := l.peek()
		if !ok {
			l.redo()
			return token.Token{}, false
		}
		if r == '"' {
			l.next()
			return token.Str(l.reduce()), true
		}
	}
}

func (l *lexer) parseInt() (token.Token, bool) {
	for {
		// ok == true
		r, _ := l.peek()
		if isNumber(r) {
			l.next()
		} else if isSymbolLetter(r) {
			return token.Token{}, false
		} else {
			if l.pos == 0 {
				return token.Token{}, false
			}
			return token.Int(l.reduce()), true
		}
	}
}

func (l *lexer) parseSymbol() (token.Token, bool) {
	for {
		r, _ := l.peek()
		if isSymbolLetter(r) {
			l.next()
		} else {
			if l.pos == 0 {
				return token.Token{}, false
			}
			return token.Symbol(l.reduce()), true
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

func (l *lexer) reset() []token.Token {
	l.s = []rune{}
	l.pos = 0
	t := l.tokens
	l.tokens = []token.Token{}
	return t
}

func (l *lexer) append(ts ...token.Token) {
	l.tokens = append(l.tokens, ts...)
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isWhite(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r'
}

func isSymbolLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || isNumber(r) || strings.ContainsRune("-+", r)
}
