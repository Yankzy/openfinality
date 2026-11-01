package domain

type ParticipantType string

const (
	ParticipantTypeBank                ParticipantType = "BANK"
	ParticipantTypePSP                 ParticipantType = "PSP"
	ParticipantTypeMobileMoneyOperator ParticipantType = "MOBILE_MONEY_OPERATOR"
	ParticipantTypeClearingInstitution ParticipantType = "CLEARING_INSTITUTION"
	ParticipantTypeCentralBank         ParticipantType = "CENTRAL_BANK"
	ParticipantTypeSettlementOperator  ParticipantType = "SETTLEMENT_OPERATOR"
	ParticipantTypeLiquidityProvider   ParticipantType = "LIQUIDITY_PROVIDER"
	ParticipantTypeOther               ParticipantType = "OTHER"
)

type ParticipantStatus string

const (
	ParticipantStatusActive    ParticipantStatus = "ACTIVE"
	ParticipantStatusSuspended ParticipantStatus = "SUSPENDED"
	ParticipantStatusRevoked   ParticipantStatus = "REVOKED"
)

// Participant represents an institutional participant in the OpenFinality network.
type Participant struct {
	ID           string            `json:"id"`
	MSPID        string            `json:"msp_id"`
	Name         string            `json:"name"`
	Type         ParticipantType   `json:"type"`
	Jurisdiction string            `json:"jurisdiction"`
	Status       ParticipantStatus `json:"status"`
}
