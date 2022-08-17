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
	switch n.Type {
	case NodeInt:
		return fmt.Sprintf("int: %v", n.Value)
	case NodeCons:
		return fmt.Sprintf(
			`cons: {
	%v,
	%v,
}`,
			strings.ReplaceAll(n.Car.String(), "\n", "\n\t"),
			strings.ReplaceAll(n.Cdr.String(), "\n", "\n\t"))
	case NodeNil:
		return "nil"
	case NodeSymbol:
		return fmt.Sprintf("symbol: %v", n.Value)
	default:
		return fmt.Sprintf("%v: %v, %v, %v", n.Type, n.Value, n.Car, n.Cdr)
	}
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
