package pkg_test

import (
	"com/alexander/debendency/pkg"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Cache", func() {

	Describe("ClearBefore", func() {

		var cacheDirectory string
		var originalDirectory string
		cacheLocation := "/internal/temp/"

		BeforeEach(func() {
			workingDir, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			temp, err := os.MkdirTemp(workingDir+cacheLocation, "")
			Expect(err).ToNot(HaveOccurred())

			cacheDirectory = temp
			originalDirectory = workingDir
			fmt.Println("New temp dir :" + temp)
			fmt.Println("Original dir :" + originalDirectory)

		})

		AfterEach(func() {
			fmt.Println("Switching to original dir :" + originalDirectory)
			err := os.Chdir(originalDirectory)
			Expect(err).ToNot(HaveOccurred())

			fmt.Println("Deleting temp dir :" + cacheDirectory)
			err = os.RemoveAll(cacheDirectory)
			Expect(err).ToNot(HaveOccurred())
		})

		When("Cache directory doesn't exist", func() {
			It("should succeed", func() {
				currentDir, err := os.Getwd()
				Expect(err).ToNot(HaveOccurred())
				cache := pkg.NewCache(pkg.Config{
					InstallerLocation: cacheDirectory,
				})

				workingDir, err := cache.ClearBefore()
				Expect(err).ToNot(HaveOccurred())

				By("Should create the cache directory", func() {
					dirEntries, err := os.ReadDir(cacheDirectory + string(os.PathSeparator) + pkg.DEBENDENCY)
					Expect(err).ToNot(HaveOccurred())
					Expect(dirEntries).To(BeEmpty())
				})
				By("Returning the original working directory", func() {
					Expect(workingDir).To(Equal(currentDir))
				})
			})

		})

		When("Cache has files in it", func() {
			// Create debendency directory
			It("Should clear the cache", func() {
				currentDir, err := os.Getwd()
				Expect(err).ToNot(HaveOccurred())
				cache := pkg.NewCache(pkg.Config{
					InstallerLocation: cacheDirectory,
				})

				err = os.Mkdir(cacheDirectory+string(os.PathSeparator)+pkg.DEBENDENCY, os.FileMode(755))
				Expect(err).ToNot(HaveOccurred())
				cacheRoot := cacheDirectory + string(os.PathSeparator) + pkg.DEBENDENCY + string(os.PathSeparator)
				// Create deb files
				err = os.WriteFile(cacheRoot+"test.deb", []byte{}, os.FileMode(755))
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile(cacheRoot+"test.txt", []byte{}, os.FileMode(755))
				Expect(err).ToNot(HaveOccurred())

				workingDir, err := cache.ClearBefore()
				Expect(err).ToNot(HaveOccurred())

				By("Should not recreate the cache directory", func() {
					// verify this still exists: cacheDirectory + string(os.PathSeparator)
					dirEntries, err := os.ReadDir(cacheDirectory + string(os.PathSeparator) + pkg.DEBENDENCY)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(dirEntries)).To(Equal(1))
				})
				By("Should change to the cache directory", func() {
					currentDir, err := os.Getwd()
					Expect(err).ToNot(HaveOccurred())
					fmt.Println("Have moved to the cache dir:" + currentDir)

					Expect(filepath.ToSlash(currentDir)).To(Equal(filepath.ToSlash(cacheDirectory + string(os.PathSeparator) + pkg.DEBENDENCY)))
				})
				By("Should clear all *.deb files", func() {
					currentDir, err := os.Getwd()
					dirEntries, err := os.ReadDir(currentDir)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(dirEntries)).To(Equal(1))
				})
				By("Returning the original working directory", func() {
					Expect(workingDir).To(Equal(currentDir))
				})
			})
		})
	})
})
