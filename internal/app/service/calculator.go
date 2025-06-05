package service

import (
	"calculator/internal/app/domain"
	"calculator/internal/app/strategy"
	"sync"
)

type Calculator struct {
	variables map[string]*domain.Node
	printVars map[string]bool
	mu        sync.RWMutex
	strategy  strategy.NodeStrategy
}

func NewCalculator(strategy strategy.NodeStrategy) *Calculator {
	return &Calculator{
		variables: make(map[string]*domain.Node),
		printVars: make(map[string]bool),
		strategy:  strategy,
	}
}

func (c *Calculator) SetStrategy(strategy strategy.NodeStrategy) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.strategy = strategy
}

func (c *Calculator) ProcessNode(node *domain.Node) (interface{}, error) {
	return c.strategy.Execute(node)
}
func (c *Calculator) buildDependencyTree(instructions []domain.Instruction) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, inst := range instructions {
		if inst.Type == domain.PrintInstruction {
			c.printVars[inst.Var] = true
		}
	}

	for _, inst := range instructions {
		if inst.Type == domain.CalcInstruction {
			if _, exists := c.variables[inst.Var]; exists {
				return domain.ErrVariableRedeclared
			}

			left, err := c.createOperand(inst.Left)
			if err != nil {
				return err
			}

			right, err := c.createOperand(inst.Right)
			if err != nil {
				return err
			}

			c.variables[inst.Var] = domain.NewOperationNode(inst.Var, inst.Op, left, right)
		}
	}
	return nil
}
func (c *Calculator) createOperand(input interface{}) (*domain.Node, error) {
	switch v := input.(type) {
	case float64:
		return domain.NewLiteralNode(int64(v)), nil
	case string:
		if node, exists := c.variables[v]; exists {
			return node, nil
		}
		return nil, domain.ErrVariableNotFound
	default:
		return nil, domain.ErrInvalidType
	}
}
func (c *Calculator) ProcessInstructions(instructions []domain.Instruction) ([]domain.ResultItem, error) {
	//сброс состояния
	c.mu.Lock()
	c.variables = make(map[string]*domain.Node)
	c.printVars = make(map[string]bool)
	c.mu.Unlock()

	if err := c.buildDependencyTree(instructions); err != nil {
		return nil, err
	}

	// Вычисляем все узлы
	c.SetStrategy(&strategy.EvaluationStrategy{})
	var wg sync.WaitGroup
	for _, node := range c.variables {
		wg.Add(1)
		go func(n *domain.Node) {
			defer wg.Done()
			c.ProcessNode(n)
		}(node)
	}
	wg.Wait()

	// Собираем результаты
	c.SetStrategy(&strategy.PrintingStrategy{})
	var results []domain.ResultItem
	for varName := range c.printVars {
		if node, exists := c.variables[varName]; exists {
			result, err := c.ProcessNode(node)
			if err != nil {
				return nil, err
			}
			results = append(results, result.(domain.ResultItem))
		}
	}

	return results, nil
}
