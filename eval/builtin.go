package eval

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hiroyaonoe/go-lisp/node"
)

var ErrInvalidArguments = errors.New("invalid arguments")

type Fun func(*Env, *node.Node) (*node.Node, error)

var builtin map[string]Fun

func init() {
	builtin = map[string]Fun{}
	builtin["+"] = wrapFun("+", doPlus)
	builtin["let"] = wrapFun("let", doLet)
	builtin["quote"] = wrapFun("quote", doQuote)
	builtin["cons"] = wrapFun("cons", doCons)
	builtin["defun"] = wrapFun("defun", doDefun)
	builtin["lambda"] = wrapFun("lambda", doLambda)
	builtin["princ"] = wrapFun("princ", doPrinc)
}

func wrapFun(name string, f Fun) Fun {
	return func(env *Env, n *node.Node) (*node.Node, error) {
		nn, err := f(env, n)
		if err != nil {
			err = fmt.Errorf("builtin %s: %w", name, err)
		}
		return nn, err
	}
}

func doPlus(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok {
		return nil, ErrInvalidArguments
	}
	ret := node.Int(0)
	for _, nn := range params {
		v, err := eval(env, nn)
		if err != nil {
			return nil, err
		}
		if node.NotIs(v, node.NodeInt) {
			return nil, ErrInvalidArguments
		}
		ret.Value = ret.Value.(int) + v.Value.(int)
	}
	return ret, nil
}

func doLet(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) < 2 {
		return nil, ErrInvalidArguments
	}
	kvs, ok := node.ListToNodes(params[0])
	if !ok {
		return nil, ErrInvalidArguments
	}
	kvmap := make(map[string]*node.Node, len(kvs))
	for _, kvlist := range kvs {
		kv, ok := node.ListToNodes(kvlist)
		if !ok || len(kv) != 2 {
			return nil, ErrInvalidArguments
		}
		k := kv[0]
		if node.NotIs(k, node.NodeSymbol) {
			return nil, ErrInvalidArguments
		}
		v, err := eval(env, kv[1])
		if err != nil {
			return nil, err
		}
		kvmap[k.Value.(string)] = v
	}

	lenv := NewEnv(env)
	for k, v := range kvmap {
		lenv.setVar(k, v)
	}
	bodys := params[1:]
	var ret *node.Node
	var err error
	for _, body := range bodys {
		ret, err = eval(lenv, body)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func doQuote(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 1 {
		return nil, ErrInvalidArguments
	}
	return params[0], nil
}

func doCons(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 2 {
		return nil, ErrInvalidArguments
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
		return nil, ErrInvalidArguments
	}

	if node.NotIs(params[0], node.NodeSymbol) {
		return nil, ErrInvalidArguments
	}
	name := params[0].Value.(string)

	args, ok := node.ListToNodes(params[1])
	if !ok {
		return nil, ErrInvalidArguments
	}
	for _, arg := range args {
		if node.NotIs(arg, node.NodeSymbol) {
			return nil, ErrInvalidArguments
		}
	}

	closure := node.Fun(env, params[1], params[2])
	global := env.Global()
	global.setFun(name, closure)

	return closure, nil
}

func doLambda(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 2 {
		return nil, ErrInvalidArguments
	}

	args, ok := node.ListToNodes(params[0])
	if !ok {
		return nil, ErrInvalidArguments
	}
	for _, arg := range args {
		if node.NotIs(arg, node.NodeSymbol) {
			return nil, ErrInvalidArguments
		}
	}

	closure := node.Fun(env, params[0], params[1])

	return closure, nil
}

func doPrinc(env *Env, n *node.Node) (*node.Node, error) {
	params, ok := node.ListToNodes(n)
	if !ok || len(params) != 1 {
		return nil, ErrInvalidArguments
	}
	v, err := eval(env, params[0])
	if err != nil {
		return nil, err
	}
	var s string
	if node.Is(v, node.NodeStr) {
		s = v.Value.(string)
	} else if node.Is(v, node.NodeInt) {
		i := v.Value.(int)
		s = strconv.Itoa(i)
	} else {
		return nil, ErrInvalidArguments
	}

	fmt.Print(s)
	return v, nil
}
