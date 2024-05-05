package commands_test

import (
	"com/alexander/debendency/pkg/commands"
	"com/alexander/debendency/pkg/commands/internal"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apt Download", func() {

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

			It("Downloads the debian package file", func() {
				output, statusCode, err := a.DownloadPackage("dos2unix")
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

		When("Package already downloaded", func() {
			successMessage := ""

			command := &internal.FakeCommand{}
			command.CommandReturns(successMessage, 0, nil)

			a := commands.Apter{
				Cmd: command,
			}

			It("Does not download the debian package file", func() {
				output, statusCode, err := a.DownloadPackage("dos2unix")

				Expect(command.CommandCallCount()).To(Equal(1))
				arg, _ := command.CommandArgsForCall(0)
				fmt.Printf("%#v\n", arg)
				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(0))
				Expect(output).Should(ContainSubstring(""))
			})
		})
	})

	//apt list dos2unix
	//Listing... Done
	//dos2unix/focal 7.4.0-2 amd64
})
