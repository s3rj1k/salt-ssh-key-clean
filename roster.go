package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type target struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Port int    `yaml:"port"`
}

func parseRoster(path string) (map[string]target, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	var roster map[string]target
	if err := yaml.NewDecoder(fd).Decode(&roster); err != nil {
		return nil, err
	}

	return roster, nil
}
