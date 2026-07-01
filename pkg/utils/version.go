package utils

import "fmt"

const (
	MajorVersion = 4
	MinorVersion = 2
	PatchVersion = 0
)

// Version returns the version string in format "major.minor.patch"
func Version() string {
	return fmt.Sprintf("%d.%d.%d",
		MajorVersion, MinorVersion, PatchVersion)
}
