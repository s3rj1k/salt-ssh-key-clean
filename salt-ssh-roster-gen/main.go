package main

import (
	"fmt"
	"log"

	"github.com/s3rj1k/jrpc2/client"
	"gopkg.in/yaml.v2"
)

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
		if el.Skip() {
			continue
		}

		id := fmt.Sprintf(
			"%s.hosting",
			el.GetFQDNWithOutPublicSuffix(),
		)

		roster[id] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    "root",
			Port:    el.СonfigurationManagement.Port,
			ThinDir: "/root/salt/",
			Timeout: 300,
		}

	}

	for _, el := range containersList {
		if el.Skip() {
			continue
		}

		id := fmt.Sprintf(
			"%s.%s.%s",
			el.GetShortFQDN(),
			el.GetShortNodeFQDN(),
			el.GetShortHostingContainerType(),
		)

		roster[id] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    "root",
			Port:    el.СonfigurationManagement.Port,
			ThinDir: "/root/salt/",
			Timeout: 300,
		}
	}

	for _, el := range serviceHostsList {
		if el.Skip() {
			continue
		}

		id := fmt.Sprintf(
			"%s.service",
			el.GetFQDNWithOutPublicSuffix(),
		)

		roster[id] = Target{
			Host:    el.СonfigurationManagement.FQDN,
			User:    "root",
			Port:    el.СonfigurationManagement.Port,
			ThinDir: "/root/salt/",
			Timeout: 300,
		}
	}

	b, err := yaml.Marshal(roster)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
