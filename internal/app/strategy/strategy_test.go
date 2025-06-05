package strategy

import (
	"calculator/internal/app/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluationStrategy_Execute(t *testing.T) {
	tests := []struct {
		name     string
		node     *domain.Node
		expected int64
		wantErr  bool
	}{
		{
			name:     "literal node",
			node:     domain.NewLiteralNode(5),
			expected: 5,
			wantErr:  false,
		},
		{
			name: "addition operation",
			node: &domain.Node{
				Op:    domain.Add,
				Left:  &domain.Node{Value: 2, Ready: make(chan struct{})},
				Right: &domain.Node{Value: 3, Ready: make(chan struct{})},
				Ready: make(chan struct{}),
			},
			expected: 5,
			wantErr:  false,
		},
		{
			name: "subtraction operation",
			node: &domain.Node{
				Op:    domain.Sub,
				Left:  &domain.Node{Value: 5, Ready: make(chan struct{})},
				Right: &domain.Node{Value: 3, Ready: make(chan struct{})},
				Ready: make(chan struct{}),
			},
			expected: 2,
			wantErr:  false,
		},
		{
			name: "multiplication operation",
			node: &domain.Node{
				Op:    domain.Mul,
				Left:  &domain.Node{Value: 2, Ready: make(chan struct{})},
				Right: &domain.Node{Value: 3, Ready: make(chan struct{})},
				Ready: make(chan struct{}),
			},
			expected: 6,
			wantErr:  false,
		},
		{
			name:     "nil node",
			node:     nil,
			expected: 0,
			wantErr:  true,
		},
		{
			name: "unknown operation",
			node: &domain.Node{
				Op:    "/",
				Left:  &domain.Node{Value: 2, Ready: make(chan struct{})},
				Right: &domain.Node{Value: 3, Ready: make(chan struct{})},
				Ready: make(chan struct{}),
			},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EvaluationStrategy{}

			// Подготовка каналов
			if tt.node != nil {
				if tt.node.Left != nil && tt.node.Left.Ready != nil {
					close(tt.node.Left.Ready)
				}
				if tt.node.Right != nil && tt.node.Right.Ready != nil {
					close(tt.node.Right.Ready)
				}
			}

			result, err := s.Execute(tt.node)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.(int64))
		})
	}
}

func TestPrintingStrategy(t *testing.T) {
	tests := []struct {
		name     string
		node     *domain.Node
		expected domain.ResultItem
		err      error
	}{
		{
			name: "valid node",
			node: &domain.Node{
				VarName: "x",
				Value:   10,
				Ready:   make(chan struct{}),
			},
			expected: domain.ResultItem{Var: "x", Value: 10},
			err:      nil,
		},
		{
			name:     "nil node",
			node:     nil,
			expected: domain.ResultItem{},
			err:      errors.New("nil node"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strat := &PrintingStrategy{}
			if tt.node != nil && tt.node.Ready != nil {
				close(tt.node.Ready)
			}

			result, err := strat.Execute(tt.node)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.expected, result.(domain.ResultItem))
			}
		})
	}
}
