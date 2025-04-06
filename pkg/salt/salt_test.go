package salt_test

import (
	"bytes"
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/salt"
	"fmt"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Salt", func() {

	var conf *pkg.Config

	BeforeEach(func() {
		conf = &pkg.Config{
			ExcludeInstalledPackages: true,
		}
	})

	When("We have a model", func() {
		dos2unix := &pkg.PackageModel{
			Name:     "dos2unix",
			Filepath: "dos2unix_7.4.0-2_amd64.deb",
			Version:  "7.4.0-2",
		}
		It("Should create a valid salt .sls file", func() {
			var b bytes.Buffer
			salt.RootPackageToSaltDefinition(&b, dos2unix, conf)

			expectedSalt := `
dos2unix:
  pkg.installed:
  {% if salt['grains.get']('offline', False) == True %}
    - sources:
      - dos2unix: "salt://dos2unix_7.4.0-2_amd64.deb"
    - refresh: False
  {% else %}
    - pkgs:
      - dos2unix: "7.4.0-2"
    - refresh: True
  {% endif %}
`
			fmt.Println("Expected")
			fmt.Println(expectedSalt)
			fmt.Println("Actual")
			fmt.Println(b.String())

			diff := cmp.Diff(b.String(), expectedSalt)
			Expect(diff).To(BeEmpty())
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
			OrderedDependencies: []*pkg.PackageModel{
				jqlib1,
			},
		}
		It("Should create a valid salt .sls file", func() {
			var b bytes.Buffer
			//Assert structure
			salt.RootPackageToSaltDefinition(&b, jq, conf)
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
    - refresh: True
  {% endif %}
`
			fmt.Println("Expected")
			fmt.Println(expectedSalt)
			fmt.Println("Actual")
			fmt.Println(b.String())

			diff := cmp.Diff(b.String(), expectedSalt)
			Expect(diff).To(BeEmpty())
		})
	})

	When("Modelling dependency list", func() {
		jqlib1 := &pkg.PackageModel{
			Name:     "jqlib1",
			Filepath: "jqlib1_1.6-1ubuntu0.20.04.1_amd64.deb",
			Version:  "1.6-1ubuntu0.20.04.1",
		}
		jq := &pkg.PackageModel{
			Name:     "jq",
			Filepath: "jq_1.6-1ubuntu0.20.04.1_amd64.deb",
			Version:  "1.6-1ubuntu0.20.04.1",
			OrderedDependencies: []*pkg.PackageModel{
				jqlib1,
			},
		}
		It("Should create a valid salt .sls file", func() {
			var b bytes.Buffer
			salt.DependenciesToSaltDefinitions(&b, []*pkg.PackageModel{jq, jqlib1}, conf)
			//Assert structure
			//Assert there is a dependency between the two declarations

			expectedSalt := `
{% if salt['grains.get']('offline', False) == True %}
jq:
  pkg.installed:
    - sources:
      - jq: "salt://jq_1.6-1ubuntu0.20.04.1_amd64.deb"
    - refresh: False
    - require:
      - pkg: jqlib1

jqlib1:
  pkg.installed:
    - sources:
      - jqlib1: "salt://jqlib1_1.6-1ubuntu0.20.04.1_amd64.deb"
    - refresh: False

{% endif %}
`
			fmt.Println("Expected")
			fmt.Println(expectedSalt)
			fmt.Println("Actual")
			fmt.Println(b.String())
			diff := cmp.Diff(b.String(), expectedSalt)
			Expect(diff).To(BeEmpty())
		})
	})

	When("Excluding installed dependencies", func() {
		libonig5 := &pkg.PackageModel{
			Name:        "libonig5",
			Filepath:    "libonig5-6.9.4-1_amd64.deb",
			Version:     "6.9.4-1",
			IsInstalled: true,
		}

		jqlib1 := &pkg.PackageModel{
			Name:     "jqlib1",
			Filepath: "jqlib1_1.6-1ubuntu0.20.04.1_amd64.deb",
			Version:  "1.6-1ubuntu0.20.04.1",
			OrderedDependencies: []*pkg.PackageModel{
				libonig5,
			},
		}
		jq := &pkg.PackageModel{
			Name:     "jq",
			Filepath: "jq_1.6-1ubuntu0.20.04.1_amd64.deb",
			Version:  "1.6-1ubuntu0.20.04.1",
			OrderedDependencies: []*pkg.PackageModel{
				jqlib1,
			},
		}
		It("Should exclude the installed dependency", func() {
			var b bytes.Buffer
			salt.DependenciesToSaltDefinitions(&b, []*pkg.PackageModel{jq, jqlib1}, conf)

			expectedSalt := `
{% if salt['grains.get']('offline', False) == True %}
jq:
  pkg.installed:
    - sources:
      - jq: "salt://jq_1.6-1ubuntu0.20.04.1_amd64.deb"
    - refresh: False
    - require:
      - pkg: jqlib1

jqlib1:
  pkg.installed:
    - sources:
      - jqlib1: "salt://jqlib1_1.6-1ubuntu0.20.04.1_amd64.deb"
    - refresh: False

{% endif %}
`
			fmt.Println("Expected")
			fmt.Println(expectedSalt)
			fmt.Println("Actual")
			fmt.Println(b.String())
			diff := cmp.Diff(b.String(), expectedSalt)
			Expect(diff).To(BeEmpty())
		})
	})
})
