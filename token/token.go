package token

type tokenType int

const (
	TokenInt tokenType = iota + 1
	TokenLParen
	TokenRParen
	TokenSymbol
)

type Token struct {
	Type  tokenType
	Value string
}

func (t *Token) String() string {
	return t.Value
}

func Int(s string) Token {
	return Token{
		Type:  TokenInt,
		Value: s,
	}
}

func LParen() Token {
	return Token{
		Type:  TokenLParen,
		Value: "(",
	}
}
func RParen() Token {
	return Token{
		Type:  TokenRParen,
		Value: ")",
	}
}

func Symbol(s string) Token {
	return Token{
		Type:  TokenSymbol,
		Value: s,
	}
}
