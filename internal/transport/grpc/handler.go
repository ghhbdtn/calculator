package grpc

import (
	"context"

	"calculator/internal/app/domain"
	"calculator/internal/app/service"
	"calculator/pkg/calculatorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	calculatorpb.UnimplementedCalculatorServiceServer
	service *service.Calculator
}

func NewHandler(service *service.Calculator) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Calculate(ctx context.Context, req *calculatorpb.CalculationRequest) (*calculatorpb.CalculationResponse, error) {
	instructions := make([]domain.Instruction, len(req.GetInstructions()))

	for i, inst := range req.GetInstructions() {
		instructions[i] = domain.Instruction{
			Type:  domain.InstructionType(inst.GetType()),
			Op:    domain.Operation(inst.GetOp()),
			Var:   inst.GetVar(),
			Left:  getValue(inst.GetLeft()),
			Right: getValue(inst.GetRight()),
		}
	}

	results, err := h.service.ProcessInstructions(instructions)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return h.toCalculationResponse(results), nil
}

func (h *Handler) toCalculationResponse(items []domain.ResultItem) *calculatorpb.CalculationResponse {
	response := &calculatorpb.CalculationResponse{
		Items: make([]*calculatorpb.ResultItem, len(items)),
	}

	for i, item := range items {
		response.Items[i] = &calculatorpb.ResultItem{
			Var:   item.Var,
			Value: item.Value,
		}
	}

	return response
}

func getValue(value interface{}) interface{} {
	switch v := value.(type) {
	case *calculatorpb.Instruction_LeftVal:
		return v.LeftVal
	case *calculatorpb.Instruction_LeftVar:
		return v.LeftVar
	case *calculatorpb.Instruction_RightVal:
		return v.RightVal
	case *calculatorpb.Instruction_RightVar:
		return v.RightVar
	default:
		return nil
	}
}
