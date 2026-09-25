package contract

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/afro-rail/afro-rail/chaincode/afrorail/domain"
)

// SmartContract provides functions for managing Afro-Rail settlement lifecycle
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of participants to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	participants := []domain.Participant{
		{ID: "NigeriaOrg", MSPID: "NigeriaMSP", Name: "NigeriaOrg Bank", Type: domain.ParticipantTypeBank, Jurisdiction: "NG", Status: domain.ParticipantStatusActive},
		{ID: "GhanaOrg", MSPID: "GhanaMSP", Name: "GhanaOrg Bank", Type: domain.ParticipantTypeBank, Jurisdiction: "GH", Status: domain.ParticipantStatusActive},
		{ID: "MoroccoOrg", MSPID: "MoroccoMSP", Name: "MoroccoOrg Bank", Type: domain.ParticipantTypeBank, Jurisdiction: "MA", Status: domain.ParticipantStatusActive},
	}

	for _, participant := range participants {
		participantJSON, err := json.Marshal(participant)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState("PARTICIPANT_"+participant.ID, participantJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateSettlementIntent creates a new settlement intent
func (s *SmartContract) CreateSettlementIntent(ctx contractapi.TransactionContextInterface, intentID string, intentJSON string) error {
	exists, err := s.IntentExists(ctx, intentID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the intent %s already exists", intentID)
	}

	var intent domain.SettlementIntent
	err = json.Unmarshal([]byte(intentJSON), &intent)
	if err != nil {
		return err
	}

	// Validate authorization and identity here...
	intent.Status = domain.SettlementStatusCreated

	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("INTENT_"+intentID, intentBytes)
}

// IntentExists returns true when intent with given ID exists in world state
func (s *SmartContract) IntentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	intentJSON, err := ctx.GetStub().GetState("INTENT_" + id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return intentJSON != nil, nil
}
