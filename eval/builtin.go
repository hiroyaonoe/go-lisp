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
}

func doPlus(env *Env, n *node.Node) (*node.Node, error) {
	ret := node.Int(0)
	for {
		if n == nil || n.Type == node.NodeNil {
			return ret, nil
		}
		if n.Type != node.NodeCons ||
			(n.Car == nil || n.Car.Type != node.NodeInt) {
			return nil, errors.New("invalid arguments for +")
		}
		ret.Value = ret.Value.(int) + n.Car.Value.(int)
		n = n.Cdr
	}
}
