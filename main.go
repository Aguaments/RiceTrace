package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"network/sdkinit"
	"os"
)

var App sdkInit.Application

func main() {

	// Data
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "Admin",
			OrgPeerNum:    1,
			OrgAnchorFile: "/home/alex/go/src/github.com/Aguaments/network/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "Admin",
			OrgPeerNum:    1,
			OrgAnchorFile: "/home/alex/go/src/github.com/Aguaments/network/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "Admin",
			OrgPeerNum:    1,
			OrgAnchorFile: "/home/alex/go/src/github.com/Aguaments/network/fixtures/channel-artifacts/Org3MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "Admin",
			OrgPeerNum:    1,
			OrgAnchorFile: "/home/alex/go/src/github.com/Aguaments/network/fixtures/channel-artifacts/Org4MSPanchors.tx",
		},
	}
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "channel1",
		ChannelConfig:    "/home/alex/go/src/github.com/Aguaments/network/fixtures/channel-artifacts/channel1.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer0.orderer.com",
		ChaincodeID:      "transfer",
		ChaincodePath:    "/home/alex/go/src/github.com/Aguaments/network/chaincode/go/core",
		ChaincodeVersion: "v1",
	}

	var (
		a   []string
		ret string
	)

	// Declaration sdk
	sdk, err := sdkInit.Setup("/home/alex/go/src/github.com/Aguaments/network/config.yaml", &info)
	if err != nil {
		fmt.Println(">> sdk set error")
		fmt.Printf("%s", err)
		os.Exit(-1)
	}

	// Transfer: install chaincode
	packageID, err := sdkInit.InstallCC(&info)
	if err != nil {
		fmt.Println(">> install chaincode error: ", err)
		os.Exit(-1)
	}

	// Transfer: approve chaincode
	if err = sdkInit.ApproveLifecycle(&info, 1, packageID); err != nil {
		fmt.Println(">> approve chaincode error: ", err)
		os.Exit(-1)
	}

	// Transfer: init chaincode
	if err = sdkInit.InitCC(&info, false, sdk); err != nil {
		fmt.Println(">> init chaincode error: ", err)
		os.Exit(-1)
	}

	// Transfer: Acquire handler Org2
	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[1], sdk); err != nil {
		fmt.Println(">> InitService error: ", err)
		os.Exit(-1)
	}
	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">>Transfer Org2: set successfully.")
	//Transfer: Function test
	//Read json file
	content, err := ioutil.ReadFile("./product.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	fmt.Printf("%s", string(content))

	// CreateAsset
	a = []string{"CreateAsset", string(content)}
	ret, err = App.CreateAsset(a)
	if err != nil {
		fmt.Println(">> set error: ", err)
		os.Exit(-1)
	}
	fmt.Println(">> result: ", ret)

	// UserTransferAsset
	a = []string{"UserTransferAsset", "Org3MSP", "P00001"}
	ret, err = App.UserTransferAsset(a)
	if err != nil {
		fmt.Println(">> set error: ", err)
		os.Exit(-1)
	}
	fmt.Println(">> result: ", ret)

	// Transfer: Acquire handler Org0
	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {
		fmt.Println(">> InitService error: ", err)
		os.Exit(-1)
	}
	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">>Transfer Org0: set successfully.")

	// SuperTransferAsset
	a = []string{"SuperTransferAsset", "P00001"}
	ret, err = App.SuperTransferAsset(a)
	if err != nil {
		fmt.Println(">> set error: ", err)
		os.Exit(-1)
	}
	fmt.Println(">> result: ", ret)

	// Transfer: Acquire handler Org0
	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[2], sdk); err != nil {
		fmt.Println(">> InitService error: ", err)
		os.Exit(-1)
	}
	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">>Transfer Org3: set successfully.")

	//Read json file
	content, err = ioutil.ReadFile("./manufacture.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// SelfTransfer
	a = []string{"SelfTransfer", "P00001", "5", string(content)}
	ret, err = App.SelfTransfer(a)
	if err != nil {
		fmt.Println(">> set error: ", err)
		os.Exit(-1)
	}
	fmt.Println(">> result: ", ret)

}
