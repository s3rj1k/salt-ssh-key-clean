package main

import (
	"fmt"
	"log"

	"github.com/s3rj1k/jrpc2/client"
	"gopkg.in/yaml.v2"
)

func main() {
	// create RPC client config
	rpc := client.GetConfig(defaultRPCURL)
	// set credentials
	rpc.SetBasicAuth(defaultRPCBasicAuthUser, defaultRPCBasicAuthPass)

	// create default application config
	cfg := CreateDefaultConfig()

	// read configuration data from environment
	if err := cfg.ReadFromEnvironment(); err != nil {
		log.Fatal(err)
	}

	// get list of hosting nodes
	nodeList, err := GetNodeList(rpc, defaultRPCAccessKey)
	if err != nil {
		fatal.Fatal(err)
	}

	// get list of hosting containers
	containersList, err := GetContainersList(rpc, defaultRPCAccessKey)
	if err != nil {
		fatal.Fatal(err)
	}

	// get list of service hosts
	serviceHostsList, err := GetServiceDevicesList(rpc, defaultRPCAccessKey, defaultProjectNameForGetListMethods)
	if err != nil {
		fatal.Fatal(err)
	}

	// roster data
	roster := make(
		map[string]Target,
		len(nodeList)+len(containersList)+len(serviceHostsList),
	)

	for _, el := range nodeList {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		roster[el.GetHostingNodeID(cfg.HostingNodeListSuffix)] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    cfg.RosterTargetUser,
			Port:    el.СonfigurationManagement.Port,
			ThinDir: cfg.GetRosterTargetThinDir(),
			Timeout: cfg.RosterTargetTimeout,
		}
	}

	for _, el := range containersList {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		roster[el.GetHostingContainerID(cfg.HostingContainerListSuffix)] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    cfg.RosterTargetUser,
			Port:    el.СonfigurationManagement.Port,
			ThinDir: cfg.GetRosterTargetThinDir(),
			Timeout: cfg.RosterTargetTimeout,
		}
	}

	for _, el := range serviceHostsList {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		roster[el.GetServiceDeviceID(cfg.ServiceDevicesListSuffix)] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    cfg.RosterTargetUser,
			Port:    el.СonfigurationManagement.Port,
			ThinDir: cfg.GetRosterTargetThinDir(),
			Timeout: cfg.RosterTargetTimeout,
		}
	}

	b, err := yaml.Marshal(roster)
	if err != nil {
		fatal.Fatal(err)
	}

	fmt.Println(string(b))
}
