package commands

import (
	"fmt"
	"strings"
)

// Interfaces may allow for faking when testing our native commands
type Dpkg interface {
	IdentifyDependencies(filename string) []string
	ParseDependencies(output string) []string
}

type DpkgQuery interface {
	IsInstalled(packageName string) bool
}

type Dpkger struct {
	Cmd Command
}

func (dpkg Dpkger) IdentifyDependencies(filename string) []string {
	output, exitCode, err := dpkg.Cmd.Command("dpkg", "-I", filename)

	if err != nil {
		//shit the bed
		//log stuff
		//propagate error
	}

	if exitCode != 0 {
		// shit the bed
		// log output
		// propagate an error
	}

	//parse output into an array of package names
	return dpkg.ParseDependencies(output)

}

func (dpkg Dpkger) ParseDependencies(output string) []string {

	//Find the line Pre-Depends:
	_, afterPreDepends, found := strings.Cut(output, "Pre-Depends:")
	if found {
		output = afterPreDepends
	}
	//Find the line Depends:
	_, after, found := strings.Cut(output, "Depends:")
	if found {
		depends := strings.Split(after, "\n")[0]
		dependencies := strings.Split(depends, ",")
		for i := range dependencies {
			dependencies[i] = strings.TrimSpace(dependencies[i])

			// Here we handle/ignore the optional version brackets e.g. libc6 (>= 2.4)
			if strings.Contains(dependencies[i], " ") {
				dependencies[i] = strings.Split(dependencies[i], " ")[0]
			}
			// Here we handle/ignore the qualifier on a dependency e.g. python:any
			if strings.Contains(dependencies[i], ":") {
				dependencies[i] = strings.Split(dependencies[i], ":")[0]
			}

		}
		return dependencies
	}
	return make([]string, 0)
}

type Query struct {
	Cmd Command
}

func (query Query) IsInstalled(packageName string) bool {
	output, exitCode, err := query.Cmd.Command("dpkg-query", "-W", packageName)
	if err != nil {
		fmt.Printf("Encountered an unknown error on dpkg-query: %#v\n", err)
		return false
	}
	switch exitCode {
	case 0:
		return true
	case 1:
		return false
	default:
		fmt.Printf("Encountered an unknown exit code on dpkg-query: %d, with output %s\n", exitCode, output)
		return false
	}
}
