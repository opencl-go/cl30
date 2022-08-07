package cl30

// #include "api.h"
import "C"
import "unsafe"

// Bool represents a boolean value in the OpenCL API.
// It is not guaranteed to be the same size as the bool in kernels.
type Bool C.cl_bool

const (
	// False is the Bool value representing "false".
	False Bool = C.CL_FALSE
	// True is the Bool value representing "true".
	True Bool = C.CL_TRUE
)

// BoolFrom returns the Bool equivalent of a boolean value.
func BoolFrom(b bool) Bool {
	if b {
		return True
	}
	return False
}

// ToGoBool returns false if the Bool value is False, and true otherwise.
func (b Bool) ToGoBool() bool {
	return b != False
}

// Uint represents an unsigned 32-bit integer in the OpenCL API.
type Uint C.cl_uint

// Ulong represents an unsigned 64-bit integer in the OpenCL API.
type Ulong C.cl_ulong

const (
	// NameVersionByteSize is the size, in bytes, of the NameVersion structure.
	NameVersionByteSize = unsafe.Sizeof(C.cl_name_version{})
	// NameVersionMaxNameSize is the maximum number of bytes the name component of NameVersion can have.
	// This value includes the terminating NUL character, so the effective maximum length the string can have is
	// one byte less.
	NameVersionMaxNameSize = C.CL_NAME_VERSION_MAX_NAME_SIZE
)

// NameVersionName is a convenience type for the NameVersion.Name field.
type NameVersionName [NameVersionMaxNameSize]byte

// String returns the name value as a proper string, with the terminating NUL character removed.
func (name NameVersionName) String() string {
	name[NameVersionMaxNameSize-1] = 0x00
	return C.GoString((*C.char)(unsafe.Pointer(&name[0])))
}

// NameVersion is a combination of a name with a version.
type NameVersion struct {
	// Version describes the maturity level of the identified element.
	Version Version
	// Name identifies the element.
	Name NameVersionName
}
