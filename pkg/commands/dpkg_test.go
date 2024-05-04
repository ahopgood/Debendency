package commands_test

import (
	"com/alexander/debendency/pkg/commands"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ginkgo --focus Install
var _ = Describe("Install", func() {

	When("Package file doesn't exist", func() {
		It("returns empty array", func() {
			testFile, err := os.ReadFile("internal/install/not-found.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(0))
		})
	})

	When("Package file has no dependencies", func() {
		It("returns empty array", func() {
			testFile, err := os.ReadFile("internal/install/aglfn.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(0))
		})
	})

	When("Package file has dependencies already installed", func() {
		It("returns populated array", func() {
			testFile, err := os.ReadFile("internal/install/dos2unix.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(1))
			Expect(dependencies[0]).To(Equal("libc6 (>= 2.4)"))

		})
	})

	When("Package file has dependencies some installed", func() {
		It("returns populated array", func() {
			testFile, err := os.ReadFile("internal/install/salt-minion.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(4))
			Expect(dependencies[0]).To(Equal("bsdmainutils"))
			Expect(dependencies[1]).To(Equal("dctrl-tools"))
			Expect(dependencies[2]).To(Equal("salt-common (= 3004.2+ds-1)"))
			Expect(dependencies[3]).To(Equal("python3:any"))
		})
	})
})
