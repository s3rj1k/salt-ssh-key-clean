package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultRPCURL = "https://internalrpc.mirohost.net/v1/"
	envKeyRPCURL  = "RPC_URL"

	defaultRPCBasicAuthUser = "mirohost_test"
	envKeyRPCBasicAuthUser  = "BASIC_AUTH_USER"

	defaultRPCBasicAuthPass = "809_VfghjlfK"
	envKeyRPCBasicAuthPass  = "BASIC_AUTH_PASS"

	defaultRPCAccessKey = "rdbmmzyycxv5LTj5GAL8eMibAyry/RtWV+RajHA3pMk="
	envKeyRPCAccessKey  = "ACCESS_KEY"

	defaultRosterTargetTimeout = 300
	envKeyRosterTargetTimeout  = "SSH_TIMEOUT"
)

const ( // intentionally unconfigurable in runtime
	defaultRosterTargetUser          = "root"
	defaultRosterTargetThinDirPrefix = "/"
	defaultRosterTargetThinDirSuffix = "/salt/"

	defaultProjectNameForGetListMethods = "mirohost"

	defaultHostingNodeListSuffix      = "hosting"
	defaultHostingContainerListSuffix = ""
	defaultServiceDevicesListSuffix   = "service"

	defaultEVPSShortTypeName           = "vs"
	defaultSharedHostingShortTypeName  = "sd"
	defaultSmartDedicatedShortTypeName = "sm"
	defaultUndefinedShortTypeName      = "undef"

	defaultHostStatusSkipList = "reserved,unused,deleted" // string slice separated by comma
)

// CreateDefaultConfig creates default application config.
func CreateDefaultConfig() *Config {
	c := new(Config)

	c.RPCURL = defaultRPCURL
	c.RPCBasicAuthUser = defaultRPCBasicAuthUser
	c.RPCBasicAuthPass = defaultRPCBasicAuthPass
	c.RPCAccessKey = defaultRPCAccessKey

	c.ProjectNameForGetListMethods = defaultProjectNameForGetListMethods

	c.RosterTargetUser = defaultRosterTargetUser
	c.RosterTargetThinDirPrefix = defaultRosterTargetThinDirPrefix
	c.RosterTargetThinDirSuffix = defaultRosterTargetThinDirSuffix

	c.RosterTargetTimeout = defaultRosterTargetTimeout

	c.HostingNodeListSuffix = defaultHostingNodeListSuffix
	c.HostingContainerListSuffix = defaultHostingContainerListSuffix
	c.ServiceDevicesListSuffix = defaultServiceDevicesListSuffix

	c.EVPSShortTypeName = defaultEVPSShortTypeName
	c.SharedHostingShortTypeName = defaultSharedHostingShortTypeName
	c.SmartDedicatedShortTypeName = defaultSmartDedicatedShortTypeName
	c.UndefinedShortTypeName = defaultUndefinedShortTypeName

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

	ProjectNameForGetListMethods string

	RosterFilePath string

	RosterTargetUser          string
	RosterTargetThinDirPrefix string
	RosterTargetThinDirSuffix string
	RosterTargetTimeout       int

	HostingNodeListSuffix      string
	HostingContainerListSuffix string
	ServiceDevicesListSuffix   string

	EVPSShortTypeName           string
	SharedHostingShortTypeName  string
	SmartDedicatedShortTypeName string
	UndefinedShortTypeName      string

	HostStatusSkipList map[string]struct{}
}

// ReadFromEnvironment reads configuration parameters from environment variables.
func (c *Config) ReadFromEnvironment() error {
	defer func() {
		_ = os.Unsetenv(envKeyRPCURL)
		_ = os.Unsetenv(envKeyRPCBasicAuthUser)
		_ = os.Unsetenv(envKeyRPCBasicAuthPass)
		_ = os.Unsetenv(envKeyRPCAccessKey)
		_ = os.Unsetenv(envKeyRosterTargetTimeout)
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

	if val, ok := os.LookupEnv(envKeyRosterTargetTimeout); ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("config: invalid data (%s) for %s", val, envKeyRosterTargetTimeout)
		}

		if i < 0 {
			return fmt.Errorf("config: invalid data (%d) for %s", i, envKeyRosterTargetTimeout)
		}

		c.RosterTargetTimeout = i
	}

	return nil
}

// GetRosterTargetThinDir returns assembled roster 'thin_dir'.
func (c *Config) GetRosterTargetThinDir() string {
	return filepath.Join(
		c.RosterTargetThinDirPrefix,
		c.RosterTargetUser,
		c.RosterTargetThinDirSuffix,
	) + "/"
}
