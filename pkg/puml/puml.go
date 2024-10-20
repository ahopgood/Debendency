package puml

import (
	"bytes"
	"com/alexander/debendency/pkg"
	"fmt"
	"log/slog"
	"text/template"
)

var saltTemplate = `@startuml
digraph test {
{{range .}}
	"{{.From}}\n({{.FromVersion}})"{{if .To}} -> "{{.To}}\n({{.ToVersion}})"{{end}}
{{- end}}

}

@enduml`

var pumlTemplate = `@startuml
digraph test {
{{- range .Packages}}
    "{{.Name}}" [label="{{.Name}}\n({{.Version}})"];
{{- end}}
{{range .Dependencies}}
{{if .To}}    "{{.From}}" -> "{{.To}}"{{end}}
{{- end}}
}

@enduml`

func GenerateDiagram(config *pkg.Config, modelMap map[string]*pkg.PackageModel, modelList []*pkg.PackageModel) Uml {
	dependencies := make([]Dependency, 0)
	packages := make([]Package, 0)

	for _, fromModel := range modelList {
		//slog.Info(fmt.Sprintf("%s %#v\n", fromModel.Name, fromModel))

		if config.ExcludeInstalledPackages && fromModel.IsInstalled {

		} else {
			packages = append(packages, Package{
				Name:    fromModel.Name,
				Version: fromModel.Version,
			})
		}

		if fromModel.Dependencies == nil {
			dependencies = append(dependencies, Dependency{
				From:        fromModel.Name,
				FromVersion: fromModel.Version,
			})
		}
		for _, toModel := range fromModel.OrderedDependencies {
			//if config.ExcludeInstalledPackages && (fromModel.IsInstalled && toModel.IsInstalled) {
			if config.ExcludeInstalledPackages && toModel.IsInstalled {
				//fmt.Printf()
				//Log output separately from this model
			} else {
				slog.Info(fmt.Sprintf("From %s to %s\n", fromModel.Name, toModel.Name))
				dependencies = append(dependencies, Dependency{
					From:        fromModel.Name,
					FromVersion: fromModel.Version,
					To:          toModel.Name,
					ToVersion:   toModel.Version,
				})
			}
		}
	}
	//slices.Sort(dependencies)
	puml := NewUml()

	pumlModel := PumlModel{
		Packages:     packages,
		Dependencies: dependencies,
	}
	tmpl, err := template.New("puml").Parse(pumlTemplate)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, pumlModel)
	if err != nil {
		panic(err)
	}
	puml.Buffer = b
	return puml
}

type UmlDiagram interface {
	Contents() string
}

type Uml struct {
	Buffer bytes.Buffer
}

func NewUml() Uml {
	return Uml{}
}

func (uml Uml) Contents() string {
	return uml.Buffer.String()
	//return uml.start + uml.diagram.Contents() + uml.end
}

type PumlModel struct {
	Dependencies []Dependency
	Packages     []Package
}

type Package struct {
	Name    string
	Version string
}

// Represents a dependency From one thing To another
// E.g. "salt-master" -> "salt-common"
type Dependency struct {
	From        string
	FromVersion string
	To          string
	ToVersion   string
}
