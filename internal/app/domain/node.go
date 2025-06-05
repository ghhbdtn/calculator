package domain

import "sync"

type Node struct {
	VarName   string
	Left      *Node
	Right     *Node
	Op        Operation
	Value     int64
	IsLiteral bool
	Ready     chan struct{}
	mu        sync.Mutex
}

func NewLiteralNode(value int64) *Node {
	n := &Node{
		Value:     value,
		IsLiteral: true,
		Ready:     make(chan struct{}),
	}
	close(n.Ready)
	return n
}

func NewOperationNode(varName string, op Operation, left, right *Node) *Node {
	return &Node{
		VarName: varName,
		Op:      op,
		Left:    left,
		Right:   right,
		Ready:   make(chan struct{}),
	}
}
