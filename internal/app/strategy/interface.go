package strategy

import "calculator/internal/app/domain"

type NodeStrategy interface {
	Execute(node *domain.Node) (interface{}, error)
}
