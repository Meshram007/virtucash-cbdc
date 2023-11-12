package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/virtucash-cbdc/chaincode/contracts"
)

func main() {
	cbdcContract := new(contracts.CBDCContract)

	chaincode, err := contractapi.NewChaincode(cbdcContract)

	if err != nil {
		log.Panicf("Could not create chaincode." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode. " + err.Error())
	}
}
