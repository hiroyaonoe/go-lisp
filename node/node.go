package node

type tokenType int

const (
	NodeInt = iota + 1
)
type Node struct {
	Type tokenType
	Value any
	Car *Node
	Cdr *Node
}
