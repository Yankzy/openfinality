package simulator

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/openfinality/openfinality/chaincode/afrorail/domain"
	"github.com/openfinality/openfinality/internal/settlement"
)

// SimulatorAdapter implements a deterministic external settlement simulator.
// It maintains simulated balances independently from Fabric ledger state.
type SimulatorAdapter struct {
	mu           sync.Mutex
	balances     map[string]map[string]domain.Amount // ParticipantID -> Currency -> Amount
	reservations map[string]*simulatorReservation
	executions   map[string]*simulatorExecution
}

type simulatorReservation struct {
	IdempotencyKey string
	ReservationID  string
	ParticipantID  string
	Amount         domain.Amount
	ValidUntil     time.Time
	Status         domain.ReservationStatus
}

type simulatorExecution struct {
	InstructionID string
	ExternalID    string
	Status        domain.SettlementStatus
}

func NewSimulatorAdapter() *SimulatorAdapter {
	return &SimulatorAdapter{
		balances:     make(map[string]map[string]domain.Amount),
		reservations: make(map[string]*simulatorReservation),
		executions:   make(map[string]*simulatorExecution),
	}
}

func (s *SimulatorAdapter) Name() string {
	return "AFRO_RAIL_DETERMINISTIC_SIMULATOR"
}

// SeedBalance adds an initial balance for a participant in the simulator.
func (s *SimulatorAdapter) SeedBalance(participantID, currency string, amount domain.Amount) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.balances[participantID]; !ok {
		s.balances[participantID] = make(map[string]domain.Amount)
	}
	s.balances[participantID][currency] = amount
}

func (s *SimulatorAdapter) Reserve(ctx context.Context, req settlement.ReserveRequest) (settlement.ReserveResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Idempotency check
	for _, res := range s.reservations {
		if res.IdempotencyKey == req.IdempotencyKey {
			// Return successful response if it's the exact same reservation
			return settlement.ReserveResult{
				Success:       true,
				ReservationID: res.ReservationID,
			}, nil
		}
	}

	// Check sufficient balance
	participantBalances, ok := s.balances[req.ParticipantID]
	if !ok {
		return settlement.ReserveResult{Success: false, ErrorMsg: "participant not found"}, nil
	}

	balance, ok := participantBalances[req.Amount.Currency]
	if !ok {
		return settlement.ReserveResult{Success: false, ErrorMsg: "currency balance not found"}, nil
	}

	cmp, err := balance.Cmp(req.Amount)
	if err != nil {
		return settlement.ReserveResult{Success: false, ErrorMsg: err.Error()}, nil
	}
	if cmp < 0 {
		return settlement.ReserveResult{Success: false, ErrorMsg: "insufficient funds"}, nil
	}

	// In a real system, we'd deduct or lock this balance.
	// We'll subtract from available balance.
	newBalance, err := balance.Sub(req.Amount)
	if err != nil {
		return settlement.ReserveResult{Success: false, ErrorMsg: err.Error()}, nil
	}
	s.balances[req.ParticipantID][req.Amount.Currency] = newBalance

	resID := fmt.Sprintf("res_%s", req.IdempotencyKey) // Simplistic ID generation

	s.reservations[resID] = &simulatorReservation{
		IdempotencyKey: req.IdempotencyKey,
		ReservationID:  resID,
		ParticipantID:  req.ParticipantID,
		Amount:         req.Amount,
		ValidUntil:     time.Now().Add(1 * time.Hour), // 1 hour expiry
		Status:         domain.ReservationStatusActive,
	}

	return settlement.ReserveResult{
		Success:       true,
		ReservationID: resID,
	}, nil
}

func (s *SimulatorAdapter) Execute(ctx context.Context, instr settlement.SettlementInstruction) (settlement.ExecutionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Idempotency check: duplicate instruction
	if exec, exists := s.executions[instr.InstructionID]; exists {
		return settlement.ExecutionResult{
			Success:    true,
			ExternalID: exec.ExternalID,
			Status:     exec.Status,
		}, nil
	}

	// Simplistic execution logic
	// Find payer reservations or balances, credit payee
	
	// Ensure payee exists
	if _, ok := s.balances[instr.PayeeInstitution]; !ok {
		s.balances[instr.PayeeInstitution] = make(map[string]domain.Amount)
	}
	
	payeeBalances := s.balances[instr.PayeeInstitution]
	
	currentPayeeBal, ok := payeeBalances[instr.Currency]
	if !ok {
		currentPayeeBal = domain.NewAmount(instr.Currency, 0, instr.Amount.Scale)
	}

	newBal, err := currentPayeeBal.Add(instr.Amount)
	if err != nil {
		return settlement.ExecutionResult{Success: false, ErrorMsg: err.Error()}, nil
	}
	s.balances[instr.PayeeInstitution][instr.Currency] = newBal

	extID := fmt.Sprintf("ext_%s", instr.InstructionID)
	status := domain.SettlementStatusSettled

	s.executions[instr.InstructionID] = &simulatorExecution{
		InstructionID: instr.InstructionID,
		ExternalID:    extID,
		Status:        status,
	}

	return settlement.ExecutionResult{
		Success:    true,
		ExternalID: extID,
		Status:     status,
	}, nil
}

func (s *SimulatorAdapter) Status(ctx context.Context, externalID string) (domain.SettlementStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, exec := range s.executions {
		if exec.ExternalID == externalID {
			return exec.Status, nil
		}
	}
	return "", errors.New("external execution not found")
}
