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

func (cache Cache) ClearBefore() {
	os.UserCacheDir()
	workingDir, err := os.Getwd()
	if err != nil {
		slog.Error("Pre-Cache issue", fmt.Errorf("Unable locate current working directory for cache clearance\n"))
	}
	dirFiles, err := os.ReadDir(workingDir)
	if err != nil {
		slog.Error("Pre-Cache issue", fmt.Errorf("Unable to clear cache of .deb files in current directory\n"))
	}
	for _, dirEntry := range dirFiles {
		if strings.HasSuffix(dirEntry.Name(), ".deb") {
			slog.Debug(fmt.Sprintf("Found debian installer: %s, attempting to delete\n", dirEntry.Name()))
			err := os.Remove(dirEntry.Name())
			if err != nil {
				slog.Error("Pre-Cache issue", fmt.Errorf("Unable to remove debian installer file: %s\n", dirEntry.Name()))
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
	slog.Debug("Attempting to clear cache\n")
	for _, packageModel := range modelMap {
		slog.Debug(fmt.Sprintf("Removing debian installer file: %s\n", packageModel.Filepath))
		err := os.Remove(packageModel.Filepath)
		if err != nil {
			slog.Error("Post-Cache issue", fmt.Errorf("Unable to remove debian installer file: %s\n", packageModel.Filepath))
		}
	}
}
