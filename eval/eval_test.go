package eval

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hiroyaonoe/go-lisp/node"
)

func TestEnv_Eval(t *testing.T) {
	tests := []struct {
		name    string
		env     *Env
		node    *node.Node
		want    *node.Node
		wantErr error
	}{
		{
			name: "int",
			env:  NewEnv(nil),
			node: node.Int(1),
			want: node.Int(1),
		},
		{
			name: "symbol",
			env:  NewEnv(nil),
			node: node.Symbol("a"),
			want: node.Symbol("a"),
		},
		{
			name: "引数0個のPlus",
			env:  NewEnv(nil),
			node: node.Cons(node.Symbol("+"), node.Nil()),
			want: node.Int(0),
		},
		{
			name: "引数1個のPlus",
			env:  NewEnv(nil),
			node: node.Cons(node.Symbol("+"), node.Cons(node.Int(1), node.Nil())),
			want: node.Int(1),
		},
		{
			name: "引数2個のPlus",
			env:  NewEnv(nil),
			node: node.Cons(node.Symbol("+"), node.Cons(node.Int(1), node.Cons(node.Int(2), node.Nil()))),
			want: node.Int(3),
		},
		{
			name: "引数3個のPlus",
			env:  NewEnv(nil),
			node: node.Cons(node.Symbol("+"), node.Cons(node.Int(1), node.Cons(node.Int(2), node.Cons(node.Int(3), node.Nil())))),
			want: node.Int(6),
		},
		{
			name:    "invalid arguments for +",
			env:     NewEnv(nil),
			node:    node.Cons(node.Symbol("+"), node.Cons(node.Symbol("a"), node.Nil())),
			wantErr: errors.New("invalid arguments for +"),
		},
		{
			name:    "illegal function call",
			env:     NewEnv(nil),
			node:    node.Cons(node.Int(1), node.Cons(node.Int(2), node.Nil())),
			wantErr: errors.New("illegal function call"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.env.Eval(tt.node)
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
