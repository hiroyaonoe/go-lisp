package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hiroyaonoe/go-lisp/token"
)

func Test_lexer_ReadString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    []token.Token
		wantErr error
	}{
		{
			name: "nil文字",
			s:    "",
			want: []token.Token{},
		},
		{
			name: "空白と改行",
			s:    " \n \r ",
			want: []token.Token{},
		},
		{
			name: "12",
			s:    "12",
			want: []token.Token{token.Int("12")},
		},
		{
			name: "(+ 1 2)",
			s:    "(+ 1 2)",
			want: []token.Token{token.LParen(), token.Symbol("+"), token.Int("1"), token.Int("2"), token.RParen()},
		},
		{
			name: "(+ 123 456)",
			s:    "(+ 123 456)",
			want: []token.Token{token.LParen(), token.Symbol("+"), token.Int("123"), token.Int("456"), token.RParen()},
		},
		{
			name: "  (  +  123  456  )  ",
			s:    "  (  +  123  456  )  ",
			want: []token.Token{token.LParen(), token.Symbol("+"), token.Int("123"), token.Int("456"), token.RParen()},
		},
		{
			name: "ab",
			s:    "ab",
			want: []token.Token{token.Symbol("ab")},
		},
		{
			name: "12ab34",
			s:    "12ab34",
			want: []token.Token{token.Symbol("12ab34")},
		},
		{
			name: "ab12cd",
			s:    "ab12cd",
			want: []token.Token{token.Symbol("ab12cd")},
		},
		{
			name: "a-b+c1d",
			s:    "a-b+c1d",
			want: []token.Token{token.Symbol("a-b+c1d")},
		},
		{
			name:    "あ",
			s:       "あ",
			wantErr: NewErrInvalidInput("あ"),
		},
		{
			name:    "( あ ",
			s:       "( あ ",
			wantErr: NewErrInvalidInput("あ"),
		},
		{
			name:    "12あ",
			s:       "12あ",
			wantErr: NewErrInvalidInput("あ"),
		},
		{
			name:    "aあ",
			s:       "aあ",
			wantErr: NewErrInvalidInput("あ"),
		},
		{
			name:    "あ12いう",
			s:       "あ12いう",
			wantErr: NewErrInvalidInput("あ12いう"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer()
			got, err := l.ReadString(tt.s)
			if err != nil {
				if diff := cmp.Diff(tt.wantErr.Error(), err.Error()); diff != "" {
					t.Errorf("err is mismatch (-want +got):\n%s", diff)
				}
			} else {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("[]Token value is mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
