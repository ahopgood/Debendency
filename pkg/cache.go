package pkg

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Cache struct {
	directory string
}

const DEBENDENCY = "debendency"

func NewCache(config Config) Cache {
	dir := config.InstallerLocation + "/" + DEBENDENCY
	return Cache{directory: dir}
}

func (cache Cache) ClearBefore() (string, error) {
	startingDir, err := os.Getwd()

	// check dir exists
	err = os.MkdirAll(cache.directory, 0755)
	if err != nil {
		return "", fmt.Errorf("Pre-Cache issue: unable to create cache directory: %s", cache.directory)
	}
	// change dir
	err = os.Chdir(cache.directory)
	if err != nil {
		return "", fmt.Errorf("Pre-Cache issue: unable to switch to cache directory: %s", cache.directory)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Pre-Cache issue: unable locate current working directory for cache clearance")
	}
	dirFiles, err := os.ReadDir(workingDir)
	if err != nil {
		return "", fmt.Errorf("Pre-Cache issue: unable to clear cache of .deb files in current directory")
	}
	for _, dirEntry := range dirFiles {
		if strings.HasSuffix(dirEntry.Name(), ".deb") {
			slog.Debug(fmt.Sprintf("Found debian installer: %s, attempting to delete\n", dirEntry.Name()))
			err := os.Remove(dirEntry.Name())
			if err != nil {
				slog.Error(fmt.Sprintf("Pre-Cache issue: Unable to remove debian installer file: %s", dirEntry.Name()))
			}
		}
	}
	return startingDir, nil
}
func (cache Cache) clearAfter(modelMap map[string]*PackageModel) {
	//os.Chdir()
	//os.MkdirTemp()
	//os.Mkdir()
	//os.MkdirAll()
	//os.Remove()
	slog.Debug("Attempting to clear cache\n")
	for _, packageModel := range modelMap {
		slog.Debug(fmt.Sprintf("Removing debian installer file: %s\n", packageModel.Filepath))
		err := os.Remove(packageModel.Filepath)
		if err != nil {
			slog.Error(fmt.Sprintf("Post-Cache issue: unable to remove debian installer file: %s\n", packageModel.Filepath))
		}
	}
}
