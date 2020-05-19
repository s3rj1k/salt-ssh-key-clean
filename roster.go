package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type target struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Port int    `yaml:"port"`
}

func (t target) String() string {
	var out string

	if len(t.User) > 0 {
		out = fmt.Sprintf("%s@%s", t.User, t.Host)
	} else {
		out = t.Host
	}

	if t.Port != 0 {
		out = fmt.Sprintf("[%s]:%d", out, t.Port)
	}

	return out
}

func parseRoster(path string) (map[string]target, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("roster error: %w", err)
	}

	defer fd.Close()

	var roster map[string]target
	if err := yaml.NewDecoder(fd).Decode(&roster); err != nil {
		return nil, fmt.Errorf("roster error: %w", err)
	}

	return roster, nil
}
