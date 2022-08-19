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
			closure, ok := e.getFun(name)
			if ok {
				return call(closure, cdr)
			}
			v, ok := e.getVar(name)
			if ok {
				return call(v, cdr)
			}
			return nil, fmt.Errorf("not function-binded symbol: %s", name)
		default:
			v, err := eval(e, car)
			if err != nil {
				return nil, err
			}
			return call(v, cdr)
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

func call(f *node.Node, argsNode *node.Node) (*node.Node, error) {
	if node.NotIs(f, node.NodeFun) {
		return nil, errors.New("illegal function call")
	}

	scope := NewEnv(f.Value.(*Env))
	vars, _ := node.ListToNodes(f.Car)
	args, ok := node.ListToNodes(argsNode)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}
	num := len(vars)
	if len(args) != num {
		return nil, fmt.Errorf("missing the number of arguments")
	}

	for i := 0; i < num; i++ {
		scope.setVar(vars[i].Value.(string), args[i])
	}

	return eval(scope, f.Cdr)
}
