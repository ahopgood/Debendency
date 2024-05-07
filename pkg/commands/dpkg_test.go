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
			testFile, err := os.ReadFile("internal/dpkg/not-found.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(0))
		})
	})

	When("Package file has no dependencies", func() {
		It("returns empty array", func() {
			testFile, err := os.ReadFile("internal/dpkg/aglfn.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(0))
		})
	})

	When("Package file has dependencies already installed", func() {
		It("returns populated array", func() {
			testFile, err := os.ReadFile("internal/dpkg/dos2unix.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(1))
			Expect(dependencies[0]).To(Equal("libc6"))

		})
	})

	When("Package file has dependencies some installed", func() {
		It("returns populated array", func() {
			testFile, err := os.ReadFile("internal/dpkg/salt-minion.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(4))
			Expect(dependencies[0]).To(Equal("bsdmainutils"))
			Expect(dependencies[1]).To(Equal("dctrl-tools"))
			Expect(dependencies[2]).To(Equal("salt-common"))
			Expect(dependencies[3]).To(Equal("python3"))
		})
	})

	When("Package file has pre-dependencies", func() {
		It("returns populated array", func() {
			testFile, err := os.ReadFile("internal/dpkg/virtualbox-6.0.output")
			Expect(err).ToNot(HaveOccurred())

			dpkg := commands.Dpkger{}
			dependencies := dpkg.ParseDependencies(string(testFile))
			Expect(len(dependencies)).To(Equal(30))
			Expect(dependencies).Should(ContainElements("psmisc",
				"adduser", "libc6", "libcurl4", "libdevmapper1.02.1", "libgcc1",
				"libgl1", "libopus0", "libpng16-16", "libqt5core5a", "libqt5gui5",
				"libqt5opengl5", "libqt5printsupport5", "libqt5widgets5", "libqt5x11extras5",
				"libsdl1.2debian", "libssl1.1", "libstdc++6", "libvpx5", "libx11-6",
				"libxcb1", "libxcursor1", "libxext6", "libxml2", "libxmu6", "libxt6",
				"zlib1g", "python", "python", "python",
			))
		})
	})
})
