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

func ToSaltDefinition(writer io.Writer, model *pkg.PackageModel) {

	saltTemplate := `
{{.Name}}:
  pkg.installed:
  {% if salt['grains.get']('offline', False) == True %}
    - sources:
      - {{.Name}}: "salt://{{.Filepath}}"
    - refresh: False
      {{- if .Dependencies}}
    - require:{{end}}
      {{- range .Dependencies}}
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

func ToSaltDefinitions(writer io.Writer, modelList []*pkg.PackageModel) {

	saltTemplate := `
{{- if . }}
{% if salt['grains.get']('offline', False) == True %}
{{- range . }}
{{.Name}}:
  pkg.installed:
    - sources:
      - {{.Name}}: "salt://{{.Filepath}}"
    - refresh: False
      {{- if .Dependencies}}
    - require:{{end}}
      {{- range .Dependencies}}
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
	err = tmpl.Execute(writer, modelList)
	if err != nil {
		panic(err)
	}
}
