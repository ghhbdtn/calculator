package domain

type InstructionType string

const (
	CalcInstruction  InstructionType = "calc"
	PrintInstruction InstructionType = "print"
)

type Operation string

const (
	Add Operation = "+"
	Sub Operation = "-"
	Mul Operation = "*"
)

type Instruction struct {
	Type  InstructionType
	Op    Operation
	Var   string
	Left  interface{}
	Right interface{}
}
