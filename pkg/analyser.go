package pkg

import (
	"com/alexander/debendency/pkg/commands"
	"fmt"
	"strings"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fake_dpkger.go com/alexander/debendency/pkg/commands.Dpkg
//counterfeiter:generate -o internal/fake_apter.go com/alexander/debendency/pkg/commands.Apt
//counterfeiter:generate -o internal/fake_query.go com/alexander/debendency/pkg/commands.DpkgQuery

type PackageModel struct {
	Filepath     string
	Name         string
	Version      string
	Dependencies map[string]*PackageModel
	IsInstalled  bool
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
	return packager.BuildPackage(name, modelMap)
}

func (packager Analyser) BuildPackage(name string, modelMap map[string]*PackageModel) *PackageModel {
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
	// Iterate through dependencies
	dependencyNames := packager.Dpkg.IdentifyDependencies(packageModel.Filepath)

	fmt.Printf("Dependencies %#v\n", dependencyNames)
	fmt.Printf("Dependencies length %d\n", len(dependencyNames))
	packageModel.Dependencies = make(map[string]*PackageModel, len(dependencyNames))

	for _, name := range dependencyNames {
		fmt.Printf("Building dependency [%s] for %s \n", name, packageModel.Name)
		dep := packager.BuildPackage(name, modelMap)
		packageModel.Dependencies[dep.Name] = dep
	}
	return &packageModel
}

func (packageModel *PackageModel) GetPackageFilename(name string) {
	fmt.Printf("Package download output: %#v\n", name)
	outputArray := strings.Split(name, "\n")
	fmt.Printf("Number of lines: %d\n", len(outputArray))
	for index := range outputArray {
		fmt.Println(outputArray[index])
	}

	downloadOutputLine := strings.Split(outputArray[0], " ")
	// Length should be 8
	// Get:1 http://gb.archive.ubuntu.com/ubuntu focal/universe amd64 dos2unix amd64 7.4.0-2 [374 kB]
	// Get:1 https://repo.saltproject.io/py3/ubuntu/20.04/amd64/3004 focal/main amd64 salt-master all 3004.2+ds-1 [40.9 kB]
	packageName := downloadOutputLine[4]
	fmt.Printf("PackageName: %s\n", packageName)
	arch := downloadOutputLine[5]
	fmt.Printf("Arch: %s\n", arch)
	version := downloadOutputLine[6]
	fmt.Printf("Version: %s\n", version)

	fileName := packageName + "_" + version + "_" + arch + ".deb"
	fmt.Printf("Filename: %s\n", fileName)
	//Check file exists
	packageModel.Filepath = fileName
	packageModel.Name = packageName
	packageModel.Version = version
}
