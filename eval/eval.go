package eval

import (
	"errors"
	"fmt"
	"strings"

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

func (e *Env) String() string {
	ret := "env:\n"
	ret += "\tvars:\n"
	for k, v := range e.vars {
		ret += fmt.Sprintf("\t\t[%s] %s\n", k, strings.ReplaceAll(v.String(), "\n", "\n\t\t"))
	}
	if e.env == nil {
		ret += "\tenv: nil"
	} else {
		ret += "\t"
		ret += strings.ReplaceAll(e.env.String(), "\n", "\n\t")
	}
	return ret
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

func (e *Env) getVar(s string) (*node.Node, bool) {
	env := e
	for {
		if env == nil {
			return nil, false
		}
		n, ok := env.vars[s]
		if ok {
			return n, true
		}
		env = env.env
	}
}

func (e *Env) setVar(s string, n *node.Node) *Env {
	e.vars[s] = n
	return e
}
