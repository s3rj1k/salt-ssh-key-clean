package main

// Target describes single target element of salt-ssh roster.
type Target struct {
	Host    string `yaml:"host"`
	User    string `yaml:"user"`
	Port    int    `yaml:"port"`
	ThinDir string `yaml:"thin_dir"`
	Timeout int    `yaml:"timeout"`
}
