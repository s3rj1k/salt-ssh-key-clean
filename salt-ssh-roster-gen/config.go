package main

const (
	defaultRPCURL = "https://internalrpc.mirohost.net/v1/"
	// envKeyRPCURL  = "RPC_URL"

	defaultRPCBasicAuthUser = "mirohost_test"
	// envKeyRPCBasicAuthUser  = "BASIC_AUTH_USER"

	defaultRPCBasicAuthPass = "809_VfghjlfK"
	// envKeyRPCBasicAuthPass  = "BASIC_AUTH_PASS"

	defaultRPCAccessKey = "rdbmmzyycxv5LTj5GAL8eMibAyry/RtWV+RajHA3pMk="
	// envKeyRPCAccessKey  = "ACCESS_KEY"

	defaultProjectNameForGetListMethods = "mirohost"
	// envKeyProjectNameForGetListMethods  = "PROJECT_NAME"

	defaultRosterTargetUser    = "root"
	defaultRosterTargetThinDir = "/root/salt/"
	defaultRosterTargetTimeout = 300

	defaultNodeListSuffix           = "hosting"
	defaultServiceDevicesListSuffix = "service"

	defaultEVPSShortTypeName           = "vs"
	defaultSharedHostingShortTypeName  = "sd"
	defaultSmartDedicatedShortTypeName = "sm"
	defaultUndefinedShortTypeName      = "undef"

	defaultHostStatusSkipList = "reserved,unused,deleted" // string slice separated by comma
)

/*
type config struct {
	RPCURL           string
	RPCBasicAuthUser string
	RPCBasicAuthPass string
	RPCAccessKey     string

	ProjectNameForGetListMethods string

	RosterTargetUser    string
	RosterTargetThinDir string
	RosterTargetTimeout string

	NodeListSuffix           string
	ServiceDevicesListSuffix string

	EVPSShortTypeName           string
	SharedHostingShortTypeName  string
	SmartDedicatedShortTypeName string
	UndefinedShortTypeName      string

	HostStatusSkipList map[string]struct{}
}

func (c *config) ReadFromEnvironment() {
	if val, ok := os.LookupEnv(envKeyRPCURL); ok {
		c.RPCURL = val
	} else {
		c.RPCURL = defaultRPCURL
	}

	if val, ok := os.LookupEnv(envKeyRPCBasicAuthUser); ok {
		c.RPCBasicAuthUser = val
	} else {
		c.RPCBasicAuthUser = defaultRPCBasicAuthUser
	}

	if val, ok := os.LookupEnv(envKeyRPCBasicAuthPass); ok {
		c.RPCBasicAuthPass = val
	} else {
		c.RPCBasicAuthPass = defaultRPCBasicAuthPass
	}

	if val, ok := os.LookupEnv(envKeyRPCAccessKey); ok {
		c.RPCAccessKey = val
	} else {
		c.RPCAccessKey = defaultRPCAccessKey
	}

	if val, ok := os.LookupEnv(envKeyProjectNameForGetListMethods); ok {
		c.ProjectNameForGetListMethods = val
	} else {
		c.ProjectNameForGetListMethods = defaultProjectNameForGetListMethods
	}
}
*/
