package pkg

import (
	"com/alexander/debendency/pkg/commands"
	"fmt"
	"log/slog"
	"strings"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fake_dpkger.go com/alexander/debendency/pkg/commands.Dpkg
//counterfeiter:generate -o internal/fake_apter.go com/alexander/debendency/pkg/commands.Apt
//counterfeiter:generate -o internal/fake_query.go com/alexander/debendency/pkg/commands.DpkgQuery

type PackageModel struct {
	Filepath            string
	Name                string
	Version             string
	Dependencies        map[string]*PackageModel
	OrderedDependencies []*PackageModel
	IsInstalled         bool
}

type Analyser struct {
	Apt    commands.Apt
	Dpkg   commands.Dpkg
	Config *Config
	Query  commands.DpkgQuery
}

func NewAnalyser(config *Config) Analyser {
	command := commands.LinuxCommand{}
	apter := commands.Apter{
		Cmd: command,
	}
	dpkger := commands.Dpkger{
		Cmd: command,
	}
	query := commands.Query{
		Cmd: command,
	}
	return Analyser{
		Apt:    apter,
		Dpkg:   dpkger,
		Config: config,
		Query:  query,
	}
}

func (packager Analyser) Start(name string) *PackageModel {
	modelMap := make(map[string]*PackageModel)
	modelList := make([]*PackageModel, 0)
	return packager.BuildPackage(name, modelMap, &modelList)
}

func (packager Analyser) BuildPackage(name string, modelMap map[string]*PackageModel, modelList *[]*PackageModel) *PackageModel {
	model, ok := modelMap[name]
	if ok {
		return model
	}
	// start with 1st package
	standardOut, _, err := packager.Apt.DownloadPackage(name)

	if err != nil {
		return &PackageModel{}
	}
	packageModel := PackageModel{}
	packageModel.GetPackageFilename(standardOut)

	if packager.Config.ExcludeInstalledPackages {
		packageModel.IsInstalled = packager.Query.IsInstalled(packageModel.Name)
	}
	// Add packageModel to map under package name
	// packageModel.Name
	modelMap[packageModel.Name] = &packageModel
	// Add packageModel to model List in order of addition (as a map won't record this)
	*modelList = append(*modelList, &packageModel)
	// Iterate through dependencies
	dependencyNames := packager.Dpkg.IdentifyDependencies(packageModel.Filepath)

	slog.Debug(fmt.Sprintf("Dependencies %#v\n", dependencyNames))
	slog.Debug(fmt.Sprintf("Dependencies length %d\n", len(dependencyNames)))
	packageModel.Dependencies = make(map[string]*PackageModel, len(dependencyNames))
	packageModel.OrderedDependencies = make([]*PackageModel, 0)

	for _, name := range dependencyNames {
		slog.Debug(fmt.Sprintf("Building dependency [%s] for %s \n", name, packageModel.Name))
		dep := packager.BuildPackage(name, modelMap, modelList)
		packageModel.Dependencies[dep.Name] = dep

		packageModel.OrderedDependencies = append(packageModel.OrderedDependencies, dep)
	}
	return &packageModel
}

func (packageModel *PackageModel) GetPackageFilename(name string) {
	slog.Debug(fmt.Sprintf("Package download output: %#v\n", name))

	//fetchLine := strings.Split(name, "Get:1")[1]
	// We've now found the fetch line, break it down
	if strings.Contains(name, "Get:1") {
		outputArray := strings.Split(name, "Get:1")
		slog.Debug(fmt.Sprintf("Number of lines: %d\n", len(outputArray)))
		for index := range outputArray {
			slog.Debug(outputArray[index])
		}

		downloadOutputLine := strings.Split(outputArray[1], " ")
		// Length should be 8
		// Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]
		// Get:1 https://repo.saltproject.io/py3/ubuntu/20.04/amd64/3004 focal/main amd64 salt-master all 3004.2+ds-1 [40.9 kB]
		packageName := downloadOutputLine[4]
		slog.Debug(fmt.Sprintf("PackageName: %s\n", packageName))
		arch := downloadOutputLine[5]
		slog.Debug(fmt.Sprintf("Arch: %s\n", arch))
		version := downloadOutputLine[6]
		slog.Debug(fmt.Sprintf("Version: %s\n", version))

		fileName := packageName + "_" + version + "_" + arch + ".deb"
		fileName = strings.ReplaceAll(fileName, ":", "%3a")
		slog.Debug(fmt.Sprintf("Filename: %s\n", fileName))
		//Check file exists
		packageModel.Filepath = fileName
		packageModel.Name = packageName
		packageModel.Version = version
	}
}
