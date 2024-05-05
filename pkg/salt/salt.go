package salt

import (
	"com/alexander/debendency/pkg"
	"os"
	"text/template"
)

// Functions needed for:
// First model in tree; supports online *and* offline installation
// Offline requires dependencies
// Next models only support offline installation
// Requires dependencies
func ToSaltDefinition(model *pkg.PackageModel) {

	saltTemplate := `
{{.Name}}:
  pkg.installed:
  {% if salt['grains.get']('offline', False) == True %}
    - sources:
      - {{.Name}}: "salt://{{.Filepath}}"
    - refresh: False
      {{- if .Dependencies}}
    - require:
      {{end}}
      {{- range .Dependencies}}- pkg: {{.Name}}{{end}}
  {% else %}
    - pkgs:
      - {{.Name}}: "{{.Version}}"
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
	err = tmpl.Execute(os.Stdout, model)
	if err != nil {
		panic(err)
	}
}
