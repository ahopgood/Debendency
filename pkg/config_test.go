package pkg_test

import (
	"com/alexander/debendency/pkg"
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
			config, output, err := pkg.ParseFlags("", []string{"-s", "test.deb"})
			Expect(config.GenerateSalt).To(BeTrue())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("diagram (-d)  flag set", func() {
		It("Should set Config.GenerateDiagram", func() {
			config, output, err := pkg.ParseFlags("", []string{"-d", "test.deb"})
			Expect(config.GenerateDiagram).To(BeTrue())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("installer location (-o)  flag set", func() {
		It("Should set Config.InstallerLocation", func() {
			config, output, err := pkg.ParseFlags("", []string{"-o", "somedir"})
			Expect(config.InstallerLocation).To(Equal("somedir"))
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})

	When("exclude installed packages (-e)  flag set", func() {
		It("Should set Config.ExcludeInstalledPackages", func() {
			config, output, err := pkg.ParseFlags("", []string{"-e", "test.deb"})
			Expect(config.ExcludeInstalledPackages).To(BeTrue())
			Expect(output).To(BeEmpty())
			Expect(err).To(BeNil())
		})
	})
})
