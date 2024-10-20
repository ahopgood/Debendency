package commands_test

import (
	"com/alexander/debendency/pkg/commands"
	"com/alexander/debendency/pkg/commands/internal"
	"log"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ginkgo --focus dpkg
var _ = Describe("dpkg", func() {
	BeforeEach(func() {
		log.SetFlags(log.LstdFlags)
		slog.SetLogLoggerLevel(slog.LevelDebug)
	})
	Describe("Dpkger", func() {
		When("Package file doesn't exist", func() {
			It("returns empty array", func() {
				testFile, err := os.ReadFile("internal/dpkg/not-found.output")
				Expect(err).ToNot(HaveOccurred())

				//mockCommand := internal.FakeCommand{}
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

		When("Samba", func() {
			It("returns populated array", func() {
				testFile, err := os.ReadFile("internal/dpkg/samba.output")
				Expect(err).ToNot(HaveOccurred())

				dpkg := commands.Dpkger{}
				dependencies := dpkg.ParseDependencies(string(testFile))
				Expect(len(dependencies)).To(Equal(25))
				Expect(dependencies).Should(ContainElements("adduser", "libpam-modules",
					"libpam-runtime", "lsb-base", "procps", "python3", "python3-dnspython",
					"python3-samba", "samba-common", "samba-common-bin", "tdb-tools", "python3",
					"python3", "libbsd0", "libc6", "libgnutls30", "libldb2", "libpopt0", "libpython3.8",
					"libtalloc2", "libtasn1-6", "libtdb1", "libtevent0", "libwbclient0", "samba-libs"))
			})
		})
	})

	Describe("Query", func() {
		DescribeTable("Package is installed",
			func(output string) {
				command := &internal.FakeCommand{}
				command.CommandReturns(output, 0, nil)

				query := commands.Query{
					Cmd: command,
				}
				installed := query.IsInstalled("samba")
				Expect(installed).To(BeTrue())
			},
			Entry("tzdata", "tzdata      2023c-0ubuntu0.20.04.0"),
			Entry("dos2unix", "dos2unix        7.4.0-2"),
			Entry("libicu66", "libicu66:amd64      66.1-2ubuntu2.1"),
		)

		When("Package is not installed", func() {
			It("returns false", func() {
				command := &internal.FakeCommand{}
				command.CommandReturns("dos2unix", 0, nil)

				query := commands.Query{
					Cmd: command,
				}
				installed := query.IsInstalled("samba")
				Expect(installed).To(BeFalse())
			})
		})
		When("Package is not recognised", func() {
			It("returns false", func() {
				command := &internal.FakeCommand{}
				command.CommandReturns("dpkg-query: no packages found matching python3-dnspython", 0, nil)

				query := commands.Query{
					Cmd: command,
				}
				installed := query.IsInstalled("samba")
				Expect(installed).To(BeFalse())
			})
		})
	})
})
