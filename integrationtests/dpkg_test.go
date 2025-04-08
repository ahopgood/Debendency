package integrationtests_test

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/commands"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"log/slog"
	"os"
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

	When("dpkg -I", func() {
		BeforeEach(func() {
			log.SetFlags(log.LstdFlags)
			slog.SetLogLoggerLevel(slog.LevelDebug)
		})
		var (
			packageFile string
		)
		query := commands.Dpkger{
			Cmd: commands.LinuxCommand{},
		}
		apt := commands.Apter{
			Cmd: commands.LinuxCommand{},
		}
		When("Package installed", func() {
			It("Should return true", func() {
				downloadOutput, _, err := apt.DownloadPackage("samba")
				Expect(err).ToNot(HaveOccurred())

				fmt.Println("Download command output")
				fmt.Println(downloadOutput)

				packageModel := pkg.PackageModel{}
				packageModel.GetPackageFilename(downloadOutput)

				packageFile = packageModel.Filepath
				dependencies := query.IdentifyDependencies(packageModel.Filepath)
				Expect(len(dependencies)).To(Equal(25))
			})
		})
		AfterEach(func() {
			fmt.Printf("Remove package file %s", packageFile)
			os.Remove(packageFile)
		})
	})
})
