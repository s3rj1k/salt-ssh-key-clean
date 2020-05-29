package main

import (
	"net"
)

// GetListErrInnerDataObj defines JSON-RPC error inner data object, JSON-RPC server specific.
type GetListErrInnerDataObj struct {
	Message []string `json:"message"`
}

// GetListParamsObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList object.
type GetListParamsObj struct {
	AccessKey string `json:"accessKey"`
	Project   string `json:"project,omitempty"`
}

// GetListResultInnerObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList single element of result object.
type GetListResultInnerObj struct {
	FQDN   string `json:"fqdn"`
	Status string `json:"status"`

	Type string `json:"type,omitempty"`

	IP []struct {
		IP net.IP `json:"ip"`
	} `json:"ip,omitempty"`
}

// GetListResultObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList result object.
type GetListResultObj struct {
	Data   []GetListResultInnerObj
	Method string
}

// Skip is used for filtering out invalid roster targets.
func (s GetListResultInnerObj) Skip(skip map[string]struct{}) bool {
	if _, ok := skip[s.Status]; ok {
		return true
	}

	return false
}
