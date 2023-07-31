package sdkInit

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

//Write

func (t *Application) IPFShashWrite(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

//Query

func (t *Application) AcquireIPFSHash(args []string) (string, error) {
	response, err := t.SdkEnvInfo.Client.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{}}, channel.WithTargetEndpoints("peer0.org1.com"))
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}
	return string(response.Payload), nil
}

//Change state

func (t *Application) ChangeHashState(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *Application) SupervisionQuery(args []string) (string, error) {
	response, err := t.SdkEnvInfo.Client.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{}}, channel.WithTargetEndpoints("peer0.org1.com"))
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}
	return string(response.Payload), nil
}
