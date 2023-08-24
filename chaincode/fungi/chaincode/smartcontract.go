
package chaincode

import (
	"strconv"
	"fmt"
	"time"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct{
	contractapi.Contract
}

// Fungus Asset describes basic details
type Fungus struct {
	FungusId	uint	`json:"fungusid"`
	Name		string	`json:"name"`
	Owner		string	`json:"owner"`
	Dna			uint	`json:"dna"`
	ReadyTime	uint32	`json:"readytime"`
}

// Define Key names for options
const fungusCountKey = "FungusCount"

// init the chaincode 
func (s *SmartContract) Initialize (ctx contractapi.TransactionContextInterface) error {
	
	// Check authorization
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID : %v ", err)
	}
	if clientMSPID != "Org1MSP" {
		return fmt.Errorf("client is not auothorized to initalize fuongusCount : %v ", err)
	}

	// Check contract is not already set
	fungusCount, err := ctx.GetStub().GetState(fungusCountKey)
	if err != nil {
		return fmt.Errorf("failed to get fungusCount : %v ", err)
	}
	if fungusCount != nil {
		return fmt.Errorf("fungusCount is already set : %v ", err)
	}

	// Initilize FungusCountKey to zero(0)
	err = ctx.GetStub().PutState(fungusCountKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return fmt.Errorf("failed to set fungusCount : %v ", err)
	}
	return err
}

func (s *SmartContract) CreateRandomFungus (ctx contractapi.TransactionContextInterface, name string) error {
	
	// Check ClientID
	clientId, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get ClientID : %v ", err)
	}
	exists, err := s._assetExists(ctx, clientId)
	if err != nil{
		return err
	}
	if exists {
		return fmt.Errorf("client has already created an inital fungus")
	}

	// 랜덤 DNA 생성
	dna := s._generateRandomDna(ctx,name)
	// 버섯 생성 및 원장에 저장하는 함수
	s._createFungus(ctx,name,dna)
	return nil
}

func (s *SmartContract) _createFungus (ctx contractapi.TransactionContextInterface, name string, dna uint) error {

	fungusCountBytes, err := s._getState(ctx,fungusCountKey)
	if err != nil {
		return err
	}
	fungusIdINT, _ := strconv.Atoi(string(fungusCountBytes))
	fungusId := uint(fungusIdINT) 									

	// Check ClientID
	clientId, err := ctx.GetClientIdentity().GetID()				
	if err != nil {
		return fmt.Errorf("failed to get ClientID : %v ", err)
	}

	nowTime := time.Now().Unix()
	readyTime := nowTime + 60										
	
	fungus := Fungus{
		FungusId:	fungusId,
		Name:		name,
		Owner:		clientId,
		Dna:		dna,
		ReadyTime:	uint32(readyTime),
	}

	// marshal FungusId
	assetJSON,  err := json.Marshal(fungus)
	if err != nil {
		return fmt.Errorf("failed to marshal fungus : %v ", err)
	}
	// PutState FungusId
	err = ctx.GetStub().PutState(strconv.Itoa(int(fungusId)), assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put fungus : %v ", err)
	}
	
	// fungusCount ++
	fungusId += 1
	err = ctx.GetStub().PutState(fungusCountKey, []byte(strconv.Itoa(int(fungusId))))
	if err != nil {
		return fmt.Errorf("failed to put fungus : %v ", err)
	}
	// ownerFungusCount ++ 
	err = s._updateOwnerFungusCount(ctx,clientId, 1)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) _generateRandomDna (ctx contractapi.TransactionContextInterface, name string) uint {
	
	unixTime := time.Now().Unix()
	data := strconv.Itoa(int(unixTime)) + name
	hash := sha256.New()
	hash.Write([]byte(data))
	dnaHash := uint(binary.BigEndian.Uint64(hash.Sum(nil)))

	// make 14digits dna
	dna := dnaHash % uint(math.Pow(10, 10))
	dna = dna - (dna % 100)

	return dna	
}

func (s *SmartContract) GetFungiByOwner (ctx contractapi.TransactionContextInterface) ([]*Fungus, error) {

	// ClientId를 기준으로 소유한 모든 버섯정보를 조회해와서 반환한다.
	// Check ClientID
	clientId, err := ctx.GetClientIdentity().GetID()				
	if err != nil {
		return nil, fmt.Errorf("failed to get ClientID : %v ", err)
	}

	queryString := fmt.Sprintf(`{"selector":{"owner":"%s"}}`,clientId)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil,fmt.Errorf("failed to getQueryResult : %v ", err)
	}
	defer resultsIterator.Close()

	var fungi []*Fungus

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var fungus Fungus
		err = json.Unmarshal(queryResult.Value, &fungus)
		if err != nil {
			return nil, err
		}
		fungi = append(fungi, &fungus)
	}

	return fungi, nil
}