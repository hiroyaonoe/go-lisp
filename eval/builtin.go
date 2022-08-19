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
}

func doPlus(env *Env, n *node.Node) (*node.Node, error) {
	args, ok := node.ListToNodes(n)
	if !ok {
		return nil, errors.New("invalid arguments for +")
	}
	ret := node.Int(0)
	for _, nn := range args {
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
	args, ok := node.ListToNodes(n)
	if !ok || len(args) != 2 {
		return nil, ierr
	}
	kvs, ok := node.ListToNodes(args[0])
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
	body := args[1]
	if node.Is(body, node.NodeCons) {
		return eval(lenv, body)
	}
	return nil, ierr
}

func doQuote(env *Env, n *node.Node) (*node.Node, error) {
	args, ok := node.ListToNodes(n)
	if !ok || len(args) != 1 {
		return nil, errors.New("invalid arguments for quote")
	}
	return args[0], nil
}

func doCons(env *Env, n *node.Node) (*node.Node, error) {
	args, ok := node.ListToNodes(n)
	if !ok || len(args) != 2 {
		return nil, errors.New("invalid arguments for cons")
	}
	car, err := eval(env, args[0])
	if err != nil {
		return nil, err
	}
	cdr, err := eval(env, args[1])
	if err != nil {
		return nil, err
	}
	return node.Cons(car, cdr), nil
}
