package salt

import (
	"com/alexander/debendency/pkg"
	"io"
	"text/template"
)

// Functions needed for:
// First model in tree; supports online *and* offline installation
// Offline requires dependencies
// Next models only support offline installation
// Requires dependencies

func RootModelToTemplate(model *pkg.PackageModel, template string) {

}

func DependencyToTemplate(model *pkg.PackageModel, template string) {

}

func RootPackageToSaltDefinition(writer io.Writer, model *pkg.PackageModel, conf *pkg.Config) {

	// Loop through and select only packages that aren't installed if the -e flag has been passed
	var dependenciesNotInstalled = make([]*pkg.PackageModel, 0)
	for _, dependency := range model.OrderedDependencies {
		if conf.ExcludeInstalledPackages && !dependency.IsInstalled {
			dependenciesNotInstalled = append(dependenciesNotInstalled, dependency)
		}
	}
	model.OrderedDependencies = dependenciesNotInstalled

	saltTemplate := `
{{.Name}}:
  pkg.installed:
  {% if salt['grains.get']('offline', False) == True %}
    - sources:
      - {{.Name}}: "salt://{{.Filepath}}"
    - refresh: False
      {{- if .OrderedDependencies}}
    - require:{{end}}
      {{- range .OrderedDependencies}}
      - pkg: {{.Name}}{{end}}
  {% else %}
    - pkgs:
      - {{.Name}}: "{{.Version}}"
    - refresh: True
  {% endif %}
`
	//Dependency modelled here by reverse reference e.g. libjq1 required by jq
	//- require_in:
	//	- pkg: samba-libs
	// Dependency modelled here by forward reference e.g. jq requires libjq1
	//- require:
	//	- pkg: libavahi-client3
	tmpl, err := template.New("salt").Parse(saltTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(writer, model)
	if err != nil {
		panic(err)
	}
}

func DependenciesToSaltDefinitions(writer io.Writer, modelList []*pkg.PackageModel, conf *pkg.Config) {

	// Loop through and select only packages that aren't installed if the -e flag has been passed
	var packageModelsNotInstalled = make([]*pkg.PackageModel, 0)
	for _, packageModel := range modelList {
		if conf.ExcludeInstalledPackages && !packageModel.IsInstalled {
			var dependenciesNotInstalled = make([]*pkg.PackageModel, 0)
			for _, dependency := range packageModel.OrderedDependencies {
				if !dependency.IsInstalled {
					dependenciesNotInstalled = append(dependenciesNotInstalled, dependency)
				}
			}
			packageModel.OrderedDependencies = dependenciesNotInstalled
			packageModelsNotInstalled = append(packageModelsNotInstalled, packageModel)
		}
	}

	saltTemplate := `
{{- if . }}
{% if salt['grains.get']('offline', False) == True %}
{{- range . }}
{{.Name}}:
  pkg.installed:
    - sources:
      - {{.Name}}: "salt://{{.Filepath}}"
    - refresh: False
      {{- if .OrderedDependencies}}
    - require:{{end}}
      {{- range .OrderedDependencies}}
      - pkg: {{.Name}}{{end}}
{{end}}
{% endif %}{{end}}
`
	//Dependency modelled here by reverse reference e.g. libjq1 required by jq
	//- require_in:
	//	- pkg: samba-libs
	// Dependency modelled here by forward reference e.g. jq requires libjq1
	//- require:
	//	- pkg: libavahi-client3
	tmpl, err := template.New("salt").Parse(saltTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(writer, packageModelsNotInstalled)
	if err != nil {
		panic(err)
	}
}
