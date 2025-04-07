//go:build !ethver

package version

const (
	Major = 0      // Major version component of the current release
	Minor = 3      // Minor version component of the current release
	Patch = 1      // Patch version component of the current release
	Meta  = "goat" // Version metadata to append to the version string
)
