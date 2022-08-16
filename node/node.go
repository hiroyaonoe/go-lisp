package node

type nodeType int

const (
	NodeInt nodeType = iota + 1
	NodeCons
	NodeNil
	NodeT
)

// NodeConsならValueはnil, それ以外ならCar, Cdrはnil
type Node struct {
	Type  nodeType
	Value any
	Car   *Node
	Cdr   *Node
}
