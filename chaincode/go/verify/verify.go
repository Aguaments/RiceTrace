package main

import (
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
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

type Circuit struct {
	PreImage frontend.Variable
	Hash     frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	// hash function
	m, _ := mimc.NewMiMC(api)
	m.Write(circuit.PreImage)
	// specify constraints
	// mimc(preImage) == hash
	api.AssertIsEqual(circuit.Hash, m.Sum())
	return nil
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) WriteCRS(ctx contractapi.TransactionContextInterface, pk string, vk string) (string, error) {
	MSPid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("%s", "The msp id is not set.")
	}

	if MSPid != "Org1MSP" {
		return "", fmt.Errorf("%s", "Current ORG is not a supervision org.")
	}

	err = ctx.GetStub().PutState("ProofKey", []byte(pk))
	if err != nil {
		return "", err
	}
	err = ctx.GetStub().PutState("VerifyKey", []byte(vk))
	if err != nil {
		return "", err
	}
	return "Change state successfully!", nil
}

func (s *SmartContract) Verify(ctx contractapi.TransactionContextInterface, key string, proof string) (string, error) {
	// Acquire hash
	hash, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	vkValue, err := ctx.GetStub().GetState("VerifyKey")
	if err != nil {
		return "", err
	}

	var vk groth16.VerifyingKey
	var pf groth16.Proof
	_ = json.Unmarshal(vkValue, &vk)
	_ = json.Unmarshal([]byte(proof), &pf)

	// create public witness
	publicAssignment := &Circuit{
		Hash: frontend.Variable(hash),
	}
	publicWitness, err := frontend.NewWitness(publicAssignment, ecc.BN254)
	err = groth16.Verify(pf, vk, publicWitness)
	if err != nil {
		return "", fmt.Errorf("verification failed: %v\n", err)
	}

	id := key[9:]
	value, err := ctx.GetStub().GetState(id)
	var asset Asset = unmarshal(value)
	asset.State = "1"
	err = ctx.GetStub().PutState(asset.ID, marshal(asset))
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("verification succeded\n")
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
