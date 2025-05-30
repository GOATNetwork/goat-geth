//go:build !ethver

package version

const (
	Major = 0    // Major version component of the current release
	Minor = 3    // Minor version component of the current release
	Patch = 2    // Patch version component of the current release
	Meta  = "rc" // Version metadata to append to the version string
)
