package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hiroyaonoe/go-lisp/node"
	"github.com/hiroyaonoe/go-lisp/token"
)

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []token.Token
		want    []*node.Node
		wantErr error
	}{
		{
			name:   "int",
			tokens: []token.Token{token.Int("1")},
			want:   []*node.Node{node.Int(1)},
		},
		{
			name:   "symbol",
			tokens: []token.Token{token.Symbol("aaa")},
			want:   []*node.Node{node.Symbol("aaa")},
		},
		{
			name:   "t",
			tokens: []token.Token{token.Symbol("t")},
			want:   []*node.Node{node.T()},
		},
		{
			name:   "nil",
			tokens: []token.Token{token.Symbol("nil")},
			want:   []*node.Node{node.Nil()},
		},
		{
			name:   "要素数0のlist",
			tokens: []token.Token{token.LParen(), token.RParen()},
			want:   []*node.Node{node.Nil()},
		},
		{
			name:   "要素数1のlist",
			tokens: []token.Token{token.LParen(), token.Int("1"), token.RParen()},
			want:   []*node.Node{node.Cons(node.Int(1), node.Nil())},
		},
		{
			name:   "要素数2のlist",
			tokens: []token.Token{token.LParen(), token.Int("1"), token.Int("2"), token.RParen()},
			want:   []*node.Node{node.Cons(node.Int(1), node.Cons(node.Int(2), node.Nil()))},
		},
		{
			name:   "要素数3のlist",
			tokens: []token.Token{token.LParen(), token.Symbol("aaa"), token.Int("2"), token.Int("3"), token.RParen()},
			want:   []*node.Node{node.Cons(node.Symbol("aaa"), node.Cons(node.Int(2), node.Cons(node.Int(3), node.Nil())))},
		},
		{
			name:   "int2つ並列",
			tokens: []token.Token{token.Int("1"), token.Int("2")},
			want: []*node.Node{
				node.Int(1),
				node.Int(2),
			},
		},
		{
			name:   "list2つ並列",
			tokens: []token.Token{token.LParen(), token.Int("1"), token.Int("2"), token.Int("3"), token.RParen(), token.LParen(), token.Int("4"), token.Int("5"), token.Int("6"), token.RParen()},
			want: []*node.Node{
				node.Cons(node.Int(1), node.Cons(node.Int(2), node.Cons(node.Int(3), node.Nil()))),
				node.Cons(node.Int(4), node.Cons(node.Int(5), node.Cons(node.Int(6), node.Nil()))),
			},
		},
		{
			name:   "入れ子のlist",
			tokens: []token.Token{token.LParen(), token.Symbol("a"), token.LParen(), token.Int("1"), token.RParen(), token.LParen(), token.Int("2"), token.RParen(), token.Int("3"), token.RParen()},
			want: []*node.Node{
				node.Cons(
					node.Symbol("a"),
					node.Cons(
						node.Cons(node.Int(1), node.Nil()),
						node.Cons(
							node.Cons(node.Int(2), node.Nil()),
							node.Cons(
								node.Int(3),
								node.Nil(),
							),
						),
					),
				),
			},
		},
		{
			name:    "valueが数字でないint",
			tokens:  []token.Token{token.Int("a")},
			wantErr: NewErrInvalidToken(token.Int("a")),
		},
		{
			name:    "RParenが多い",
			tokens:  []token.Token{token.LParen(), token.Int("1"), token.RParen(), token.RParen()},
			wantErr: NewErrInvalidToken(token.RParen()),
		},
		{
			name:    "RParenが足りない",
			tokens:  []token.Token{token.LParen(), token.Int("1")},
			wantErr: ErrNeedNextTokens,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, err := p.Parse(tt.tokens)
			if err != nil {
				if diff := cmp.Diff(tt.wantErr.Error(), err.Error()); diff != "" {
					t.Errorf("err is mismatch (-want +got):\n%s", diff)
				}
			} else {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("Node value is mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
