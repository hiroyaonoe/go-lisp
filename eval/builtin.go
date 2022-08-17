package eval

import (
	"errors"
	"fmt"

	"github.com/hiroyaonoe/go-lisp/node"
)

type Fun func(*Env, *node.Node) (*node.Node, error)

var builtin map[string]Fun

func init() {
	builtin = map[string]Fun{}
	builtin["+"] = doPlus
	builtin["let"] = doLet
}

func doPlus(env *Env, n *node.Node) (*node.Node, error) {
	ret := node.Int(0)
	for {
		if n == nil || n.Type == node.NodeNil {
			return ret, nil
		}
		if n.Type != node.NodeCons {
			return nil, errors.New("invalid arguments for +")
		}
		v, err := env.eval(n.Car)
		if err != nil {
			return nil, err
		}
		if v == nil || v.Type != node.NodeInt {
			return nil, errors.New("invalid arguments for +")
		}
		ret.Value = ret.Value.(int) + v.Value.(int)
		n = n.Cdr
	}
}

func doLet(env *Env, n *node.Node) (*node.Node, error) {
	ierr := errors.New("invalid arguments for let")
	if n == nil || n.Type != node.NodeCons {
		return nil, ierr
	}
	kvs := n.Car
	kvmap := map[string]*node.Node{}
	for {
		if kvs == nil || kvs.Type == node.NodeNil {
			break
		}
		if kvs.Type != node.NodeCons {
			return nil, ierr
		}
		kv := kvs.Car
		if kv == nil || kv.Type != node.NodeCons {
			return nil, ierr
		}
		if kv.Car.Type != node.NodeSymbol {
			return nil, ierr
		}
		if kv.Cdr.Type != node.NodeCons {
			return nil, ierr
		}
		k := kv.Car.Value.(string)

		v, err := env.eval(kv.Cdr.Car)
		if err != nil {
			return nil, err
		}
		kvmap[k] = v
		kvs = kvs.Cdr
	}
	lenv := NewEnv(env)
	for k, v := range kvmap {
		fmt.Println(k, v)
		lenv.setVar(k, v)
	}
	if n.Cdr.Type != node.NodeCons {
		return nil, ierr
	}
	return lenv.eval(n.Cdr.Car)
}
