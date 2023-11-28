package main

import "fmt"

var (
	Version    string = "1.0.0"
	Commit     string
	CommitDate string
	TreeState  string
)

func VersionString() (out string) {
	out = fmt.Sprintf("rig v%s", Version)
	if Commit == "" {
		return
	} else {
		out += fmt.Sprintf("-%s", Commit)
	}

	if TreeState != "" {
		out += fmt.Sprintf(" %s", TreeState)
	}
	if CommitDate != "" {
		out += fmt.Sprintf(" (%s)", CommitDate)
	}

	return
}
