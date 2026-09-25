package settlement

import (
	"context"

	"github.com/openfinality/openfinality/chaincode/afrorail/domain"
)

type ReserveRequest struct {
	IdempotencyKey     string
	ReservationID      string
	SettlementIntentID string
	ParticipantID      string
	Amount             domain.Amount
}

type ReserveResult struct {
	Success       bool
	ReservationID string
	ErrorMsg      string
}

type SettlementInstruction struct {
	InstructionID      string
	SettlementIntentID string
	PayerInstitution   string
	PayeeInstitution   string
	Currency           string
	Amount             domain.Amount
	ExternalRail       string
	ExternalReferences []string
}

type ExecutionResult struct {
	Success            bool
	ExternalID         string
	Status             domain.SettlementStatus
	ErrorMsg           string
}

// SettlementAdapter defines the interface for external settlement coordination.
type SettlementAdapter interface {
	Name() string

	Reserve(
		ctx context.Context,
		request ReserveRequest,
	) (ReserveResult, error)

	Execute(
		ctx context.Context,
		instruction SettlementInstruction,
	) (ExecutionResult, error)

	Status(
		ctx context.Context,
		externalID string,
	) (domain.SettlementStatus, error)
}
