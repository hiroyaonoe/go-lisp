package eval

import (
	"errors"

	"github.com/hiroyaonoe/go-lisp/node"
)

type Fun func(*Env, *node.Node) (*node.Node, error)

var builtin map[string]Fun

func init() {
	builtin = map[string]Fun{}
	builtin["+"] = doPlus
	builtin["let"] = doLet
	builtin["quote"] = doQuote
	builtin["cons"] = doCons
	builtin["defun"] = doDefun
}

func doPlus(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok {
		return nil, errors.New("invalid arguments for +")
	}
	ret := node.Int(0)
	for _, nn := range params {
		v, err := eval(env, nn)
		if err != nil {
			return nil, err
		}
		if node.NotIs(v, node.NodeInt) {
			return nil, errors.New("invalid arguments for +")
		}
		ret.Value = ret.Value.(int) + v.Value.(int)
	}
	return ret, nil
}

func doLet(env *Env, n *node.Node) (*node.Node, error) {
	ierr := errors.New("invalid arguments for let")
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 2 {
		return nil, ierr
	}
	kvs, ok := node.ListToNodes(params[0])
	if !ok {
		return nil, ierr
	}
	kvmap := make(map[string]*node.Node, len(kvs))
	for _, kvlist := range kvs {
		kv, ok := node.ListToNodes(kvlist)
		if !ok || len(kv) != 2 {
			return nil, ierr
		}
		k := kv[0]
		if node.NotIs(k, node.NodeSymbol) {
			return nil, ierr
		}
		kvmap[k.Value.(string)] = kv[1]
	}

	lenv := NewEnv(env)
	for k, v := range kvmap {
		lenv.setVar(k, v)
	}
	body := params[1]
	if node.Is(body, node.NodeCons) {
		return eval(lenv, body)
	}
	return nil, ierr
}

func doQuote(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 1 {
		return nil, errors.New("invalid arguments for quote")
	}
	return params[0], nil
}

func doCons(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 2 {
		return nil, errors.New("invalid arguments for cons")
	}
	car, err := eval(env, params[0])
	if err != nil {
		return nil, err
	}
	cdr, err := eval(env, params[1])
	if err != nil {
		return nil, err
	}
	return node.Cons(car, cdr), nil
}

func doDefun(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 3 {
		return nil, errors.New("invalid arguments for cons")
	}

	if node.NotIs(params[0], node.NodeSymbol) {
		return nil, errors.New("invalid arguments for cons")
	}
	name := params[0].Value.(string)

	args, ok := node.ListToNodes(params[1])
	if !ok {
		return nil, errors.New("invalid arguments for cons")
	}
	for _, arg := range args {
		if node.NotIs(arg, node.NodeSymbol) {
			return nil, errors.New("invalid arguments for cons")
		}
	}

	closure := node.Fun(env, params[1], params[2])
	global := env.Global()
	global.setFun(name, closure)

	return closure, nil
}
