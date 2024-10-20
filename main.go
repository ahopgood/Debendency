package main

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/puml"
	"com/alexander/debendency/pkg/salt"
	"flag"
	"io/fs"
	"log/slog"
	"os"

	"fmt"
)

func main() {

	conf, flagOutput, flagErr := pkg.ParseFlags(os.Args[0], os.Args[1:])

	pkg.ConfigureLogger(*conf)

	// Specifically handle the case where we are asked for the help prompt or an error returns the help prompt
	if flagErr == flag.ErrHelp {
		slog.Error(flagOutput)
		os.Exit(2)
	} else if flagErr != nil {
		slog.Error("got error:", flagErr)
		slog.Error(
			"output:\n", flagOutput)
		os.Exit(1)
	}

	slog.Debug("%#v\n", conf)

	cache := pkg.Cache{}
	cache.ClearBefore()

	packageModelMap := make(map[string]*pkg.PackageModel)
	packageModelList := make([]*pkg.PackageModel, 0)
	firstPackage := pkg.NewAnalyser(conf).BuildPackage(conf.PackageName, packageModelMap, &packageModelList)

	slog.Info(fmt.Sprintf("Package list %#v", packageModelList))
	slog.Info(fmt.Sprintf("Package map %#v", packageModelMap))
	if true == conf.GenerateDiagram {
		// Need to create the file output here

		pumlDiagramString := puml.GenerateDiagram(conf, packageModelMap, packageModelList).Contents()
		filename := fmt.Sprintf("%s.puml", conf.PackageName)
		slog.Debug(pumlDiagramString)
		err := os.WriteFile(filename, []byte(pumlDiagramString), fs.ModePerm)
		if err != nil {
			slog.Error("Main error:", fmt.Errorf("Issue writing puml diagram to file: %s\n%#v\n", filename, err))
		}
	}

	if true == conf.GenerateSalt {
		salt.ToSaltDefinition(firstPackage)
	}
}
