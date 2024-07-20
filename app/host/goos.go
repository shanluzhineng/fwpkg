package host

type GOOS_OS string

const (
	goos_windows GOOS_OS = "windows"
	goos_linux   GOOS_OS = "linux"
)

func (os GOOS_OS) IsWindows() bool {
	return os == goos_windows
}

func (os GOOS_OS) IsLinux() bool {
	return os == goos_linux
}
