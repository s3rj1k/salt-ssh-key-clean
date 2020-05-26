package main

import (
	"fmt"
	"strings"
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

// GetListResultObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList single element of result object.
type GetListResultObj struct {
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

// Skip is used for filtering out invalid roster targets.
func (s GetListResultObj) Skip(skip map[string]struct{}) bool {
	if _, ok := skip[s.Status]; ok {
		return true
	}

	if !s.СonfigurationManagement.Enabled {
		return true
	}

	return false
}

// GetShortFQDN returns short FQDN of a target.
func (s GetListResultObj) GetShortFQDN() string {
	return GetShortFQDN(s.FQDN)
}

// GetFQDNWithoutPublicSuffix returns FQDN with public suffix stripped.
func (s GetListResultObj) GetFQDNWithoutPublicSuffix() string {
	return GetFQDNWithOutPublicSuffix(s.FQDN)
}

// GetShortNodeFQDN returns short node FQDN of a target.
func (s GetListResultObj) GetShortNodeFQDN() string {
	return GetShortFQDN(s.Node)
}

// GetShortHostingContainerType returns short type of hosting container
func (s GetListResultObj) GetShortHostingContainerType() string {
	switch strings.ToLower(strings.TrimSpace(s.Type)) {
	case "vps":
		return defaultEVPSShortTypeName
	case "shared":
		return defaultSharedHostingShortTypeName
	case "smart":
		return defaultSmartDedicatedShortTypeName
	default:
		return defaultUndefinedShortTypeName
	}
}

// GetHostingNodeID returns roster target ID for hosting node.
func (s GetListResultObj) GetHostingNodeID(suff string) string {
	return strings.TrimSuffix(
		fmt.Sprintf(
			"%s.%s",
			s.GetFQDNWithoutPublicSuffix(),
			suff,
		), ".",
	)
}

// GetHostingContainerID returns roster target ID for hosting container.
func (s GetListResultObj) GetHostingContainerID(suff string) string {
	return strings.TrimSuffix(
		fmt.Sprintf(
			"%s.%s.%s.%s",
			s.GetShortFQDN(),
			s.GetShortNodeFQDN(),
			s.GetShortHostingContainerType(),
			suff,
		), ".",
	)
}

// GetServiceDeviceID returns roster target ID for service host.
func (s GetListResultObj) GetServiceDeviceID(suff string) string {
	return strings.TrimSuffix(
		fmt.Sprintf(
			"%s.%s",
			s.GetFQDNWithoutPublicSuffix(),
			suff,
		), ".",
	)
}
