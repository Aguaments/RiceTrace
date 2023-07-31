package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"strconv"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) (string, error) {

	lastnum, err := ctx.GetStub().GetState("lastnum")

	if err != nil {
		return "", err
	}
	if lastnum != nil {
		return "The lastnum is not empty.", err
	}
	err = ctx.GetStub().PutState("lastnum", []byte("0"))
	if err != nil {
		return "Put state error", err
	}
	return "Initialization successful! Lastnum： 0", nil
}

func (s *SmartContract) IPFShashWrite(ctx contractapi.TransactionContextInterface, hash string, fileName string) (string, error) {
	//Acquire the current MSPid
	curMSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("error: %v.\n", err)
	}

	num, err := ctx.GetStub().GetState("lastnum")
	if err != nil {
		return "", err
	}
	if num == nil {
		return "The lastnum is empty.", nil
	}
	ctx.GetStub().PutState(curMSPid+"_"+hash, []byte("0")) //用于更改审核状态，便于后续的所有权转移工作

	value, _ := strconv.Atoi(string(num))
	value += 1
	sValue := strconv.Itoa(value)
	ctx.GetStub().PutState(sValue, []byte(curMSPid+"_"+hash+"_"+fileName)) //将该哈希值写入到链上，通过lastnum从后向前进行索引，来获取存在链上的hash值
	ctx.GetStub().PutState("lastnum", []byte(sValue))                      //存入hash值后将lastnum加1
	return "lastnum: " + sValue, nil
}

func (s *SmartContract) AcquireIPFSHash(ctx contractapi.TransactionContextInterface) (string, error) {

	MSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "The msp id is not set", err
	}

	if MSPid != "Org1MSP" {
		return "Current ORG is not a supervision org.", nil
	}

	num, err := ctx.GetStub().GetState("lastnum")

	if err != nil {
		return "", err
	}

	var hashList []string

	for nums, _ := strconv.Atoi(string(num)); nums > 0; nums-- {
		value, err := ctx.GetStub().GetState(strconv.Itoa(nums))
		if err != nil {
			return "", err
		}
		state, err := ctx.GetStub().GetState(string(value[:53]))
		if err != nil {
			return "", err
		}
		state_, err := strconv.Atoi(string(state))
		if state_ != 0 {
			break
		} else {
			hashList = append(hashList, string(value))
		}
	}
	return fmt.Sprintf("%s", hashList), nil
}

// 更改当前hash值对应的状态
func (s *SmartContract) ChangeHashState(ctx contractapi.TransactionContextInterface, orgMSPid string, hash string, state string) (string, error) {

	MSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("%s", "The msp id is not set.")
	}

	if MSPid != "Org1MSP" {
		return "", fmt.Errorf("%s", "Current ORG is not a supervision org.")
	}

	value, err := ctx.GetStub().GetState(orgMSPid + "_" + hash)
	if err != nil {
		return "", fmt.Errorf("error: %s, %s", err, orgMSPid+"_"+hash)
	}
	if value == nil {
		return "", fmt.Errorf("error: %s, %s", err, orgMSPid+"_"+hash)
	}

	err = ctx.GetStub().PutState(orgMSPid+"_"+hash, []byte(state))
	if err != nil {
		return "", err
	}
	return "Change state successfully!", nil

}

func (s *SmartContract) WriteVerifyHash(ctx contractapi.TransactionContextInterface, key string, hash string) (string, error) {
	MSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("%s", "The msp id is not set.")
	}

	if MSPid != "Org1MSP" {
		return "", fmt.Errorf("%s", "Current ORG is not a supervision org.")
	}

	err = ctx.GetStub().PutState(key, []byte(hash))
	if err != nil {
		return "", err
	}
	return "Verifying hash write successfully!", nil
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, hash string) (string, error) {

	value, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func main() {
	con, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error to create smartcontract: %v", err)
	}

	if err = con.Start(); err != nil {
		log.Panicf("Error to start chaincode: %v", err)
	}
}
