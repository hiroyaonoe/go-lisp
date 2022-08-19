package eval

import (
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

func (e *Env) Eval(n *node.Node) (*node.Node, error) {
	return eval(e, n)
}

func (e *Env) String() string {
	ret := "env:\n"
	ret += "\tvars:\n"
	for k, v := range e.vars {
		ret += fmt.Sprintf("\t\t[%s] %s\n", k, strings.ReplaceAll(v.String(), "\n", "\n\t\t"))
	}
	ret += "\tfuns:\n"
	for k, v := range e.funs {
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

func (e *Env) Global() *Env {
	env := e
	for {
		if env.env == nil {
			return env
		}
		env = env.env
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

func (e *Env) getFun(s string) (*node.Node, bool) {
	env := e
	for {
		if env == nil {
			return nil, false
		}
		n, ok := env.funs[s]
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

func (e *Env) setFun(s string, n *node.Node) *Env {
	e.funs[s] = n
	return e
}
