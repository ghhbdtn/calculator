package domain

import "errors"

var (
	ErrVariableRedeclared = errors.New("variable redeclared")
	ErrVariableNotFound   = errors.New("variable not found")
	ErrInvalidType        = errors.New("invalid type")
	ErrInvalidOperation   = errors.New("invalid operation")
)
