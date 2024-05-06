package puml

import (
	"com/alexander/debendency/pkg"
	"fmt"
)

func GenerateDiagram(config *pkg.Config, modelMap map[string]*pkg.PackageModel) Uml {
	dependencies := make([]Dependency, 0)
	// Loop through the model dependencies?
	for key, fromModel := range modelMap {
		fmt.Printf("%s %#v\n", key, fromModel)
		for _, toModel := range fromModel.Dependencies {
			if config.ExcludeInstalledPackages && (fromModel.IsInstalled || toModel.IsInstalled) {
				//fmt.Printf()
				//Log output separately from this model
			} else {
				fmt.Printf("From %s to %s\n", fromModel.Name, toModel.Name)
				dependencies = append(dependencies, Dependency{
					From: fromModel.Name,
					To:   toModel.Name,
				})
			}
		}
	}
	puml := NewUml(
		NewDigraph(dependencies),
	)
	return puml
}

type Uml struct {
	start   string
	diagram UmlDiagram
	end     string
}

func NewUml(umlDiagram UmlDiagram) Uml {
	return Uml{
		start:   "@startuml\n",
		diagram: umlDiagram,
		end:     "@enduml",
	}
}

func (uml Uml) Contents() string {
	return uml.start + uml.diagram.Contents() + uml.end
}

type UmlDiagram interface {
	Contents() string
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

// Represents a dependency From one thing To another
// E.g. "salt-master" -> "salt-common"
type Dependency struct {
	From string
	To   string
}

// Uml Diagram implementation
func (d Digraph) Contents() string {
	fmt.Println("Building diagraph contents")
	output := d.start + "\n"
	for _, value := range d.dependencies {
		output = output + "\t" + "\"" + value.From + "\"" + " -> " + "\"" + value.To + "\"" + "\n"

	}
	output = output + "\n" + d.end + "\n"
	return output
}
