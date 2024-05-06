package pkg

import (
	"bytes"
	"flag"
)

type Config struct {
	PackageName              string
	GenerateSalt             bool
	GenerateDiagram          bool
	InstallerLocation        string
	ExcludeInstalledPackages bool
}

func ParseFlags(programName string, args []string) (config *Config, output string, err error) {

	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	// This is so we can return any flag parsing output
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var conf Config
	flags.StringVar(&conf.PackageName, "p", "", ".deb package name to calculate dependencies for")
	flags.BoolVar(&conf.GenerateSalt, "s", false, "output dependencies as salt code")
	flags.BoolVar(&conf.GenerateDiagram, "d", false, "output dependencies as a diagram")
	flags.StringVar(&conf.InstallerLocation, "o", "~/.debendency/cache", "cache directory to save installer files to")
	flags.BoolVar(&conf.ExcludeInstalledPackages, "e", false, "exclude already installed packages from output")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	return &conf, buf.String(), nil
}
