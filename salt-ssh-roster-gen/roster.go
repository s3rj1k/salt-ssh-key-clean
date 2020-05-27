package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Roster target related constants
const (
	IgnoreBackupTargetInRosterTargetValue = "UNKNOWN"
)

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
		} `yaml:"grains"`
	} `yaml:"minion_opts"`
}

// CreateTarget returns roster target objected created from input data.
func CreateTarget(el GetListResultInnerObj, cfg *Config, roles ...string) Target {
	var t Target

	t.Host = el.СonfigurationManagement.FQDN
	t.User = cfg.RosterTargetUser
	t.Port = el.СonfigurationManagement.Port
	t.ThinDir = cfg.GetRosterTargetThinDir()
	t.Timeout = cfg.RosterTargetTimeout

	roles = append(roles, el.Type)
	t.MinionOpts.Grains.Roles = FilterStringSlice(roles)

	if el.CTID != nil {
		t.MinionOpts.Grains.Virtual.CTID = *el.CTID
		t.MinionOpts.Grains.Virtual.Parent = el.Node
	}

	if !strings.EqualFold(el.Backup, IgnoreBackupTargetInRosterTargetValue) {
		t.MinionOpts.Grains.Backup.Target = el.Backup
		t.MinionOpts.Grains.Backup.CreateBackups = el.CreateBackups
	}

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
