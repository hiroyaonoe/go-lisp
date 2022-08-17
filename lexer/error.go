package lexer

import "fmt"

type ErrInvalidInput struct {
	l string
}

func NewErrInvalidInput(l string) ErrInvalidInput {
	return ErrInvalidInput{l: l}
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input: '%s'", e.l)
}
