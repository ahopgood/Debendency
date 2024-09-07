package puml_test

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/puml"
	"fmt"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Puml", func() {

	jq := &pkg.PackageModel{
		Name:        "jq",
		Version:     "1.6-1ubuntu0.20.04.1",
		Filepath:    "jq_1.6-1ubuntu0.20.04.1_amd64.deb",
		IsInstalled: false,
	}
	libjq1 := &pkg.PackageModel{
		Name:        "libjq1",
		Version:     "1.6-1ubuntu0.20.04.1",
		Filepath:    "libjq1_1.6-1ubuntu0.20.04.1_amd64.deb",
		IsInstalled: false,
	}
	libonig5 := &pkg.PackageModel{
		Name:        "libonig5",
		Version:     "6.9.4-1",
		Filepath:    "libonig5_6.9.4-1_amd64.deb",
		IsInstalled: false,
	}

	libc6 := &pkg.PackageModel{
		Name:        "libc6",
		Version:     "2.31-0ubuntu9.15",
		Filepath:    "libc6_2.31-0ubuntu9.15_amd64.deb",
		IsInstalled: true,
	}

	libgcc_s1 := &pkg.PackageModel{
		Name:        "libgcc-s1",
		Version:     "10.5.0-1ubuntu1~20.04",
		Filepath:    "libgcc-s1_10.5.0-1ubuntu1~20.04_amd64.deb",
		IsInstalled: true,
	}

	gcc_10_base := &pkg.PackageModel{
		Name:         "gcc-10-base",
		Version:      "10.5.0-1ubuntu1~20.04",
		Filepath:     "gcc-10-base_10.5.0-1ubuntu1~20.04_amd64.deb",
		IsInstalled:  true,
		Dependencies: map[string]*pkg.PackageModel{},
	}

	libcrypt1 := &pkg.PackageModel{
		Name:         "libcrypt1",
		Version:      "1:4.4.10-10ubuntu4",
		Filepath:     "libcrypt1_1%3a4.4.10-10ubuntu4_amd64.deb",
		IsInstalled:  true,
		Dependencies: map[string]*pkg.PackageModel{},
	}

	modelMap := map[string]*pkg.PackageModel{
		"jq":          jq,
		"libjq1":      libjq1,
		"libonig5":    libonig5,
		"libc6":       libc6,
		"libgcc-s1":   libgcc_s1,
		"gcc-10-base": gcc_10_base,
		"libcrypt1":   libcrypt1,
	}

	jq.Dependencies = map[string]*pkg.PackageModel{
		"libjq1": libjq1,
		"libc6":  libc6,
	}
	jq.OrderedDependencies = []*pkg.PackageModel{
		libjq1,
		libc6,
	}

	libjq1.Dependencies = map[string]*pkg.PackageModel{
		"libonig5": libonig5,
		"libc6":    libc6,
	}
	libjq1.OrderedDependencies = []*pkg.PackageModel{
		libonig5,
		libc6,
	}

	libonig5.Dependencies = map[string]*pkg.PackageModel{
		"libc6": libc6,
	}
	libonig5.OrderedDependencies = []*pkg.PackageModel{
		libc6,
	}

	libc6.Dependencies = map[string]*pkg.PackageModel{
		"libgcc-s1": libgcc_s1,
		"libcrypt1": libcrypt1,
	}
	libc6.OrderedDependencies = []*pkg.PackageModel{
		libgcc_s1,
		libcrypt1,
	}

	libgcc_s1.Dependencies = map[string]*pkg.PackageModel{
		"gcc-10-base": gcc_10_base,
		"libc6":       libc6,
	}
	libgcc_s1.OrderedDependencies = []*pkg.PackageModel{
		gcc_10_base,
		libc6,
	}
	//Important to test we don't create an empty relation
	gcc_10_base.Dependencies = nil

	modelList := []*pkg.PackageModel{
		jq,
		libc6,
		libonig5,
		libjq1,
		libgcc_s1,
		gcc_10_base,
		libcrypt1,
	}

	When("Generate Diagram", func() {
		When("ModelMap is empty", func() {
			config := &pkg.Config{
				ExcludeInstalledPackages: false,
			}
			It("Should show an empty diagraph", func() {
				testFile, err := os.ReadFile("internal/EmptyDiagram.puml")
				Expect(err).ToNot(HaveOccurred())

				emptyModelMap := map[string]*pkg.PackageModel{}
				emptyModelList := []*pkg.PackageModel{}

				pumlDiagram := puml.GenerateDiagram(config, emptyModelMap, emptyModelList).Contents()
				Expect(pumlDiagram).To(Equal(string(testFile)))
				//Expect(cmp.Diff(pumlDiagram, string(testFile))).To(BeEmpty())
			})
		})

		When("Including Installed Packages", func() {
			config := &pkg.Config{
				ExcludeInstalledPackages: false,
			}
			It("Should include all relationships", func() {
				testFile, err := os.ReadFile("internal/IncludeAllDependencies.puml")
				Expect(err).ToNot(HaveOccurred())

				pumlDiagram := puml.GenerateDiagram(config, modelMap, modelList).Contents()
				fmt.Println("Expected")
				fmt.Println(string(testFile))
				fmt.Println("Actual")
				fmt.Println(pumlDiagram)
				//Expect(pumlDiagram).To(Equal(string(testFile)))
				diff := cmp.Diff(pumlDiagram, string(testFile))
				Expect(diff).To(BeEmpty())
			})
		})

		When("Excluding Installed Packages", func() {
			config := &pkg.Config{
				ExcludeInstalledPackages: true,
			}
			It("Should exclude all relationships where either origin or destination package is installed", func() {
				testFile, err := os.ReadFile("internal/ExcludeInstalledDependencies.puml")
				Expect(err).ToNot(HaveOccurred())

				pumlDiagram := puml.GenerateDiagram(config, modelMap, modelList).Contents()
				fmt.Println("Expected")
				fmt.Println(string(testFile))
				fmt.Println("Actual")
				fmt.Println(pumlDiagram)
				//Expect(pumlDiagram).To(Equal(string(testFile)))
				diff := cmp.Diff(pumlDiagram, string(testFile))
				Expect(diff).To(BeEmpty())
			})
		})

		When("Model has no dependencies", func() {
			config := &pkg.Config{
				ExcludeInstalledPackages: false,
			}

			docker := &pkg.PackageModel{
				Name:        "docker-ce",
				Version:     "18.09.7~3-0~ubuntu-xenial",
				Filepath:    "18.09.7~3-0~ubuntu-xenial_amd64.deb",
				IsInstalled: false,
			}

			modelMap := map[string]*pkg.PackageModel{
				"docker-ce": docker,
			}
			modelList := []*pkg.PackageModel{
				docker,
			}

			It("Should show root package", func() {
				testFile, err := os.ReadFile("internal/NoDependencies.puml")
				Expect(err).ToNot(HaveOccurred())

				pumlDiagram := puml.GenerateDiagram(config, modelMap, modelList).Contents()

				fmt.Println("Expected")
				fmt.Println(string(testFile))
				fmt.Println("Actual")
				fmt.Println(pumlDiagram)
				//Expect(pumlDiagram).To(Equal(string(testFile)))
				diff := cmp.Diff(pumlDiagram, string(testFile))
				Expect(diff).To(BeEmpty())
			})
		})
	})
})
