package pkg

import (
	"fmt"
	"os"
	"strings"
)

type Cache struct {
	directory string
}

func (cache Cache) ClearBefore() {
	os.UserCacheDir()
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Unable locate current working directory for cache clearance\n")
	}
	dirFiles, err := os.ReadDir(workingDir)
	if err != nil {
		fmt.Errorf("Unable to clear cache of .deb files in current directory\n")
	}
	for _, dirEntry := range dirFiles {
		if strings.HasSuffix(dirEntry.Name(), ".deb") {
			fmt.Printf("Found debian installer: %s, attempting to delete\n", dirEntry.Name())
			err := os.Remove(dirEntry.Name())
			if err != nil {
				fmt.Errorf("Unable to remove debian installer file: %s\n", dirEntry.Name())
			}
		}
	}
}
func (cache Cache) clearAfter(modelMap map[string]*PackageModel) {
	//os.Chdir()
	//os.MkdirTemp()
	//os.Mkdir()
	//os.MkdirAll()
	//os.Remove()
	fmt.Printf("Attempting to clear cache\n")
	for _, packageModel := range modelMap {
		fmt.Printf("Removing debiand installer file: %s\n", packageModel.Filepath)
		err := os.Remove(packageModel.Filepath)
		if err != nil {
			fmt.Errorf("Unable to remove debian installer file: %s\n", packageModel.Filepath)
		}
	}
}
