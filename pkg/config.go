package pkg

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"strings"
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

	// If package flag not provided
	if len(strings.TrimSpace(flags.Lookup("p").Value.String())) == 0 {
		//return nil, buf.String(), flag.ErrHelp
		fmt.Fprintf(flags.Output(), "Usage of %s:\n", programName)
		flags.PrintDefaults()
		return nil, buf.String(), errors.New("flag: need to specify a package name via -p flag")
	}

	return &conf, buf.String(), nil
}
