package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Roster target related constants
const (
	IgnoreBackupTargetInRosterTargetValue = "UNKNOWN"
)

// IPRecord describes IP VLAN relation.
type IPRecord struct {
	IP      net.IP  `yaml:"IP"`
	Gateway net.IP  `yaml:"Gateway,omitempty"`
	Network *string `yaml:"Network,omitempty"` // https://github.com/golang/go/issues/12803
	IsIPv4  bool    `yaml:"IsIPv4"`
	IsIPv6  bool    `yaml:"IsIPv6"`
	VlanID  int     `yaml:"VlanID,omitempty"`
}

// Target describes single target element of salt-ssh roster.
type Target struct {
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	Port       int    `yaml:"port"`
	ThinDir    string `yaml:"thin_dir"`
	Timeout    int    `yaml:"timeout"`
	MinionOpts struct {
		Grains struct {
			Roles []string `yaml:"roles"`

			Virtual struct {
				Parent string `yaml:"parent"`
				CTID   string `yaml:"ctid"`
			} `yaml:"virtual,omitempty"`

			Backup struct {
				Target        string `yaml:"target"`
				CreateBackups bool   `yaml:"createBackups"`
			} `yaml:"backup,omitempty"`

			Network struct {
				IP []IPRecord `yaml:"IP,omitempty"`

				IPv6Hextet string `yaml:"v6Hextet"`
			} `yaml:"network"`
		} `yaml:"grains"`
	} `yaml:"minion_opts"`
}

// CreateTarget returns roster target objected created from input data.
func CreateTarget(el GetListResultInnerObj, cfg *Config, roles ...string) Target {
	var t Target

	t.Host = strings.TrimSpace(el.СonfigurationManagement.FQDN)
	t.User = cfg.RosterTargetUser
	t.Port = el.СonfigurationManagement.Port
	t.ThinDir = cfg.GetRosterTargetThinDir()
	t.Timeout = cfg.RosterTargetTimeout

	roles = append(roles, strings.TrimSpace(el.Type))
	t.MinionOpts.Grains.Roles = FilterStringSlice(roles)

	if el.CTID != nil {
		t.MinionOpts.Grains.Virtual.CTID = strings.TrimSpace(*el.CTID)
		t.MinionOpts.Grains.Virtual.Parent = strings.TrimSpace(el.Node)
	}

	if !strings.EqualFold(strings.TrimSpace(el.Backup), IgnoreBackupTargetInRosterTargetValue) {
		t.MinionOpts.Grains.Backup.Target = strings.TrimSpace(el.Backup)
		t.MinionOpts.Grains.Backup.CreateBackups = el.CreateBackups
	}

	ip := make([]IPRecord, 0, len(el.IP))

	for _, el := range el.IP {
		if el.IP == nil {
			continue
		}

		network := new(string)
		if el.Network != nil {
			if _, ntwrk, err := net.ParseCIDR(*el.Network); err == nil {
				str := ntwrk.String()
				network = &str
			}
		}

		ip = append(ip, IPRecord{
			IP:      el.IP,
			IsIPv4:  IsIPv4(el.IP),
			IsIPv6:  IsIPv6(el.IP),
			VlanID:  el.VlanID,
			Gateway: el.Gateway,
			Network: network,
		})
	}

	t.MinionOpts.Grains.Network.IP = ip

	t.MinionOpts.Grains.Network.IPv6Hextet = strings.TrimSpace(el.IPV6Hextet)

	return t
}

// Roster describes full list of targets inside roster.
type Roster struct {
	Data map[string]Target
}

// CreateNewRoster returns new roster object initialized to defined capacity.
func CreateNewRoster(cap int) *Roster {
	roster := new(Roster)

	roster.Data = make(
		map[string]Target, cap,
	)

	return roster
}

// SaveToFile writes YAML encoded roster data to file.
func (r *Roster) SaveToFile(path string) error {
	path = filepath.Clean(path)

	// create file or open as RW
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("roster error: %w", err)
	}

	// close fd on exit
	defer fd.Close()

	// create new encoder
	enc := yaml.NewEncoder(fd)

	// flush data
	defer enc.Close()

	// encode to fd
	if err := enc.Encode(r.Data); err != nil {
		return fmt.Errorf("roster error: %w", err)
	}

	return nil
}
