package cl30

import (
	"fmt"
	"math"
)

// Version represents a major.minor.patch version combination, encoded in a 32-bit unsigned integer value.
type Version uint32

const (
	// VersionMin is the lowest value that Version can represent. It has all components set to zero.
	VersionMin Version = 0
	// VersionMax is the highest value that Version can represent.
	// It has all components set to their maximally possible value.
	VersionMax Version = math.MaxUint32

	versionMajorBits = 10
	versionMinorBits = 10
	versionPatchBits = 12
	versionMajorMask = (1 << versionMajorBits) - 1
	versionMinorMask = (1 << versionMinorBits) - 1
	versionPatchMask = (1 << versionPatchBits) - 1
)

// VersionOf returns a Version that has the three provided components encoded.
// No particular limits-checking is performed, the provided values are cast and shifted into the final field.
func VersionOf(major, minor, patch int) Version {
	return Version(
		((uint32(major) & versionMajorMask) << (versionMinorBits + versionPatchBits)) |
			((uint32(minor) & versionMinorMask) << versionPatchBits) | (uint32(patch) & versionPatchMask))
}

// String returns a common representation of <major> <.> <minor> <.> <patch> of the version, with the components
// in decimal format.
func (ver Version) String() string {
	return fmt.Sprintf("%d.%d.%d", ver.Major(), ver.Minor(), ver.Patch())
}

// Major returns the major component value.
func (ver Version) Major() int {
	return int((uint32(ver) >> (versionMinorBits + versionPatchBits)) & versionMajorMask)
}

// Minor returns the minor component value.
func (ver Version) Minor() int {
	return int((uint32(ver) >> versionPatchBits) & versionMinorMask)
}

// Patch returns the patch component value.
func (ver Version) Patch() int {
	return int(uint32(ver) & versionPatchMask)
}
