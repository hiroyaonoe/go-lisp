package parser

import (
	"errors"
	"fmt"
	"github.com/hiroyaonoe/go-lisp/token"
)

var (
	ErrNeedNextTokens = errors.New("must read next tokens")
)

type ErrInvalidToken struct {
	t token.Token
}

func NewErrInvalidToken(t token.Token) ErrInvalidToken {
	return ErrInvalidToken{t: t}
}

func (e ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token: '%v'", e.t)
}
