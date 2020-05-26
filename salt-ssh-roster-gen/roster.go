package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Target describes single target element of salt-ssh roster.
type Target struct {
	Host    string `yaml:"host"`
	User    string `yaml:"user"`
	Port    int    `yaml:"port"`
	ThinDir string `yaml:"thin_dir"`
	Timeout int    `yaml:"timeout"`
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
