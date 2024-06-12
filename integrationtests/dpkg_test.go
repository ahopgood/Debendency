package integrationtests

import (
	"com/alexander/debendency/pkg/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dpkg", func() {

	When("dpkg-query", func() {
		//dpkg-query -W libgcc-s1
		//libgcc-s1:amd64 10.5.0-1ubuntu1~20.04
		//vagrant@vagrant:/vagrant$ echo $?
		//0
		query := commands.Query{
			Cmd: commands.LinuxCommand{},
		}
		When("Package installed", func() {
			It("Should return true", func() {
				installed := query.IsInstalled("libc6")
				Expect(installed).To(BeTrue())
			})
		})

		//dpkg-query -W libgcc-s1f
		//dpkg-query: no packages found matching libgcc-s1f
		//vagrant@vagrant:/vagrant$ echo $?
		//1
		When("Package not installed", func() {
			It("Should return false", func() {
				installed := query.IsInstalled("unknown")
				Expect(installed).To(BeFalse())
			})
		})
	})
})
