package strategy

import (
	"calculator/internal/app/domain"
	"errors"
)

type EvaluationStrategy struct{}

func (s *EvaluationStrategy) Execute(node *domain.Node) (interface{}, error) {
	if node == nil {
		return nil, errors.New("nil node")
	}

	if node.IsLiteral {
		return node.Value, nil
	}

	<-node.Left.Ready
	<-node.Right.Ready

	var result int64
	switch node.Op {
	case domain.Add:
		result = node.Left.Value + node.Right.Value
	case domain.Sub:
		result = node.Left.Value - node.Right.Value
	case domain.Mul:
		result = node.Left.Value * node.Right.Value
	default:
		return nil, errors.New("unknown operation")
	}

	node.Value = result
	close(node.Ready)
	return result, nil
}
