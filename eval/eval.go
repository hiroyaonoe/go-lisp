package eval

import (
	"errors"

	"github.com/hiroyaonoe/go-lisp/node"
)

type Env struct {
	vars map[string]*node.Node
	funs map[string]*node.Node
	env  *Env
}

func NewEnv(env *Env) *Env {
	return &Env{
		vars: map[string]*node.Node{},
		funs: map[string]*node.Node{},
		env:  env,
	}
}

func (e *Env) Eval(n *node.Node) (*node.Node, error) {
	return e.eval(n)
}

func (e *Env) eval(n *node.Node) (*node.Node, error) {
	switch n.Type {
	case node.NodeCons:
		car := n.Car
		cdr := n.Cdr

		switch car.Type {
		case node.NodeSymbol:
			name := car.Value.(string)
			fn, ok := builtin[name]
			if ok {
				return fn(e, cdr)
			}

			return nil, errors.New("not implemented")
		default:
			return nil, errors.New("illegal function call")
		}
	default:
		return n, nil
	}
}
