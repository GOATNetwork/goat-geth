//go:build !ethver

package version

const (
	Major = 0          // Major version component of the current release
	Minor = 1          // Minor version component of the current release
	Patch = 9          // Patch version component of the current release
	Meta  = "testnet3" // Version metadata to append to the version string
)
