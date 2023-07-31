package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"strconv"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	IpfsHash string            `json:"ipfshash"`
	Weight   int               `json:"weight"`
	Quality  string            `json:"quality"`
	Owner    string            `json:"owner"`
	Dist     string            `json:"dist"`
	State    string            `json:"state"`
	DateLine map[string]string `json:"dateline"`
	Trace    map[string]string `json:"trace"`
	Content  string            `json:"content"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, obj string) (string, error) {
	// Acquire current timestamp and mspid
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}
	creMSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}

	asset := unmarshal([]byte(obj))

	response := ctx.GetStub().InvokeChaincode("supervision", [][]byte{[]byte("Query"), []byte(creMSPid + "_" + asset.IpfsHash)}, "channel1")

	if string(response.Payload) != "1" {
		return "", fmt.Errorf("error: The product does not meet the requirements: %s, %s", creMSPid+"_"+asset.IpfsHash, response.Payload)
	}

	tm := timestamp.AsTime().Format(time.StampMilli)
	asset.Owner = creMSPid
	asset.Dist = ""
	asset.State = ""
	asset.DateLine = map[string]string{}
	asset.Trace = map[string]string{}
	asset.DateLine[tm] = creMSPid
	asset.Trace[asset.Owner] = fmt.Sprintf("TIMESTAMP:%s|ID:%s|NAME:%s|IPFSHASH:%s|WEIGHT:%d|QUALITY:%s|CONTENT:%s",
		tm, asset.ID, asset.Name, asset.IpfsHash, asset.Weight, asset.Quality, asset.Content)

	err = ctx.GetStub().PutState(asset.ID, marshal(asset))

	if err != nil {
		return "", err
	}
	return "Successfully. Asset :" + fmt.Sprintf("%v", asset), nil
}

func (s *SmartContract) UserTransferAsset(ctx contractapi.TransactionContextInterface, dist string, id string) (string, error) {
	//Acquire current time
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}
	//Acquire the current MSPid
	curMSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}

	if curMSPid == "Org1MSP" {
		return "", fmt.Errorf("%s", "Current mspid is SUPERVISION mspid. Not allow.\n")
	}

	value, err := ctx.GetStub().GetState(id)
	asset := unmarshal(value)
	ownMSPid := asset.Owner

	if curMSPid != ownMSPid {
		return "", fmt.Errorf("Current mspid is not match creator mspid. Current identity:%s, Creator identity:%s.\n", curMSPid, ownMSPid)
	}

	asset.Owner = "Org1MSP"
	asset.Dist = dist
	asset.DateLine[timestamp.AsTime().Format(time.StampMilli)] = "Org1MSP"

	err = ctx.GetStub().PutState(asset.ID, marshal(asset))
	if err != nil {
		return "", fmt.Errorf("error: %v. Current State: %v\n", err, asset)
	}
	return fmt.Sprintf("Successfully. Current asset :%v.\n", asset), nil
}

func (s *SmartContract) SuperTransferAsset(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	//Acquire current time
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}
	//Acquire the current MSPid
	curMSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}

	value, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}
	if value == nil {
		return "", fmt.Errorf("error: %s.\n", "The value for the key does not exist.")
	}

	asset := unmarshal(value)

	if curMSPid != "Org1MSP" {
		return "", fmt.Errorf("Current mspid is not match the SUPERVISION MSP. Current identity:%s.\n", curMSPid)
	}
	if asset.Owner != "Org1MSP" {
		return "", fmt.Errorf("The current asset not belong to the SUPERVISION MSP. Current identity:%s.\n", asset.Owner)
	}

	tm := timestamp.AsTime().Format(time.StampMilli)

	if asset.State != "1" {
		return "", fmt.Errorf("error: %v. Current State is not statisfied", err)
	}
	asset.Owner = asset.Dist
	asset.DateLine[tm] = asset.Dist
	asset.Trace[asset.Dist] = fmt.Sprintf("TIMESTAMP:%s|ID:%s|NAME:%s|IPFSHASH:%s|WEIGHT:%d|QUALITY:%s|CONTENT:%s",
		tm, asset.ID, asset.Name, asset.IpfsHash, asset.Weight, asset.Quality, asset.Content)
	asset.Dist = ""

	asset1 := asset

	err = ctx.GetStub().PutState(asset1.ID, marshal(asset1))
	if err != nil {
		return "", fmt.Errorf("error: %v. Current State: %v\n", err, asset1)
	}

	return fmt.Sprintf("Successfully. Current Asset: %v.\n", asset), nil
}

func (s *SmartContract) SelfTransfer(ctx contractapi.TransactionContextInterface, id string, cons string, content string) (string, error) {
	//Acquire current time
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}
	//Acquire the current MSPid
	curMSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}

	value, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}
	if value == nil {
		return "", fmt.Errorf("error: %s.\n", "The value for the key does not exist")
	}

	orgAsset := unmarshal(value)

	if orgAsset.Owner != curMSPid {
		return "", fmt.Errorf("error: The current asset does not belong to the contract caller. Asser Owner: %s. Current caller: %s", orgAsset.Owner, curMSPid)
	}

	newAsset := unmarshal([]byte(content))

	response := ctx.GetStub().InvokeChaincode("supervision", [][]byte{[]byte("Query"), []byte(curMSPid + "_" + newAsset.IpfsHash)}, "channel1")
	if string(response.Payload) != "1" {
		return "", fmt.Errorf("error: The product does not meet the requirements")
	}
	consumption, _ := strconv.Atoi(cons)

	if consumption > orgAsset.Weight && consumption < newAsset.Weight { //消耗大于原料量或者新产品的量多于消耗的量则自转移失败
		return "", fmt.Errorf("error: Consumption is greater than the amount of raw materials or the amount of new products is greater than the amount consumed. Consumption: %d. Original product weight: %d. New product weight: %d",
			consumption, orgAsset.Weight, newAsset.Weight)
	}

	//If the product is exhausted, delete the modified product
	newAsset.DateLine = map[string]string{}
	newAsset.Trace = map[string]string{}

	orgAsset.Weight -= consumption
	if orgAsset.Weight <= 0 {
		newAsset.DateLine = orgAsset.DateLine
		newAsset.Trace = orgAsset.Trace
		err = ctx.GetStub().DelState(orgAsset.ID)
		if err != nil {
			return "", fmt.Errorf("error: %v. Current State: %v\n", err, orgAsset)
		}
	}
	err = ctx.GetStub().PutState(orgAsset.ID, marshal(orgAsset))
	if err != nil {
		return "", fmt.Errorf("error: %v. Current State: %v\n", err, orgAsset)
	}

	//Deal with new products
	tm := timestamp.AsTime().Format(time.StampMilli)
	newAsset.Owner = curMSPid
	newAsset.Dist = ""
	newAsset.DateLine[tm] = curMSPid
	newAsset.Trace[newAsset.Owner] = fmt.Sprintf("TIMESTAMP:%s|ID:%s|NAME:%s|IPFSHASH:%s|WEIGHT:%d|QUALITY:%s|CONTENT:%s",
		tm, newAsset.ID, newAsset.Name, newAsset.IpfsHash, newAsset.Weight, newAsset.Quality, newAsset.Content)

	err = ctx.GetStub().PutState(newAsset.ID, marshal(newAsset))

	if err != nil {
		return "", err
	}
	return "Successfully. Asset :" + fmt.Sprintf("%v", newAsset), nil
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, hash string) (string, error) {

	value, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func unmarshal(b []byte) Asset {
	var a Asset
	if err := json.Unmarshal(b, &a); err != nil {
		log.Fatal(err)
	}
	return a
}

func marshal(a Asset) []byte {
	b, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func main() {
	con, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error to create smartcontract: %v", err)
	}

	if err = con.Start(); err != nil {
		log.Panicf("Error to start chaincode: %v.\n", err)
	}
}
