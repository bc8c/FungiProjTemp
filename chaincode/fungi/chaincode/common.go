package chaincode

import (	
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) _assetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) _getState(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}
	return assetJSON, nil
}

func (s *SmartContract) _updateOwnerFungusCount(ctx contractapi.TransactionContextInterface, clientId string, increment int ) error {
	countByte, err := s._getState(ctx, clientId)
	if countByte == nil {
		ctx.GetStub().PutState(clientId, []byte(strconv.Itoa(1)))
		return nil
	}
	if err !=nil {
		return err
	}
	ownerFungusCount, _ := strconv.Atoi(string(countByte[:]))
	ownerFungusCount += increment
	err = ctx.GetStub().PutState(clientId, []byte(strconv.Itoa(ownerFungusCount)))
	if err != nil {
		return fmt.Errorf("failed to put OwnerFungusCount state: %v", err)
	}
	return nil
}
