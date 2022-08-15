package token

type tokenType int

const (
	TokenInt tokenType = iota + 1
	TokenLParen
	TokenRParen
	TokenPlus
)

type Token struct {
	Type tokenType
	Value string
}
