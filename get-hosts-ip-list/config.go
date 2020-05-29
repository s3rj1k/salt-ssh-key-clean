package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultRPCURL = "https://internalrpc.mirohost.net/v1/"
	envKeyRPCURL  = "RPC_URL"

	defaultRPCBasicAuthUser = "salt-ssh"
	envKeyRPCBasicAuthUser  = "BASIC_AUTH_USER"

	defaultRPCBasicAuthPass = "OeCoo0koog7iecaf"
	envKeyRPCBasicAuthPass  = "BASIC_AUTH_PASS"

	defaultRPCAccessKey = "ivC8n+Ll3G9VJNaK4s5KEQiH1acpRTYDU7834eUwifc2"
	envKeyRPCAccessKey  = "ACCESS_KEY"
)

const ( // intentionally unconfigurable in runtime
	defaultProjectNameForGetListMethods = "mirohost"

	defaultHostStatusSkipList = "reserved,unused,deleted" // string slice separated by comma
)

// CreateDefaultConfig creates default application config.
func CreateDefaultConfig() *Config {
	c := new(Config)

	c.RPCURL = defaultRPCURL
	c.RPCBasicAuthUser = defaultRPCBasicAuthUser
	c.RPCBasicAuthPass = defaultRPCBasicAuthPass
	c.RPCAccessKey = defaultRPCAccessKey

	// seed default status skip list
	c.HostStatusSkipList = make(map[string]struct{})
	for _, el := range strings.Split(defaultHostStatusSkipList, ",") {
		c.HostStatusSkipList[el] = struct{}{}
	}

	return c
}

// Config defines application configuration object.
type Config struct {
	RPCURL           string
	RPCBasicAuthUser string
	RPCBasicAuthPass string
	RPCAccessKey     string

	HostStatusSkipList map[string]struct{}
}

// ReadFromEnvironment reads configuration parameters from environment variables.
func (c *Config) ReadFromEnvironment() error {
	defer func() {
		_ = os.Unsetenv(envKeyRPCURL)
		_ = os.Unsetenv(envKeyRPCBasicAuthUser)
		_ = os.Unsetenv(envKeyRPCBasicAuthPass)
		_ = os.Unsetenv(envKeyRPCAccessKey)
	}()

	if val, ok := os.LookupEnv(envKeyRPCURL); ok {
		if !strings.HasPrefix(val, "https://") || !strings.HasPrefix(val, "http://") {
			return fmt.Errorf("config: invalid data (%s) for %s", val, envKeyRPCURL)
		}

		c.RPCURL = val
	}

	if val, ok := os.LookupEnv(envKeyRPCBasicAuthUser); ok {
		c.RPCBasicAuthUser = val
	}

	if val, ok := os.LookupEnv(envKeyRPCBasicAuthPass); ok {
		c.RPCBasicAuthPass = val
	}

	if val, ok := os.LookupEnv(envKeyRPCAccessKey); ok {
		c.RPCAccessKey = val
	}

	return nil
}
