package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/s3rj1k/jrpc2/client"
	"gopkg.in/yaml.v2"
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

const (
	defaultRPCAccessKey                 = "rdbmmzyycxv5LTj5GAL8eMibAyry/RtWV+RajHA3pMk="
	defaultRPCLink                      = "https://internalrpc.mirohost.net/v1/"
	defaultRPCBasicAuthPass             = "809_VfghjlfK"
	defaultRPCBasicAuthUser             = "mirohost_test"
	defaultProjectNameForGetListMethods = "mirohost"
)

// GetListErrInnerDataObj defines JSON-RPC error inner data object, JSON-RPC server specific.
type GetListErrInnerDataObj struct {
	Message []string `json:"message"`
}

// GetListResultObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList result object.
type GetListResultObj []struct {
	FQDN   string `json:"fqdn"`
	Node   string `json:"node,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status"`

	СonfigurationManagement struct {
		Enabled bool   `json:"enabled"`
		FQDN    string `json:"fqdn"`
		Port    int    `json:"port"`
	} `json:"configurationManagement"`
}

// GetListParamsObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList object.
type GetListParamsObj struct {
	AccessKey string `json:"accessKey"`
	Project   string `json:"project,omitempty"`
}

func getListWrapper(c *client.Config, key, method, project string) (GetListResultObj, error) {
	// prepare results object
	var resultObj GetListResultObj

	// JSON-RPC params field
	paramsObj := GetListParamsObj{
		AccessKey: key,
	}

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
func GetServiceDevicesList(c *client.Config, key, project string) (GetListResultObj, error) {
	return getListWrapper(c, key, "getServiceDevicesList", project)
}

// GetNodeList calls remote JSON-RPC server to get list of hosting nodes.
func GetNodeList(c *client.Config, key string) (GetListResultObj, error) {
	return getListWrapper(c, key, "getNodesList", "")
}

// GetContainersList calls remote JSON-RPC server to get list of hosting containers.
func GetContainersList(c *client.Config, key string) (GetListResultObj, error) {
	return getListWrapper(c, key, "getContainersList", "")
}

func main() {
	// create RPC client config
	rpc := client.GetConfig(defaultRPCLink)
	// set credentials
	rpc.SetBasicAuth(defaultRPCBasicAuthUser, defaultRPCBasicAuthPass)

	// get list of hosting nodes
	nodeList, err := GetNodeList(rpc, defaultRPCAccessKey)
	if err != nil {
		log.Fatal(err)
	}

	// get list of hosting containers
	containersList, err := GetContainersList(rpc, defaultRPCAccessKey)
	if err != nil {
		log.Fatal(err)
	}

	// get list of service hosts
	serviceHostsList, err := GetServiceDevicesList(rpc, defaultRPCAccessKey, defaultProjectNameForGetListMethods)
	if err != nil {
		log.Fatal(err)
	}

	roster := make(map[string]Target, len(nodeList)+len(containersList)+len(serviceHostsList))

	for _, el := range nodeList {
		if strings.EqualFold(el.Status, "reserved") {
			continue
		}

		if strings.EqualFold(el.Status, "unused") {
			continue
		}

		if strings.EqualFold(el.Status, "deleted") {
			continue
		}

		if el.СonfigurationManagement.Enabled {
			roster[el.FQDN] = Target{
				Host:    el.СonfigurationManagement.FQDN,
				User:    "root",
				Port:    el.СonfigurationManagement.Port,
				ThinDir: "/root/salt/",
				Timeout: 300,
			}
		}
	}

	for _, el := range containersList {
		if strings.EqualFold(el.Status, "reserved") {
			continue
		}

		if strings.EqualFold(el.Status, "unused") {
			continue
		}

		if strings.EqualFold(el.Status, "deleted") {
			continue
		}

		if el.СonfigurationManagement.Enabled {
			roster[el.FQDN] = Target{
				Host:    el.СonfigurationManagement.FQDN,
				User:    "root",
				Port:    el.СonfigurationManagement.Port,
				ThinDir: "/root/salt/",
				Timeout: 300,
			}
		}
	}

	for _, el := range serviceHostsList {
		if strings.EqualFold(el.Status, "reserved") {
			continue
		}

		if strings.EqualFold(el.Status, "unused") {
			continue
		}

		if strings.EqualFold(el.Status, "deleted") {
			continue
		}

		if el.СonfigurationManagement.Enabled {
			roster[el.FQDN] = Target{
				Host:    el.СonfigurationManagement.FQDN,
				User:    "root",
				Port:    el.СonfigurationManagement.Port,
				ThinDir: "/root/salt/",
				Timeout: 300,
			}
		}
	}

	b, err := yaml.Marshal(roster)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
