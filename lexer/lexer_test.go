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
			want: []token.Token{token.LParen(), token.Plus(), token.Int("1"), token.Int("2"), token.RParen()},
		},
		{
			name: "(+ 123 456)",
			s:    "(+ 123 456)",
			want: []token.Token{token.LParen(), token.Plus(), token.Int("123"), token.Int("456"), token.RParen()},
		},
		{
			name: "  (  +  123  456  )  ",
			s:    "  (  +  123  456  )  ",
			want: []token.Token{token.LParen(), token.Plus(), token.Int("123"), token.Int("456"), token.RParen()},
		},
		{
			name:    "a",
			s:       "a",
			wantErr: NewErrInvalidInput("a"),
		},
		{
			name:    "( a ",
			s:       "( a ",
			wantErr: NewErrInvalidInput("a"),
		},
		{
			name:    "12a",
			s:       "12a",
			wantErr: NewErrInvalidInput("12a"),
		},
		{
			name:    "abc",
			s:       "abc",
			wantErr: NewErrInvalidInput("abc"),
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
