package sdkInit

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *Application) CreateAsset(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *Application) UserTransferAsset(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *Application) SuperTransferAsset(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *Application) SelfTransfer(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *Application) UserQuery(args []string) (string, error) {
	response, err := t.SdkEnvInfo.Client.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}
	return string(response.Payload), nil
}
