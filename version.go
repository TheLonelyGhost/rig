package main

import "fmt"

var (
	Version    string = "1.0.0"
	TreeState  string = "dirty"
	Commit     string
	CommitDate string
)

func VersionString() (out string) {
	out = fmt.Sprintf("rig %s", Version)
	if Commit != "" {
		out += fmt.Sprintf(" %s", Commit)
	}
	if TreeState != "" {
		out += fmt.Sprintf(" (%s)", TreeState)
	}

	return
}
