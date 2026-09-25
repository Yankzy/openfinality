package contract

import (
	"encoding/json"
	"fmt"

	"github.com/openfinality/openfinality/chaincode/afrorail/domain"
	"github.com/hyperledger/fabric-chaincode-go/pkg/statebased"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

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

func (s *SmartContract) CreateSettlementIntent(ctx contractapi.TransactionContextInterface, intentID string, counterpartyMSPID string) error {
	existingBytes, err := ctx.GetStub().GetState("INTENT_" + intentID)
	if err != nil {
		return fmt.Errorf("failed to read world state: %v", err)
	}
	if existingBytes != nil {
		return fmt.Errorf("intent %s already exists", intentID)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP: %v", err)
	}

	intent := domain.SettlementIntent{
		ID:                intentID,
		InitiatorMSPID:    clientMSP,
		CounterpartyMSPID: counterpartyMSPID,
		Status:            domain.SettlementStatusCreated,
	}
	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("INTENT_"+intentID, intentBytes)
	if err != nil {
		return err
	}

	// Add state-based endorsement: Initiator AND Counterparty
	ep, err := statebased.NewStateEP(nil)
	if err != nil {
		return err
	}
	err = ep.AddOrgs(statebased.RoleTypeMember, clientMSP, counterpartyMSPID)
	if err != nil {
		return err
	}
	policyBytes, err := ep.Policy()
	if err != nil {
		return err
	}
	return ctx.GetStub().SetStateValidationParameter("INTENT_"+intentID, policyBytes)
}

func (s *SmartContract) AcceptSettlementIntent(ctx contractapi.TransactionContextInterface, intentID string) error {
	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP: %v", err)
	}

	intentBytes, err := ctx.GetStub().GetState("INTENT_" + intentID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if intentBytes == nil {
		return fmt.Errorf("the intent %s does not exist", intentID)
	}

	var intent domain.SettlementIntent
	err = json.Unmarshal(intentBytes, &intent)
	if err != nil {
		return err
	}

	// Idempotency
	if intent.Status == domain.SettlementStatusAccepted {
		return nil
	}

	if intent.Status != domain.SettlementStatusCreated {
		return fmt.Errorf("ERR_INVALID_STATE_TRANSITION: intent %s cannot be accepted in status %s", intentID, intent.Status)
	}

	// Authorization
	if clientMSP != intent.CounterpartyMSPID {
		return fmt.Errorf("ERR_UNAUTHORIZED_COUNTERPARTY: only %s can accept, but called by %s", intent.CounterpartyMSPID, clientMSP)
	}

	if clientMSP == intent.InitiatorMSPID {
		return fmt.Errorf("ERR_UNAUTHORIZED_COUNTERPARTY: initiator cannot accept its own intent")
	}

	intent.Status = domain.SettlementStatusAccepted

	intentBytes, err = json.Marshal(intent)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("INTENT_"+intentID, intentBytes)
}

func (s *SmartContract) ReserveCapacity(ctx contractapi.TransactionContextInterface, intentID string) error {
	// Ideally, ReserveCapacity should also check if caller is Initiator or Counterparty
	// For this milestone, we ensure only involved parties can reserve
	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP: %v", err)
	}

	intentBytes, err := ctx.GetStub().GetState("INTENT_" + intentID)
	if err != nil {
		return fmt.Errorf("failed to read world state: %v", err)
	}
	if intentBytes == nil {
		return fmt.Errorf("intent not found")
	}

	var intent domain.SettlementIntent
	err = json.Unmarshal(intentBytes, &intent)
	if err != nil {
		return err
	}

	if clientMSP != intent.InitiatorMSPID && clientMSP != intent.CounterpartyMSPID {
		return fmt.Errorf("ERR_UNAUTHORIZED_RESERVATION_OWNER: caller %s is not part of this settlement", clientMSP)
	}

	return s.updateIntentStatus(ctx, intentID, domain.SettlementStatusReserved)
}

func (s *SmartContract) EndorseInstruction(ctx contractapi.TransactionContextInterface, intentID string) error {
	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP: %v", err)
	}
	
	// We can let either involved party endorse
	intentBytes, err := ctx.GetStub().GetState("INTENT_" + intentID)
	if err != nil {
		return fmt.Errorf("failed to read world state: %v", err)
	}
	if intentBytes == nil {
		return fmt.Errorf("intent not found")
	}

	var intent domain.SettlementIntent
	err = json.Unmarshal(intentBytes, &intent)
	if err != nil {
		return err
	}

	if clientMSP != intent.InitiatorMSPID && clientMSP != intent.CounterpartyMSPID {
		return fmt.Errorf("ERR_UNAUTHORIZED_INITIATOR: caller %s is not part of this settlement", clientMSP)
	}
	return s.updateIntentStatus(ctx, intentID, domain.SettlementStatusSettlementPending)
}

func (s *SmartContract) CommitSettlementProof(ctx contractapi.TransactionContextInterface, intentID string) error {
	// Must be settlement adapter. For now, check if identity has "role=settlement-adapter" attribute.
	// Or we can just check if clientMSP == "OrdererMSP" if adapter runs in that org, 
	// but the prompt says: "Only an authorized settlement-adapter/service identity may record external settlement confirmation. Use a certificate attribute or explicitly registered participant role. For example: role=settlement-adapter"
	
	val, ok, err := ctx.GetClientIdentity().GetAttributeValue("role")
	if err != nil {
		return fmt.Errorf("failed to get client attributes: %v", err)
	}
	if !ok || val != "settlement-adapter" {
		return fmt.Errorf("ERR_UNAUTHORIZED_SETTLEMENT_ADAPTER: caller does not have role=settlement-adapter")
	}

	return s.updateIntentStatus(ctx, intentID, domain.SettlementStatusFinal)
}

func (s *SmartContract) QueryIntent(ctx contractapi.TransactionContextInterface, intentID string) (string, error) {
	intentBytes, err := ctx.GetStub().GetState("INTENT_" + intentID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if intentBytes == nil {
		return "", fmt.Errorf("the intent %s does not exist", intentID)
	}
	return string(intentBytes), nil
}

func (s *SmartContract) updateIntentStatus(ctx contractapi.TransactionContextInterface, intentID string, newStatus domain.SettlementStatus) error {
	intentBytes, err := ctx.GetStub().GetState("INTENT_" + intentID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if intentBytes == nil {
		return fmt.Errorf("the intent %s does not exist", intentID)
	}

	var intent domain.SettlementIntent
	err = json.Unmarshal(intentBytes, &intent)
	if err != nil {
		return err
	}

	if intent.Status == newStatus {
		return nil
	}

	intent.Status = newStatus

	intentBytes, err = json.Marshal(intent)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("INTENT_"+intentID, intentBytes)
}
