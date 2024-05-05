package salt_test

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/salt"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Salt", func() {

	When("We have a model", func() {
		dos2unix := &pkg.PackageModel{
			Name:     "dos2unix",
			Filepath: "dos2unix_7.4.0-2_amd64.deb",
			Version:  "7.4.0-2",
		}
		It("Should create a valid salt .sls file", func() {
			salt.ToSaltDefinition(dos2unix)
			//Assert structure
			//Assert no dependency
		})
	})

	When("We have a model with a dependency", func() {
		jqlib1 := &pkg.PackageModel{
			Name:     "jqlib1",
			Filepath: "jqlib1_1.6-1ubuntu0.20.04.1_amd64.deb",
			Version:  "1.6-1ubuntu0.20.04.1",
		}
		jq := &pkg.PackageModel{
			Name:     "jq",
			Filepath: "jq_1.6-1ubuntu0.20.04.1_amd64.deb",
			Version:  "1.6-1ubuntu0.20.04.1",
			Dependencies: map[string]*pkg.PackageModel{
				"jqlib1": jqlib1,
			},
		}
		It("Should create a valid salt .sls file", func() {
			salt.ToSaltDefinition(jq)
			//Assert structure
			//Assert there is a dependency between the two declarations

			expectedSalt := `
jq:
  pkg.installed:
  {% if salt['grains.get']('offline', False) == True %}
    - sources:
      - jq: "salt://jq_1.6-1ubuntu0.20.04.1_amd64.deb"
    - refresh: False
    - require:
      - pkg: jqlib1
  {% else %}
    - pkgs:
      - jq: "1.6-1ubuntu0.20.04.1"
  {% endif %}
`
		})

	})
})
