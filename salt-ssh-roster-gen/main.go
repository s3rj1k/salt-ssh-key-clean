package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/s3rj1k/jrpc2/client"
)

func main() {
	// create default application config
	cfg := CreateDefaultConfig()

	// path to generated roster file
	flag.StringVar(&cfg.RosterFilePath, "roster", "/tmp/roster", "defines an location for the default roster file")

	// custom help
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Printf("  environment variables:\n")
		fmt.Printf("    %s - sets RPC endpoint URL\n", envKeyRPCURL)
		fmt.Printf("    %s - sets RPC endpoint BasicAuth Username\n", envKeyRPCBasicAuthUser)
		fmt.Printf("    %s - sets RPC endpoint BasicAuth Password\n", envKeyRPCBasicAuthPass)
		fmt.Printf("    %s - sets RPC endpoint Access key\n", envKeyRPCAccessKey)
		fmt.Printf("    %s - sets roster target ssh timeout\n", envKeyRosterTargetTimeout)
	}

	flag.Parse()

	// create RPC client config
	rpc := client.GetConfig(defaultRPCURL)
	// set credentials
	rpc.SetBasicAuth(defaultRPCBasicAuthUser, defaultRPCBasicAuthPass)

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

	if err := roster.SaveToFile(cfg.RosterFilePath); err != nil {
		fatal.Fatal(err)
	}
}
