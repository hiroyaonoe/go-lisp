package eval

import "github.com/hiroyaonoe/go-lisp/node"

type hash struct {
	id    string
	value any
}

type Env []hash

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) Eval(ast *node.Node) (hash, error) {
	return hash{}, nil
}
