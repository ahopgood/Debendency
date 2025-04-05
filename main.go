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

	// Specifically handle the case where we are asked for the help prompt or an error returns the help prompt
	if flagErr == flag.ErrHelp {
		fmt.Printf("%s\n", flagOutput)
		os.Exit(2)
	} else if flagErr != nil {
		fmt.Printf("got error:%#v", flagErr)
		fmt.Printf(
			"output:%s\n", flagOutput)
		os.Exit(1)
	}

	fmt.Printf("%#v\n", conf)
	pkg.ConfigureLogger(*conf)

	slog.Debug("%#v\n", conf)

	cache := pkg.NewCache(*conf)
	// enter cache dir and clear it
	startingDir, err := cache.ClearBefore()
	if err != nil {
		os.Exit(1)
	}

	packageModelMap := make(map[string]*pkg.PackageModel)
	packageModelList := make([]*pkg.PackageModel, 0)
	// download packages and build map
	pkg.NewAnalyser(conf).BuildPackage(conf.PackageName, packageModelMap, &packageModelList)

	// move out of cache dir back to calling directory
	err = os.Chdir(startingDir)
	if err != nil {
		slog.Error("Main error:", fmt.Errorf("Failed to change back to calling dir: %s\n%#v\n", startingDir, err))
		os.Exit(1)
	}

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

	// Pass in an io.Writer, in the case of os.Stdout for standard out
	if true == conf.GenerateSalt {
		// Root package
		salt.ToSaltDefinition(os.Stdout, packageModelList[0])
		// Print off dependencies
		salt.ToSaltDefinitions(os.Stdout, packageModelList[1:])
	}
}
