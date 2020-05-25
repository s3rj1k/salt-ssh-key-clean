package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/s3rj1k/jrpc2/client"
	"gopkg.in/yaml.v2"
)

func main() {
	// create RPC client config
	rpc := client.GetConfig(defaultRPCURL)
	// set credentials
	rpc.SetBasicAuth(defaultRPCBasicAuthUser, defaultRPCBasicAuthPass)

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

	// seed default status skip list
	hostStatusSkipList := make(map[string]struct{})
	for _, el := range strings.Split(defaultHostStatusSkipList, ",") {
		hostStatusSkipList[el] = struct{}{}
	}

	// roster data
	roster := make(
		map[string]Target,
		len(nodeList)+len(containersList)+len(serviceHostsList),
	)

	for _, el := range nodeList {
		if el.Skip(hostStatusSkipList) {
			continue
		}

		id := fmt.Sprintf(
			"%s.%s",
			defaultNodeListSuffix,
			el.GetFQDNWithOutPublicSuffix(),
		)

		roster[id] = Target{
			Host: el.СonfigurationManagement.FQDN,
			User: defaultRosterTargetUser,
			Port: el.СonfigurationManagement.Port,
			ThinDir: filepath.Join(
				defaultRosterTargetThinDirPrefix,
				defaultRosterTargetUser,
				defaultRosterTargetThinDirSuffix,
			) + "/",
			Timeout: defaultRosterTargetTimeout,
		}
	}

	for _, el := range containersList {
		if el.Skip(hostStatusSkipList) {
			continue
		}

		id := fmt.Sprintf(
			"%s.%s.%s",
			el.GetShortFQDN(),
			el.GetShortNodeFQDN(),
			el.GetShortHostingContainerType(),
		)

		roster[id] = Target{
			Host: el.СonfigurationManagement.FQDN,
			User: defaultRosterTargetUser,
			Port: el.СonfigurationManagement.Port,
			ThinDir: filepath.Join(
				defaultRosterTargetThinDirPrefix,
				defaultRosterTargetUser,
				defaultRosterTargetThinDirSuffix,
			) + "/",
			Timeout: defaultRosterTargetTimeout,
		}
	}

	for _, el := range serviceHostsList {
		if el.Skip(hostStatusSkipList) {
			continue
		}

		id := fmt.Sprintf(
			"%s.%s",
			defaultServiceDevicesListSuffix,
			el.GetFQDNWithOutPublicSuffix(),
		)

		roster[id] = Target{
			Host: el.СonfigurationManagement.FQDN,
			User: defaultRosterTargetUser,
			Port: el.СonfigurationManagement.Port,
			ThinDir: filepath.Join(
				defaultRosterTargetThinDirPrefix,
				defaultRosterTargetUser,
				defaultRosterTargetThinDirSuffix,
			) + "/",
			Timeout: defaultRosterTargetTimeout,
		}
	}

	b, err := yaml.Marshal(roster)
	if err != nil {
		fatal.Fatal(err)
	}

	fmt.Println(string(b))
}
