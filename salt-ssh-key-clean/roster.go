package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// target describes single roster target.
type target struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Port int    `yaml:"port"`
}

/*
  cmd = [
    'ssh-keygen -R [{}]:{}'.format(fqdn, port),
    'ssh-keygen -R [{}]:{}'.format(ip, port),
    'ssh-keygen -R [{}]:{},[{}]:{}'.format(fqdn, port, ip, port),
    'ssh-keygen -R {}'.format(fqdn),
    'ssh-keygen -R {}'.format(ip),
    'ssh-keygen -R {},{}'.format(fqdn, ip)
  ]
*/

// RemoveList creates a list of combinations for `ssh-keygen -R` command.
func (t target) RemoveList() []string {
	t.Host = strings.TrimSpace(t.Host)

	if t.Port == defaultSSHPort {
		return []string{t.Host}
	}

	return []string{fmt.Sprintf("[%s]:%d", t.Host, t.Port)}
}

func (t target) String() string {
	var out string

	if len(t.User) > 0 {
		out = fmt.Sprintf("%s@%s", t.User, t.Host)
	} else {
		out = t.Host
	}

	if t.Port > 0 && t.Port != defaultSSHPort {
		out = fmt.Sprintf("[%s]:%d", out, t.Port)
	}

	return out
}

func parseRoster(path string) (map[string]target, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("roster: %w", err)
	}

	defer fd.Close()

	var roster map[string]target
	if err := yaml.NewDecoder(fd).Decode(&roster); err != nil {
		return nil, fmt.Errorf("roster: %w", err)
	}

	return roster, nil
}

func getTargetsFromRoster(path string) (<-chan target, error) {
	roster, err := parseRoster(path)
	if err != nil {
		return nil, err
	}

	targets := make(chan target, len(roster))
	for _, target := range roster {
		targets <- target
	}

	close(targets)

	return targets, nil
}
