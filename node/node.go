package node

import (
	"fmt"
	"strings"
)

type nodeType int

const (
	NodeInt nodeType = iota + 1
	NodeCons
	NodeNil
	NodeT
	NodeSymbol
)

// NodeConsならValueはnil, それ以外ならCar, Cdrはnil
type Node struct {
	Type  nodeType
	Value any
	Car   *Node
	Cdr   *Node
}

func (n *Node) String() string {
	l, ok := ListToNodes(n)
	if ok {
		if len(l) == 0 {
			return "nil"
		}
		ret := make([]string, len(l))
		for i, v := range l {
			ret[i] = v.String()
		}
		return fmt.Sprintf("(%v)", strings.Join(ret, " "))
	}
	if Is(n, NodeCons) {
		return fmt.Sprintf("(%v . %v)", n.Car, n.Cdr)
	}
	return fmt.Sprintf("%v", n.Value)
}

func ListToNodes(n *Node) ([]*Node, bool) {
	ret := make([]*Node, 0, 2)
	for {
		if Is(n, NodeNil) {
			return ret, true
		}
		if Is(n, NodeCons) {
			ret = append(ret, n.Car)
			n = n.Cdr
			continue
		}
		return nil, false
	}

}

func Is(n *Node, t nodeType) bool {
	return n != nil && n.Type == t
}

func NotIs(n *Node, t nodeType) bool {
	return !Is(n, t)
}

func Int(i int) *Node {
	return &Node{
		Type:  NodeInt,
		Value: i,
	}
}

func Cons(car *Node, cdr *Node) *Node {
	return &Node{
		Type: NodeCons,
		Car:  car,
		Cdr:  cdr,
	}
}

func Nil() *Node {
	return &Node{Type: NodeNil}
}

func T() *Node {
	return &Node{Type: NodeT}
}

func Symbol(v string) *Node {
	return &Node{
		Type:  NodeSymbol,
		Value: v,
	}
}

func List(ns ...*Node) *Node {
	ret := Nil()
	for i := len(ns) - 1; i >= 0; i-- {
		ret = &Node{
			Type: NodeCons,
			Car:  ns[i],
			Cdr:  ret,
		}
	}
	return ret
}
