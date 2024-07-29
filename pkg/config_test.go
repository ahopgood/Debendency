package pkg_test

import (
	"com/alexander/debendency/pkg"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	When("package (-p) flag set", func() {
		It("Should set Config.PackageName", func() {
			config, output, err := pkg.ParseFlags("", []string{"-p", "test.deb"})
			Expect(config.PackageName).To(Equal("test.deb"))
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("generate salt (-s) flag set", func() {
		It("Should set Config.GenerateSalt", func() {
			config, output, err := pkg.ParseFlags("", []string{"-s", "-p", "test.deb"})
			Expect(config.GenerateSalt).To(BeTrue())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("diagram (-d)  flag set", func() {
		It("Should set Config.GenerateDiagram", func() {
			config, output, err := pkg.ParseFlags("", []string{"-d", "-p", "test.deb"})
			Expect(config.GenerateDiagram).To(BeTrue())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("installer location (-o)  flag set", func() {
		It("Should set Config.InstallerLocation", func() {
			config, output, err := pkg.ParseFlags("", []string{"-o", "somedir", "-p", "test.deb"})
			Expect(config.InstallerLocation).To(Equal("somedir"))
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("exclude installed packages (-e)  flag set", func() {
		It("Should set Config.ExcludeInstalledPackages", func() {
			config, output, err := pkg.ParseFlags("", []string{"-e", "-p", "test.deb"})
			Expect(config.ExcludeInstalledPackages).To(BeTrue())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("no flags except package set", func() {
		It("Should provide defaults", func() {
			config, output, err := pkg.ParseFlags("", []string{"-p", "test.deb"})
			expectedConfig := pkg.Config{
				PackageName:              "test.deb",
				GenerateSalt:             false,
				GenerateDiagram:          false,
				InstallerLocation:        "~/.debendency/cache",
				ExcludeInstalledPackages: false,
			}
			Expect(cmp.Diff(config, &expectedConfig)).To(BeEmpty())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("package flag not set", func() {
		It("should trigger the help output", func() {
			config, output, err := pkg.ParseFlags("dummy program", []string{})

			By("not initialising the config with defaults", func() {
				Expect(config).To(BeNil())
			})

			By("printing out usage", func() {
				Expect(output).To(Not(BeEmpty()))
				//fmt.Print(output)
			})

			By("returning the flag error", func() {
				Expect(err.Error()).To(Equal("flag: need to specify a package name via -p flag"))
			})
		})
	})

	When("help (-h) flag set", func() {
		It("should trigger the help output", func() {
			config, output, err := pkg.ParseFlags("", []string{"-h"})

			By("not initialising the config with defaults", func() {
				Expect(config).To(BeNil())
			})

			By("printing out usage", func() {
				Expect(output).To(Not(BeEmpty()))
			})

			By("returning the flag error", func() {
				Expect(err.Error()).To(Equal("flag: help requested"))
			})
		})
	})

	When("unknown flag (-z) flag set", func() {
		It("should trigger the help output", func() {
			config, output, err := pkg.ParseFlags("", []string{"-z"})

			By("not initialising the config with defaults", func() {
				Expect(config).To(BeNil())
			})

			By("printing out usage", func() {
				Expect(output).To(Not(BeEmpty()))
			})

			By("returning the flag error", func() {
				Expect(err.Error()).To(Equal("flag provided but not defined: -z"))
			})
		})
	})
})
