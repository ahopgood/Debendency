package main

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/puml"
	"com/alexander/debendency/pkg/salt"
	"flag"
	"io/fs"
	"os"

	"fmt"
)

func main() {

	conf, flagOutput, flagErr := pkg.ParseFlags(os.Args[0], os.Args[1:])

	// Specifically handle the case where we are asked for the help prompt or an error returns the help prompt
	if flagErr == flag.ErrHelp {
		fmt.Println(flagOutput)
		os.Exit(2)
	} else if flagErr != nil {
		fmt.Println("got error:", flagErr)
		fmt.Println("output:\n", flagOutput)
		os.Exit(1)
	}

	fmt.Printf("%#v\n", conf)

	cache := pkg.Cache{}
	cache.ClearBefore()

	packageModelMap := make(map[string]*pkg.PackageModel)
	packageModelList := make([]*pkg.PackageModel, 0)
	firstPackage := pkg.NewAnalyser(conf).BuildPackage(conf.PackageName, packageModelMap, packageModelList)

	if true == conf.GenerateDiagram {
		// Need to create the file output here

		pumlDiagramString := puml.GenerateDiagram(conf, packageModelMap, packageModelList).Contents()
		fmt.Println(pumlDiagramString)
		err := os.WriteFile(packageModelList[0].Name, []byte(pumlDiagramString), fs.ModePerm)
		if err != nil {
			fmt.Errorf("Issue writing puml diagram to file: %\n", packageModelList[0].Name, err)
		}
	}

	if true == conf.GenerateSalt {
		salt.ToSaltDefinition(firstPackage)
	}
}
