package commands

// Interfaces may allow for faking when testing our native commands
type Apt interface {
	DownloadPackage(name string) (string, int, error)
}

type Apter struct {
	Cmd Command
}

func (apter Apter) DownloadPackage(name string) (string, int, error) {
	return apter.Cmd.Command("apt", "download", name)
}
