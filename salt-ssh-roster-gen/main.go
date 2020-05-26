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
	hostingNodeList, err := GetNodeList(rpc, defaultRPCAccessKey)
	if err != nil {
		fatal.Fatal(err)
	}

	// get list of hosting containers
	hostingContainersList, err := GetContainersList(rpc, defaultRPCAccessKey)
	if err != nil {
		fatal.Fatal(err)
	}

	// get list of service hosts
	serviceHostsList, err := GetServiceDevicesList(rpc, defaultRPCAccessKey, defaultProjectNameForGetListMethods)
	if err != nil {
		fatal.Fatal(err)
	}

	// roster data
	roster := CreateNewRoster(
		len(hostingNodeList) + len(hostingContainersList) + len(serviceHostsList),
	)

	// add hosting nodes to roster
	for _, el := range hostingNodeList {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		id := el.GetHostingNodeID(cfg.HostingNodeListSuffix)

		roster.Data[id] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    cfg.RosterTargetUser,
			Port:    el.СonfigurationManagement.Port,
			ThinDir: cfg.GetRosterTargetThinDir(),
			Timeout: cfg.RosterTargetTimeout,
		}
	}

	// add hosting containers to roster
	for _, el := range hostingContainersList {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		id := el.GetHostingContainerID(cfg.HostingContainerListSuffix)

		roster.Data[id] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    cfg.RosterTargetUser,
			Port:    el.СonfigurationManagement.Port,
			ThinDir: cfg.GetRosterTargetThinDir(),
			Timeout: cfg.RosterTargetTimeout,
		}
	}

	// add service hosts to roster
	for _, el := range serviceHostsList {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		id := el.GetServiceDeviceID(cfg.ServiceDevicesListSuffix)

		roster.Data[id] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    cfg.RosterTargetUser,
			Port:    el.СonfigurationManagement.Port,
			ThinDir: cfg.GetRosterTargetThinDir(),
			Timeout: cfg.RosterTargetTimeout,
		}
	}

	b, err := yaml.Marshal(roster.Data)
	if err != nil {
		fatal.Fatal(err)
	}

	fmt.Println(string(b))
}
