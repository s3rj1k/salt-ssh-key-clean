package main

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
