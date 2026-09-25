package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/openfinality/openfinality/chaincode/afrorail/contract"
)

func main() {
	smartContract := new(contract.SmartContract)

	chaincode, err := contractapi.NewChaincode(smartContract)
	if err != nil {
		log.Panicf("Error creating OpenFinality chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting OpenFinality chaincode: %v", err)
	}
}
