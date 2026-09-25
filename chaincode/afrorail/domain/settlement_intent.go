package domain

import "time"

type SettlementStatus string

const (
	SettlementStatusCreated           SettlementStatus = "CREATED"
	SettlementStatusValidated         SettlementStatus = "VALIDATED"
	SettlementStatusAccepted          SettlementStatus = "ACCEPTED"
	SettlementStatusReserved          SettlementStatus = "RESERVED"
	SettlementStatusSettlementPending SettlementStatus = "SETTLEMENT_PENDING"
	SettlementStatusSettled           SettlementStatus = "SETTLED"
	SettlementStatusFinal             SettlementStatus = "FINAL"
	
	// Failure states
	SettlementStatusRejected  SettlementStatus = "REJECTED"
	SettlementStatusExpired   SettlementStatus = "EXPIRED"
	SettlementStatusFailed    SettlementStatus = "FAILED"
	SettlementStatusCancelled SettlementStatus = "CANCELLED"
)

// SettlementIntent represents the canonical request submitted to OpenFinality.
type SettlementIntent struct {
	ID                  string           `json:"id"`
	InitiatorParticipantID    string           `json:"initiator_participant_id"`
	InitiatorMSPID            string           `json:"initiator_mspid"`
	CounterpartyParticipantID string           `json:"counterparty_participant_id"`
	CounterpartyMSPID         string           `json:"counterparty_mspid"`
	SourceCurrency            string           `json:"source_currency"`
	DestinationCurrency string           `json:"destination_currency"`
	SourceAmount        Amount           `json:"source_amount"`
	DestinationAmount   Amount           `json:"destination_amount"`
	PurposeCode         string           `json:"purpose_code"`
	ExternalReference   string           `json:"external_reference"`
	EvidenceHash        string           `json:"evidence_hash"`
	PolicyVersion       string           `json:"policy_version"`
	CreatedAt           time.Time        `json:"created_at"`
	ExpiresAt           time.Time        `json:"expires_at"`
	Status              SettlementStatus `json:"status"`
}
