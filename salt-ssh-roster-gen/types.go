package main

import (
	"fmt"
	"net"
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

// GetListResultInnerObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList single element of result object.
type GetListResultInnerObj struct {
	CTID *string `json:"ctid,omitempty"`
	Node string  `json:"node,omitempty"`
	Type string  `json:"type,omitempty"`

	FQDN   string `json:"fqdn"`
	Status string `json:"status"`

	Backup        string `json:"backup,omitempty"`
	CreateBackups bool   `json:"createBackups,omitempty"`

	СonfigurationManagement struct {
		Enabled bool   `json:"enabled"`
		FQDN    string `json:"fqdn"`
		Port    int    `json:"port"`
	} `json:"configurationManagement"`

	IP []struct {
		VlanID int    `json:"vlanID,omitempty"`
		IP     net.IP `json:"ip"`
	} `json:"ip,omitempty"`

	// THIS IS A HACK !!!
	IPV6Hextet string `json:"v6Hextet"`
}

// GetListResultObj defines JSON-RPC GetServiceDevicesList/GetNodesList/GetContainersList result object.
type GetListResultObj struct {
	Data   []GetListResultInnerObj
	Method string
}

// GetCombinedRoles returns target roles.
func (s GetListResultInnerObj) GetCombinedRoles(roles ...string) []string {
	if s.CTID == nil {
		roles = append(roles, "physical")
	} else {
		roles = append(roles, "virtual")
	}

	roles = append(roles, s.Type)

	return FilterStringSlice(roles)
}

// Skip is used for filtering out invalid roster targets.
func (s GetListResultInnerObj) Skip(skip map[string]struct{}) bool {
	if _, ok := skip[s.Status]; ok {
		return true
	}

	if !s.СonfigurationManagement.Enabled {
		return true
	}

	return false
}

// GetShortFQDN returns short FQDN of a target.
func (s GetListResultInnerObj) GetShortFQDN() string {
	return GetShortFQDN(s.FQDN)
}

// GetFQDNWithoutPublicSuffix returns FQDN with public suffix stripped.
func (s GetListResultInnerObj) GetFQDNWithoutPublicSuffix() string {
	return GetFQDNWithOutPublicSuffix(s.FQDN)
}

// GetShortNodeFQDN returns short node FQDN of a target.
func (s GetListResultInnerObj) GetShortNodeFQDN() string {
	return GetShortFQDN(s.Node)
}

// GetShortHostingContainerType returns short type of hosting container
func (s GetListResultInnerObj) GetShortHostingContainerType() string {
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
func (s GetListResultInnerObj) GetHostingNodeID(suff string) string {
	return strings.TrimSuffix(
		fmt.Sprintf(
			"%s.%s",
			s.GetFQDNWithoutPublicSuffix(),
			suff,
		), ".",
	)
}

// GetHostingContainerID returns roster target ID for hosting container.
func (s GetListResultInnerObj) GetHostingContainerID(suff string) string {
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
func (s GetListResultInnerObj) GetServiceDeviceID(suff string) string {
	return strings.TrimSuffix(
		fmt.Sprintf(
			"%s.%s",
			s.GetFQDNWithoutPublicSuffix(),
			suff,
		), ".",
	)
}

// GetID is a wrapper function to get target ID based on called method.
func (s GetListResultInnerObj) GetID(cfg *Config, method string) (string, error) {
	switch method {
	case GetServiceDevicesListMethodName:
		return s.GetServiceDeviceID(cfg.ServiceDevicesListSuffix), nil
	case GetNodeListMethodName:
		return s.GetHostingNodeID(cfg.HostingNodeListSuffix), nil
	case GetContainersListMethodName:
		return s.GetHostingContainerID(cfg.HostingContainerListSuffix), nil
	}

	return "", fmt.Errorf("unknown RPC method: %s", method)
}

// GetRoles is a wrapper function to get list of roles for a target based on called method.
func (s GetListResultInnerObj) GetRoles(cfg *Config, method string) ([]string, error) {
	switch method {
	case GetServiceDevicesListMethodName:
		return s.GetCombinedRoles(cfg.ServiceDevicesListSuffix), nil
	case GetNodeListMethodName:
		return s.GetCombinedRoles(cfg.HostingNodeListSuffix), nil
	case GetContainersListMethodName:
		return s.GetCombinedRoles(cfg.HostingContainerListSuffix, cfg.HostingNodeListSuffix), nil
	}

	return nil, fmt.Errorf("unknown RPC method: %s", method)
}
