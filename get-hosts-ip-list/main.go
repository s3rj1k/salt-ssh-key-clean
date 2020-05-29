package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/s3rj1k/jrpc2/client"
)

func main() {
	// create default application config
	cfg := CreateDefaultConfig()

	// custom help
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Printf("  environment variables:\n")
		fmt.Printf("    %s - sets RPC endpoint URL\n", envKeyRPCURL)
		fmt.Printf("    %s - sets RPC endpoint BasicAuth Username\n", envKeyRPCBasicAuthUser)
		fmt.Printf("    %s - sets RPC endpoint BasicAuth Password\n", envKeyRPCBasicAuthPass)
		fmt.Printf("    %s - sets RPC endpoint Access key\n", envKeyRPCAccessKey)
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

	// get list of service hosts
	serviceHostsList, err := GetServiceDevicesList(rpc, defaultRPCAccessKey, defaultProjectNameForGetListMethods)
	if err != nil {
		fatal.Fatal(err)
	}

	for _, el := range hostingNodeList.Data {
		if el.IP == nil {
			continue
		}

		addresses := make([]string, 0, len(el.IP))

		for _, ip := range el.IP {
			addresses = append(addresses, ip.IP.String())
		}

		addresses = FilterStringSlice(addresses)

		fmt.Printf(
			"Role: Hosting FQDN: %s IP: %s\n",
			el.FQDN,
			strings.Join(addresses, ", "),
		)
	}

	for _, el := range serviceHostsList.Data {
		if el.IP == nil {
			continue
		}

		addresses := make([]string, 0, len(el.IP))

		for _, ip := range el.IP {
			addresses = append(addresses, ip.IP.String())
		}

		addresses = FilterStringSlice(addresses)

		fmt.Printf(
			"Role: Service FQDN: %s IP: %s\n",
			el.FQDN,
			strings.Join(addresses, ", "),
		)
	}
}
