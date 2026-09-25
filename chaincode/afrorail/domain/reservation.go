package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusPending  ReservationStatus = "PENDING"
	ReservationStatusActive   ReservationStatus = "ACTIVE"
	ReservationStatusConsumed ReservationStatus = "CONSUMED"
	ReservationStatusExpired  ReservationStatus = "EXPIRED"
	ReservationStatusReleased ReservationStatus = "RELEASED"
)

// Reservation represents an institution's cryptographically recorded commitment
// that specified settlement capacity has been reserved for an instruction.
type Reservation struct {
	ID                 string            `json:"id"`
	SettlementIntentID string            `json:"settlement_intent_id"`
	ParticipantID      string            `json:"participant_id"`
	Currency           string            `json:"currency"`
	Amount             Amount            `json:"amount"`
	ValidUntil         time.Time         `json:"valid_until"`
	Status             ReservationStatus `json:"status"`
}
