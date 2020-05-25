package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/s3rj1k/jrpc2/client"
)

/*

curl -X POST https://internalrpc.mirohost.net/v1/ -u user:password -H "Content-Type:application/json" -d '{
    "jsonrpc":"2.0",
    "method": "getServiceDevicesList",
    "params": {
        "project":"mirohost",
        "accessKey":"rpcAccessKey"
    }
}' | jq

curl -X POST https://internalrpc.mirohost.net/v1/ -u user:password -H "Content-Type:application/json" -d '{
    "jsonrpc":"2.0",
    "method": "getNodesList",
    "params": {
        "accessKey":"rpcAccessKey"
    }
}' | jq

curl -X POST https://internalrpc.mirohost.net/v1/ -u user:password -H "Content-Type:application/json" -d '{
    "jsonrpc":"2.0",
    "method": "getContainersList",
    "params": {
        "accessKey":"rpcAccessKey"
    }
}' | jq

*/

func getListWrapper(c *client.Config, key, method, project string) ([]GetListResultObj, error) {
	// prepare results object
	resultObj := make([]GetListResultObj, 0)

	// JSON-RPC params field
	paramsObj := GetListParamsObj{
		AccessKey: key,
	}

	// set project name
	if len(project) > 0 {
		paramsObj.Project = project
	}

	// convert params object to bytes
	paramsData, err := json.Marshal(paramsObj)
	if err != nil {
		return resultObj, fmt.Errorf("method=%s params error: %s", method, err.Error())
	}

	// sent JSON-RPC request and get back response
	resultRawData, err := c.Call(method, paramsData)
	if err != nil {
		// error data object
		var errObj GetListErrInnerDataObj

		// decode inner (JSON-RPC) error message
		if errInnerDataObjError := json.Unmarshal(resultRawData, &errObj); errInnerDataObjError == nil {
			return resultObj, fmt.Errorf("%s: %s", err.Error(), strings.Join(errObj.Message, ", "))
		}

		// return plain error
		return resultObj, err
	}

	// decode result
	err = json.Unmarshal(resultRawData, &resultObj)
	if err != nil {
		return resultObj, fmt.Errorf("method=%s result error: %s", method, err.Error())
	}

	return resultObj, nil
}

// GetServiceDevicesList calls remote JSON-RPC server to get list of service hosts (physical and virtual).
func GetServiceDevicesList(c *client.Config, key, project string) ([]GetListResultObj, error) {
	return getListWrapper(c, key, "getServiceDevicesList", project)
}

// GetNodeList calls remote JSON-RPC server to get list of hosting nodes.
func GetNodeList(c *client.Config, key string) ([]GetListResultObj, error) {
	return getListWrapper(c, key, "getNodesList", "")
}

// GetContainersList calls remote JSON-RPC server to get list of hosting containers.
func GetContainersList(c *client.Config, key string) ([]GetListResultObj, error) {
	return getListWrapper(c, key, "getContainersList", "")
}
