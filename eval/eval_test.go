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
			name:    "not binded symbol",
			env:     NewEnv(nil),
			node:    node.Symbol("a"),
			wantErr: errors.New("not variable-binded symbol: a"),
		},
		{
			name: "引数0個のPlus",
			env:  NewEnv(nil),
			node: node.List(node.Symbol("+")),
			want: node.Int(0),
		},
		{
			name: "引数1個のPlus",
			env:  NewEnv(nil),
			node: node.List(node.Symbol("+"), node.Int(1)),
			want: node.Int(1),
		},
		{
			name: "引数2個のPlus",
			env:  NewEnv(nil),
			node: node.List(node.Symbol("+"), node.Int(1), node.Int(2)),
			want: node.Int(3),
		},
		{
			name: "引数3個のPlus",
			env:  NewEnv(nil),
			node: node.List(node.Symbol("+"), node.Int(1), node.Int(2), node.Int(3)),
			want: node.Int(6),
		},
		{
			name: "変数束縛",
			env:  NewEnv(nil).setVar("x", node.Int(1)),
			node: node.Symbol("x"),
			want: node.Int(1),
		},
		{
			name: "ローカル変数1個定義",
			env:  NewEnv(nil),
			node: node.List(
				node.Symbol("let"),
				node.List(
					node.List(
						node.Symbol("x"),
						node.Int(1),
					),
				),
				node.List(
					node.Symbol("+"),
					node.Int(2),
					node.Symbol("x"),
				),
			),
			want: node.Int(3),
		},
		{
			name: "ローカル変数2個定義",
			env:  NewEnv(nil),
			node: node.List(
				node.Symbol("let"),
				node.List(
					node.List(
						node.Symbol("x"),
						node.Int(1),
					),
					node.List(
						node.Symbol("y"),
						node.Int(2),
					),
				),
				node.List(
					node.Symbol("+"),
					node.Symbol("x"),
					node.Symbol("y"),
				),
			),
			want: node.Int(3),
		},
		{
			name:    "invalid arguments for +",
			env:     NewEnv(nil),
			node:    node.List(node.Symbol("+"), node.Nil()),
			wantErr: errors.New("builtin +: invalid arguments"),
		},
		{
			name:    "illegal function call",
			env:     NewEnv(nil),
			node:    node.List(node.Int(1), node.Int(2)),
			wantErr: errors.New("illegal function call"),
		},
		{
			name: "quote",
			env:  NewEnv(nil),
			node: node.List(
				node.Symbol("quote"),
				node.List(
					node.Symbol("+"),
					node.Int(1),
					node.Symbol("x"),
				),
			),
			want: node.List(
				node.Symbol("+"),
				node.Int(1),
				node.Symbol("x"),
			),
		},
		{
			name: "cons",
			env:  NewEnv(nil),
			node: node.List(
				node.Symbol("cons"),
				node.List(
					node.Symbol("+"),
					node.Int(1),
					node.Int(2),
				),
				node.List(
					node.Symbol("+"),
					node.Int(3),
					node.Int(4),
				),
			),
			want: node.Cons(
				node.Int(3),
				node.Int(7),
			),
		},
		{
			name: "defun",
			env:  NewEnv(nil),
			node: node.List(
				node.Symbol("defun"),
				node.Symbol("aa"),
				node.List(
					node.Symbol("x"),
					node.Symbol("y"),
				),
				node.List(
					node.Symbol("+"),
					node.Symbol("x"),
					node.Symbol("y"),
				),
			),
			want: setNodeFunToEnv(
				t,
				"aa",
				node.Fun(
					NewEnv(nil),
					node.List(
						node.Symbol("x"),
						node.Symbol("y"),
					),
					node.List(
						node.Symbol("+"),
						node.Symbol("x"),
						node.Symbol("y"),
					),
				),
			),
		},
		{
			name: "defun",
			env: setNodeFunToEnv(
				t,
				"aa",
				node.Fun(
					NewEnv(nil),
					node.List(
						node.Symbol("x"),
						node.Symbol("y"),
					),
					node.List(
						node.Symbol("+"),
						node.Symbol("x"),
						node.Symbol("y"),
					),
				),
			).Value.(*Env),
			node: node.List(
				node.Symbol("aa"),
				node.Int(2),
				node.Int(3),
			),
			want: node.Int(5),
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
				opt := cmp.AllowUnexported(Env{})
				if diff := cmp.Diff(tt.want, got, opt); diff != "" {
					t.Errorf("Node value is mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func setNodeFunToEnv(t *testing.T, s string, n *node.Node) *node.Node {
	t.Helper()
	env := n.Value.(*Env)
	env.setFun(s, n)
	return n
}
