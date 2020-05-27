package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/s3rj1k/jrpc2/client"
)

func main() {
	if err := CheckIfRunUnderRoot(); err != nil {
		fatal.Fatal(err)
	}

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
		fmt.Printf("    DEBUG = TRUE - sets debug logging\n")
	}

	flag.Parse()

	// create RPC client config
	rpc := client.GetConfig(defaultRPCURL)
	// set credentials
	rpc.SetBasicAuth(defaultRPCBasicAuthUser, defaultRPCBasicAuthPass)

	// read configuration data from environment
	if err := cfg.ReadFromEnvironment(); err != nil {
		fatal.Fatal(err)
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
		len(hostingNodeList.Data) + len(hostingContainersList.Data) + len(serviceHostsList.Data),
	)

	// add hosting nodes to roster
	for _, el := range hostingNodeList.Data {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		id, err := el.GetID(cfg, hostingNodeList.Method)
		if err != nil {
			debug.Println(err)

			continue
		}

		roles, err := el.GetRoles(cfg, hostingNodeList.Method)
		if err != nil {
			debug.Println(err)

			continue
		}

		roster.Data[id] = CreateTarget(el, cfg, roles...)
	}

	// add hosting containers to roster
	for _, el := range hostingContainersList.Data {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		id, err := el.GetID(cfg, hostingContainersList.Method)
		if err != nil {
			debug.Println(err)

			continue
		}

		roles, err := el.GetRoles(cfg, hostingContainersList.Method)
		if err != nil {
			debug.Println(err)

			continue
		}

		roster.Data[id] = CreateTarget(el, cfg, roles...)
	}

	// add service hosts to roster
	for _, el := range serviceHostsList.Data {
		if el.Skip(cfg.HostStatusSkipList) {
			continue
		}

		id, err := el.GetID(cfg, serviceHostsList.Method)
		if err != nil {
			debug.Println(err)

			continue
		}

		roles, err := el.GetRoles(cfg, serviceHostsList.Method)
		if err != nil {
			debug.Println(err)

			continue
		}

		roster.Data[id] = CreateTarget(el, cfg, roles...)
	}

	if err := roster.SaveToFile(cfg.RosterFilePath); err != nil {
		fatal.Fatal(err)
	}
}
