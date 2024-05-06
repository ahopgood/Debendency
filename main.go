package main

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/puml"
	"com/alexander/debendency/pkg/salt"
	"flag"
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
	firstPackage := pkg.NewAnalyser(conf).BuildPackage(conf.PackageName, packageModelMap)

	if true == conf.GenerateDiagram {
		fmt.Println(puml.GenerateDiagram(conf, packageModelMap).Contents())
	}

	if true == conf.GenerateSalt {
		salt.ToSaltDefinition(firstPackage)
	}
}
