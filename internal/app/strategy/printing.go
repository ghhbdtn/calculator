package strategy

import (
	"calculator/internal/app/domain"
	"errors"
)

type PrintingStrategy struct{}

func (s *PrintingStrategy) Execute(node *domain.Node) (interface{}, error) {
	if node == nil {
		return nil, errors.New("nil node")
	}

	<-node.Ready
	return domain.ResultItem{
		Var:   node.VarName,
		Value: node.Value,
	}, nil
}
