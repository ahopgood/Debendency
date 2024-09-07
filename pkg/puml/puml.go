package puml

import (
	"bytes"
	"com/alexander/debendency/pkg"
	"fmt"
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
		fmt.Printf("%s %#v\n", fromModel.Name, fromModel)

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
				fmt.Printf("From %s to %s\n", fromModel.Name, toModel.Name)
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
	puml := NewUml(
		NewDigraph(dependencies),
	)

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
	start   string
	diagram UmlDiagram
	end     string
	Buffer  bytes.Buffer
}

func NewUml(umlDiagram UmlDiagram) Uml {
	return Uml{
		start:   "@startuml\n",
		diagram: umlDiagram,
		end:     "@enduml",
	}
}

func (uml Uml) Contents() string {
	return uml.Buffer.String()
	//return uml.start + uml.diagram.Contents() + uml.end
}

type Digraph struct {
	start        string
	dependencies []Dependency
	end          string
}

func NewDigraph(dependencies []Dependency) Digraph {
	return Digraph{
		start:        "digraph test {\n",
		dependencies: dependencies,
		end:          "}\n",
	}
}

// Uml Diagram implementation
func (d Digraph) Contents() string {
	fmt.Println("Building diagraph contents")
	output := d.start + "\n"
	for _, value := range d.dependencies {
		output = output + "\t" + "\"" + value.From + "\\n(" + value.FromVersion + ")\"" +
			" -> " + "\"" + value.To + "\\n(" + value.ToVersion + ")\"\n"
	}
	output = output + "\n" + d.end + "\n"
	return output
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
