package commands_test

import (
	"com/alexander/debendency/pkg/commands"
	"com/alexander/debendency/pkg/commands/internal"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Download", func() {
	// var mockCommand =
	// When("OS not Linux", func() {
	// 	writer := strings.Builder{}
	// 	command := &internal.FakeCommand{}

	// 	command.CommandReturns(writer, 130, nil)
	// 	It("returns error", func() {
	// 		// vagrant@vagrant:~$ apt download test
	// 		// E: Unable to locate package test
	// 		// echo $?
	// 		// 130
	// 		a := salt.Apter{
	// 			Cmd: command,
	// 		}
	// 		_, _, err := a.DownloadPackage("dos2unix")
	// 		Expect(command.CommandCallCount()).To(Equal(1))
	// 		Expect(err).To(HaveOccurred())

	// 	})
	// 	AfterEach(func() {
	// 		os.Remove("dos2unix_7.4.0-2_amd64.deb")
	// 	})
	// })

	// Possibly create a table test here using: command, standard out, standard error, exit status?
	When("Using Fake command", func() {
		When("Package not recognised", func() {
			errorMessage := "E: Unable to locate package test"

			// writer := strings.Builder{}
			// writer.Write([]byte(errorMessage))
			command := &internal.FakeCommand{}

			command.CommandReturns(errorMessage, 130, nil)

			a := commands.Apter{
				Cmd: command,
			}

			It("fails with error message", func() {
				// vagrant@vagrant:~$ apt download test
				// E: Unable to locate package test
				// echo $?
				// 130

				output, statusCode, err := a.DownloadPackage("test")

				Expect(command.CommandCallCount()).To(Equal(1))
				Expect(err).Error().NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(130))
				Expect(output).To(Equal(errorMessage))
			})

		})

		When("Package recognised", func() {
			successMessage := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]\nFetched 374 kB in 0s (4,447 kB/s)"

			// writer := strings.Builder{}
			// writer.Write([]byte(successMessage))
			command := &internal.FakeCommand{}

			command.CommandReturns(successMessage, 0, nil)

			a := commands.Apter{
				Cmd: command,
			}

			output, statusCode, err := a.DownloadPackage("dos2unix")

			It("Downloads the debian package file", func() {
				// vagrant@vagrant:~$ apt download dos2unix
				// Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]
				// Fetched 374 kB in 0s (4,447 kB/s)
				// echo $?
				// 0

				Expect(command.CommandCallCount()).To(Equal(1))
				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(0))
				Expect(output).To(Equal(successMessage))
			})
		})
	})

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

			output, statusCode, err := a.DownloadPackage("dos2unix")

			It("Downloads the debian package file", func() {
				// var successMessage = `
				// Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]
				// Fetched 374 kB in 0s (4,447 kB/s)
				// `

				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(0))
				Expect(output).Should(ContainSubstring("dos2unix"))
			})
		})
		AfterEach(func() {
			os.Remove("dos2unix_7.4.0-2_amd64.deb")
		})
	})
})
