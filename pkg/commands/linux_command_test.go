package commands_test

import (
	"com/alexander/debendency/pkg/commands"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commands", func() {
	// Possibly create a table test here using: command, standard out, standard error, exit status?
	When("Command fails", func() {

		cmd := commands.LinuxCommand{}

		It("returns error, non-zero exit code and standard error", func() {
			// vagrant@vagrant:~$ apt download test
			// E: Unable to locate package test
			// echo $?
			// 100

			var errorMessage = "\nWARNING: apt does not have a stable CLI interface. Use with caution in scripts.\n\nE: Unable to locate package test\n"

			output, statusCode, err := cmd.Command("apt", "download", "test")

			Expect(err).To(HaveOccurred())
			Expect(statusCode).To(Equal(100))
			Expect(output).To(Equal(errorMessage))
		})

	})

	When("Command succeeds", func() {
		cmd := commands.LinuxCommand{}

		output, statusCode, err := cmd.Command("apt", "download", "dos2unix")

		It("returns zero exit code and standard output", func() {
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
