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
type Feed struct {
	FeedId		uint	`json:"feedid"`
	Name		string	`json:"name"`
	Dna			uint	`json:"dna"`
}

// Define Key names for options
const feedsCountKey = "feedsCount"

// init the chaincode 
func (s *SmartContract) Initialize (ctx contractapi.TransactionContextInterface) error {
	
	// Check authorization
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID : %v ", err)
	}
	if clientMSPID != "Org2MSP" {
		return fmt.Errorf("client is not auothorized to initalize fuongusCount : %v ", err)
	}

	// Check contract is not already set
	feedsCount, err := ctx.GetStub().GetState(feedsCountKey)
	if err != nil {
		return fmt.Errorf("failed to get fungusCount : %v ", err)
	}
	if feedsCount != nil {
		return fmt.Errorf("feedsCount is already set : %v ", err)
	}

	// Initilize feedsCount to zero(0)
	err = ctx.GetStub().PutState(feedsCountKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return fmt.Errorf("failed to set feedsCount : %v ", err)
	}
	return err
}

func (s *SmartContract) CreateRandomFeed (ctx contractapi.TransactionContextInterface, name string) error {
	// Check authorization
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID : %v ", err)
	}
	if clientMSPID != "Org2MSP" {
		return fmt.Errorf("client is not auothorized to initalize fuongusCount : %v ", err)
	}
	
	// 랜덤 DNA 생성
	dna := s._generateRandomDna(ctx,name)
	// 버섯 생성 및 원장에 저장하는 함수
	s._createFeed(ctx,name,dna)
	if err !=nil {
		return fmt.Errorf("failed to crateFeeds : %v", err)
	}
	return nil
}

func (s *SmartContract) _createFeed (ctx contractapi.TransactionContextInterface, name string, dna uint) error {

	feedsCountBytes, err := s._getState(ctx,feedsCountKey)
	if err != nil {
		return err
	}
	feedsIdINT, _ := strconv.Atoi(string(feedsCountBytes))
	feedsId := uint(feedsIdINT) 									
	
	feed := Feed{
		FeedId:		feedsId,
		Name:		name,
		Dna:		dna,
	}

	// marshal FeedsId
	assetJSON,  err := json.Marshal(feed)
	if err != nil {
		return fmt.Errorf("failed to marshal Feed : %v ", err)
	}
	// PutState FeedsId
	err = ctx.GetStub().PutState(strconv.Itoa(int(feedsId)), assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put feed : %v ", err)
	}
	
	// FeedsCount ++
	feedsId += 1
	err = ctx.GetStub().PutState(feedsCountKey, []byte(strconv.Itoa(int(feedsId))))
	if err != nil {
		return fmt.Errorf("failed to put feeds : %v ", err)
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

func (s *SmartContract) GetFeed (ctx contractapi.TransactionContextInterface, feedId uint) (*Feed, error) {

	feedsCountBytes, err := s._getState(ctx, strconv.Itoa(int(feedId)))
	if err != nil {
		return nil, fmt.Errorf("failed to get feed : %v",err)
	}

	var feed Feed
	err = json.Unmarshal(feedsCountBytes, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal feed : %v",err)
	}

	return &feed, nil
}
