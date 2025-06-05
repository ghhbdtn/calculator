package service

import (
	"calculator/internal/app/domain"
	"calculator/internal/app/strategy"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessInstructions(t *testing.T) {
	tests := []struct {
		name         string
		instructions []domain.Instruction
		expected     []domain.ResultItem
		expectedErr  error
	}{
		{
			name: "simple calculation with print",
			instructions: []domain.Instruction{
				{Type: domain.CalcInstruction, Var: "x", Op: domain.Add, Left: 2.0, Right: 3.0},
				{Type: domain.PrintInstruction, Var: "x"},
			},
			expected: []domain.ResultItem{
				{Var: "x", Value: 5},
			},
			expectedErr: nil,
		},
		{
			name: "multiple calculations",
			instructions: []domain.Instruction{
				{Type: domain.CalcInstruction, Var: "x", Op: domain.Add, Left: 2.0, Right: 3.0},
				{Type: domain.CalcInstruction, Var: "y", Op: domain.Mul, Left: "x", Right: 4.0},
				{Type: domain.PrintInstruction, Var: "y"},
			},
			expected: []domain.ResultItem{
				{Var: "y", Value: 20},
			},
			expectedErr: nil,
		},
		{
			name: "print non-existent variable",
			instructions: []domain.Instruction{
				{Type: domain.PrintInstruction, Var: "x"},
			},
			expected:    []domain.ResultItem(nil),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewCalculator(&strategy.EvaluationStrategy{})
			results, err := calc.ProcessInstructions(tt.instructions)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, results)
		})
	}
}

func TestProcessInstructions_StateReset(t *testing.T) {
	calc := NewCalculator(&strategy.EvaluationStrategy{})

	instructions1 := []domain.Instruction{
		{Type: domain.CalcInstruction, Var: "x", Op: domain.Add, Left: 2.0, Right: 3.0},
		{Type: domain.PrintInstruction, Var: "x"},
	}
	_, err := calc.ProcessInstructions(instructions1)
	assert.NoError(t, err)

	instructions2 := []domain.Instruction{
		{Type: domain.CalcInstruction, Var: "y", Op: domain.Mul, Left: 4.0, Right: 5.0},
		{Type: domain.PrintInstruction, Var: "y"},
	}
	results, err := calc.ProcessInstructions(instructions2)
	assert.NoError(t, err)
	assert.Equal(t, []domain.ResultItem{{Var: "y", Value: 20}}, results)

	assert.NotContains(t, calc.variables, "x")
	assert.NotContains(t, calc.printVars, "x")
}
