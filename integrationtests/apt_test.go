package integrationtests

import (
	"com/alexander/debendency/pkg/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Apt Download", func() {
	When("Using LinuxCommand", func() {
		// Possibly create a table test here using: command, standard out, standard error, exit status?
		When("Package not recognised", func() {
			a := commands.Apter{
				Cmd: commands.LinuxCommand{},
			}

			It("fails with error message", func() {
				// vagrant@vagrant:~$ apt download test
				// E: Unable to locate package test
				// echo $?
				// 100

				var errorMessage = "\nWARNING: apt does not have a stable CLI interface. Use with caution in scripts.\n\nE: Unable to locate package test\n"

				output, statusCode, err := a.DownloadPackage("test")

				Expect(err).To(HaveOccurred())
				Expect(statusCode).To(Equal(100))
				Expect(output).To(Equal(errorMessage))
			})

		})

		When("Package recognised", func() {
			a := commands.Apter{
				Cmd: commands.LinuxCommand{},
			}

			It("Downloads the debian package file", func() {
				output, statusCode, err := a.DownloadPackage("dos2unix")
				// var successMessage = `
				// Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]
				// Fetched 374 kB in 0s (4,447 kB/s)
				// `

				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(0))
				Expect(output).Should(ContainSubstring("dos2unix"))
			})
		})

		When("Package already downloaded", func() {
			a := commands.Apter{
				Cmd: commands.LinuxCommand{},
			}

			It("Does not download the debian package file", func() {
				_, _, _ = a.DownloadPackage("dos2unix")
				output, statusCode, err := a.DownloadPackage("dos2unix")

				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(0))
				Expect(output).Should(ContainSubstring(""))
			})
		})
		AfterEach(func() {
			os.Remove("dos2unix_7.4.0-2_amd64.deb")
		})
	})
})
