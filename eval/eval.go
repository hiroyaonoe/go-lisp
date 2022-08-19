package eval

import (
	"errors"
	"fmt"

	"github.com/hiroyaonoe/go-lisp/node"
)

func eval(e *Env, n *node.Node) (*node.Node, error) {
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
			// TODO: lambda関数
			return nil, fmt.Errorf("not function-binded symbol: %s", name)
		default:
			return nil, errors.New("illegal function call")
		}
	case node.NodeSymbol:
		name := n.Value.(string)
		v, ok := e.getVar(name)
		if ok {
			return v, nil
		}
		return nil, fmt.Errorf("not variable-binded symbol: %s", name)
	default:
		return n, nil
	}
}
