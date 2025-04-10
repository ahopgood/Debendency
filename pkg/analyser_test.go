package pkg_test

import (
	"com/alexander/debendency/pkg"
	"com/alexander/debendency/pkg/internal"
	"com/alexander/debendency/pkg/puml"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"log/slog"
	"os"
)

var _ = Describe("Analyser", func() {
	When("GetPackageFilename", func() {
		When("Downloaded Successfully", func() {
			It("Should construct file name correctly", func() {
				successMessage := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]\nFetched 374 kB in 0s (4,447 kB/s)"

				model := pkg.PackageModel{}
				model.GetPackageFilename(successMessage)
				Expect(model.Name).To(Equal("dos2unix"))
				Expect(model.Version).To(Equal("7.4.0-2"))
				Expect(model.Filepath).To(Equal("dos2unix_7.4.0-2_amd64.deb"))
			})
		})
		When("Already downloaded", func() {
			It("Should construct file name correctly", func() {
				successMessage := ""

				model := pkg.PackageModel{}
				model.GetPackageFilename(successMessage)
				Expect(model.Name).To(Equal(""))
				Expect(model.Version).To(Equal(""))
				Expect(model.Filepath).To(Equal(""))
			})
		})
		When("Package has a colon in its name", func() {
			It("Should web encode the colon as %3a in the filename", func() {
				log.SetFlags(log.LstdFlags)
				slog.SetLogLoggerLevel(slog.LevelDebug)
				//successMessage := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal-updates/main amd64 samba amd64 2:4.15.13+dfsg-0ubuntu0.20.04.8 [1,167 kB]\nFetched 1,167 kB in 0s (3,918 kB/s)"
				successMessage := `	Get:1 http://gb.archive.ubuntu.com/ubuntu focal-updates/main amd64 samba amd64 2:4.15.13+dfsg-0ubuntu0.20.04.8 [1,167 kB]
				Fetched 1,167 kB in 0s (9,447 kB/s)
				`
				model := pkg.PackageModel{}
				model.GetPackageFilename(successMessage)
				Expect(model.Name).To(Equal("samba"))
				Expect(model.Version).To(Equal("2:4.15.13+dfsg-0ubuntu0.20.04.8"))
				Expect(model.Filepath).To(Equal("samba_2%3a4.15.13+dfsg-0ubuntu0.20.04.8_amd64.deb"))
			})
		})
		When("apt provides a warning", func() {
			It("Should ignore the warning and find the package name", func() {
				var warningMessage = "\nWARNING: apt does not have a stable CLI interface. Use with caution in scripts.\n\nGet:1 http://gb.archive.ubuntu.com/ubuntu focal-updates/main amd64 samba amd64 2:4.15.13+dfsg-0ubuntu0.20.04.8 [1,167 kB]\nFetched 1,167 kB in 0s (12.7 MB/s)\n"

				log.SetFlags(log.LstdFlags)
				slog.SetLogLoggerLevel(slog.LevelDebug)
				model := pkg.PackageModel{}
				model.GetPackageFilename(warningMessage)
				Expect(model.Name).To(Equal("samba"))
				Expect(model.Version).To(Equal("2:4.15.13+dfsg-0ubuntu0.20.04.8"))
				Expect(model.Filepath).To(Equal("samba_2%3a4.15.13+dfsg-0ubuntu0.20.04.8_amd64.deb"))
			})
		})
	})
	AfterEach(func() {
		os.Remove("dos2unix_7.4.0-2_amd64.deb")
	})

	dos2unix := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]\nFetched 374 kB in 0s (4,447 kB/s)"
	jq := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal-updates/universe amd64 jq amd64 1.6-1ubuntu0.20.04.1 [50.2 kB]\nFetched 50.2 kB in 0s (1,666 kB/s)\n"
	jqFile := "jq_1.6-1ubuntu0.20.04.1_amd64.deb"
	libjq1 := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal-updates/universe amd64 libjq1 amd64 1.6-1ubuntu0.20.04.1 [121 kB]\nFetched 121 kB in 0s (2,898 kB/s)\n"
	libjq1File := "libjq1_1.6-1ubuntu0.20.04.1_amd64.deb"
	libonig5 := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 libonig5 amd64 6.9.4-1 [142 kB]\nFetched 142 kB in 0s (3,228 kB/s)\n"
	libonig5File := "libonig5_6.9.4-1_amd64.deb"
	libc6 := "Get:1 http://gb.archive.ubuntu.com/ubuntu focal-updates/main amd64 libc6 amd64 2.31-0ubuntu9.15 [2,723 kB]\nFetched 2,723 kB in 0s (20.9 MB/s)\n"
	libc6File := "libc6_2.31-0ubuntu9.15_amd64.deb"

	libgsoap := "Get:1 http://gb.archive.ubuntu.com/ubuntu jammy/universe amd64 libgsoap-2.8.117 amd64 2.8.117-2build1 [269 kB]\nFetched 269 kB in 0s (2,756 kB/s)\n"
	libgsoapFile := "libgsoap-2.8.117_2.8.117-2build1_amd64.deb"

	When("BuildPackage", func() {
		When("Package does not exist", func() {
			It("Produces nothing", func() {
				// Fail to download the debendency
				apter := &internal.FakeApt{}
				apter.DownloadPackageReturns("", 0, errors.New("not found"))

				dpkger := &internal.FakeDpkg{}
				dpkgQuery := &internal.FakeDpkgQuery{}

				config := &pkg.Config{
					ExcludeInstalledPackages: false,
				}

				packager := pkg.Analyser{
					Apt:    apter,
					Dpkg:   dpkger,
					Config: config,
					Query:  dpkgQuery,
				}

				modelMap := make(map[string]*pkg.PackageModel)
				modelList := make([]*pkg.PackageModel, 0)
				model := packager.BuildPackage("", modelMap, &modelList)
				By("Not adding a model to the map", func() {
					Expect(len(modelMap)).To(Equal(0))
				})
				By("Returning an empty model", func() {
					Expect(model.Name).To(BeEmpty())
					Expect(model.Version).To(BeEmpty())
					Expect(model.Filepath).To(BeEmpty())
				})
				By("Apt being only invoked once", func() {
					Expect(apter.DownloadPackageCallCount()).To(Equal(1))
				})
				By("Dpkg not being invoked", func() {
					Expect(dpkger.IdentifyDependenciesCallCount()).To(Equal(0))
				})
				By("DpkgQuery not being invoked", func() {
					Expect(dpkgQuery.IsInstalledCallCount()).To(Equal(0))
				})
			})
		})

		When("Package has no dependencies", func() {
			It("Produces a single model", func() {
				// Successfully download the debendency
				apter := &internal.FakeApt{}
				apter.DownloadPackageReturns(dos2unix, 0, nil)

				// No dependencies found via dpkg -I
				dpkger := &internal.FakeDpkg{}
				dpkger.IdentifyDependenciesReturns([]string{})

				dpkgQuery := &internal.FakeDpkgQuery{}

				config := &pkg.Config{
					ExcludeInstalledPackages: false,
				}

				packager := pkg.Analyser{
					Apt:    apter,
					Dpkg:   dpkger,
					Config: config,
					Query:  dpkgQuery,
				}

				modelMap := make(map[string]*pkg.PackageModel)
				modelList := make([]*pkg.PackageModel, 0)
				model := packager.BuildPackage("", modelMap, &modelList)

				By("Adding a model to the map", func() {
					Expect(len(modelMap)).To(Equal(1))
				})
				By("Returning a populated model", func() {
					Expect(model.Name).To(Equal("dos2unix"))
					Expect(model.Version).To(Equal("7.4.0-2"))
					Expect(model.Filepath).To(Equal("dos2unix_7.4.0-2_amd64.deb"))
				})
				By("Invoking Apt and Dpkg only once", func() {
					Expect(apter.DownloadPackageCallCount()).To(Equal(1))
					Expect(dpkger.IdentifyDependenciesCallCount()).To(Equal(1))
				})
				By("DpkgQuery not being invoked", func() {
					Expect(dpkgQuery.IsInstalledCallCount()).To(Equal(0))
				})
			})
		})

		When("Package has one dependency", func() {
			It("Produces two models", func() {
				// Successfully download the debendency
				apter := &internal.FakeApt{}
				apter.DownloadPackageReturnsOnCall(0, dos2unix, 0, nil)
				apter.DownloadPackageReturnsOnCall(1, libc6, 0, nil)

				// No dependencies found via dpkg -I
				dpkger := &internal.FakeDpkg{}
				dpkger.IdentifyDependenciesReturnsOnCall(0, []string{"libc6"})
				dpkger.IdentifyDependenciesReturnsOnCall(1, []string{})

				dpkgQuery := &internal.FakeDpkgQuery{}

				config := &pkg.Config{
					ExcludeInstalledPackages: false,
				}

				packager := pkg.Analyser{
					Apt:    apter,
					Dpkg:   dpkger,
					Config: config,
					Query:  dpkgQuery,
				}

				modelMap := make(map[string]*pkg.PackageModel)
				modelList := make([]*pkg.PackageModel, 0)

				_ = packager.BuildPackage("", modelMap, &modelList)

				By("Invoking Apt and Dpkg twice", func() {
					Expect(apter.DownloadPackageCallCount()).To(Equal(2))
					Expect(dpkger.IdentifyDependenciesCallCount()).To(Equal(2))
				})
				By("Adding two models to the map", func() {
					Expect(len(modelMap)).To(Equal(2))
				})
				By("Adding two models to the global list", func() {
					Expect(len(modelList)).To(Equal(2))
				})
				By("Returning dos2unix model", func() {
					model := modelMap["dos2unix"]
					Expect(model.Name).To(Equal("dos2unix"))
					Expect(model.Version).To(Equal("7.4.0-2"))
					Expect(model.Filepath).To(Equal("dos2unix_7.4.0-2_amd64.deb"))
				})
				By("Returning libc6 model", func() {
					model := modelMap["libc6"]
					Expect(model.Name).To(Equal("libc6"))
					Expect(model.Version).To(Equal("2.31-0ubuntu9.15"))
					Expect(model.Filepath).To(Equal("libc6_2.31-0ubuntu9.15_amd64.deb"))
				})
				By("Adding libc6 to the dos2unix dependencies map", func() {
					model := modelMap["dos2unix"]
					Expect(model.Dependencies["libc6"]).To(Not(BeNil()))
				})
				By("Adding libc6 to the dos2unix ordered dependencies list", func() {
					model := modelMap["dos2unix"]
					Expect(model.Dependencies["libc6"]).To(Not(BeNil()))
					Expect(len(model.OrderedDependencies)).To(Equal(1))
				})
				By("DpkgQuery not being invoked", func() {
					Expect(dpkgQuery.IsInstalledCallCount()).To(Equal(0))
				})
			})
		})

		When("Package has shared dependencies", func() {
			It("Produces four models", func() {
				// Successfully download the debendency
				apter := &internal.FakeApt{}
				apter.DownloadPackageReturnsOnCall(0, jq, 0, nil)
				apter.DownloadPackageReturnsOnCall(1, libjq1, 0, nil)
				apter.DownloadPackageReturnsOnCall(2, libonig5, 0, nil)
				apter.DownloadPackageReturnsOnCall(3, libc6, 0, nil)

				// No dependencies found via dpkg -I
				dpkger := &internal.FakeDpkg{}
				dpkger.IdentifyDependenciesReturnsOnCall(0, []string{"libjq1", "libc6"})
				dpkger.IdentifyDependenciesReturnsOnCall(1, []string{"libonig5", "libc6"})
				dpkger.IdentifyDependenciesReturnsOnCall(2, []string{"libc6"})
				dpkger.IdentifyDependenciesReturnsOnCall(3, []string{})

				dpkgQuery := &internal.FakeDpkgQuery{}

				config := &pkg.Config{
					ExcludeInstalledPackages: false,
				}

				packager := pkg.Analyser{
					Apt:    apter,
					Dpkg:   dpkger,
					Config: config,
					Query:  dpkgQuery,
				}

				modelMap := make(map[string]*pkg.PackageModel)
				modelList := make([]*pkg.PackageModel, 0)

				_ = packager.BuildPackage("", modelMap, &modelList)

				By("Invoking Apt and Dpkg twice", func() {
					Expect(apter.DownloadPackageCallCount()).To(Equal(4))
					Expect(dpkger.IdentifyDependenciesCallCount()).To(Equal(4))
				})
				By("Adding four models to the map", func() {
					Expect(len(modelMap)).To(Equal(4))
				})
				By("Returning jq model", func() {
					model := modelMap["jq"]
					Expect(model.Name).To(Equal("jq"))
					Expect(model.Version).To(Equal("1.6-1ubuntu0.20.04.1"))
					Expect(model.Filepath).To(Equal(jqFile))
				})
				By("Returning libjq1 model", func() {
					model := modelMap["libjq1"]
					Expect(model.Name).To(Equal("libjq1"))
					Expect(model.Version).To(Equal("1.6-1ubuntu0.20.04.1"))
					Expect(model.Filepath).To(Equal(libjq1File))
				})
				By("Returning libonig5 model", func() {
					model := modelMap["libonig5"]
					Expect(model.Name).To(Equal("libonig5"))
					Expect(model.Version).To(Equal("6.9.4-1"))
					Expect(model.Filepath).To(Equal(libonig5File))
				})

				By("Returning libc6 model", func() {
					model := modelMap["libc6"]
					Expect(model.Name).To(Equal("libc6"))
					//Expect(model.Version).To(Equal("2.31-0ubuntu9.15"))
					Expect(model.Filepath).To(Equal(libc6File))
				})
				By("Adding libc6 to the jq dependencies map", func() {
					model := modelMap["jq"]
					Expect(model.Dependencies["libc6"]).To(Not(BeNil()))
				})
				By("Adding libjq1 to the jq dependencies map", func() {
					model := modelMap["libjq1"]
					Expect(model.Dependencies["libc6"]).To(Not(BeNil()))
				})
				By("Adding libonig5 to the  dependencies map", func() {
					model := modelMap["libonig5"]
					Expect(model.Dependencies["libc6"]).To(Not(BeNil()))
				})
				By("Producing a puml diagram", func() {

					fmt.Println(puml.GenerateDiagram(&pkg.Config{}, modelMap, modelList))
				})
				By("DpkgQuery not being invoked", func() {
					Expect(dpkgQuery.IsInstalledCallCount()).To(Equal(0))
				})
			})
		})

		When("Package is libgsoap", func() {
			It("Produces a single model", func() {
				// Successfully download the debendency
				apter := &internal.FakeApt{}
				apter.DownloadPackageReturns(libgsoap, 0, nil)

				// No dependencies found via dpkg -I
				dpkger := &internal.FakeDpkg{}
				dpkger.IdentifyDependenciesReturns([]string{})

				dpkgQuery := &internal.FakeDpkgQuery{}

				config := &pkg.Config{
					ExcludeInstalledPackages: false,
				}

				packager := pkg.Analyser{
					Apt:    apter,
					Dpkg:   dpkger,
					Config: config,
					Query:  dpkgQuery,
				}

				modelMap := make(map[string]*pkg.PackageModel)
				modelList := make([]*pkg.PackageModel, 0)
				model := packager.BuildPackage("", modelMap, &modelList)

				By("Adding a model to the map", func() {
					Expect(len(modelMap)).To(Equal(1))
				})
				By("Returning a populated model", func() {
					Expect(model.Name).To(Equal("libgsoap-2.8.117"))
					Expect(model.Version).To(Equal("2.8.117-2build1"))
					Expect(model.Filepath).To(Equal(libgsoapFile))
				})
				By("Invoking Apt and Dpkg only once", func() {
					Expect(apter.DownloadPackageCallCount()).To(Equal(1))
					Expect(dpkger.IdentifyDependenciesCallCount()).To(Equal(1))
				})
				By("DpkgQuery not being invoked", func() {
					Expect(dpkgQuery.IsInstalledCallCount()).To(Equal(0))
				})
			})
		})
	})
})
