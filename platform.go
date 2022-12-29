package cl30

// #cgo LDFLAGS: -lOpenCL
// #include "api.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// PlatformID references one of the available OpenCL platforms of the system.
// It allows applications to query OpenCL devices, device configuration information, and to create OpenCL contexts
// using one or more devices.
// Retrieve a list of available platforms with the function PlatformIDs().
type PlatformID uintptr

func (id PlatformID) handle() C.cl_platform_id {
	return *(*C.cl_platform_id)(unsafe.Pointer(&id))
}

// String provides a readable presentation of the platform identifier.
// It is based on the numerical value of the underlying pointer.
func (id PlatformID) String() string {
	return fmt.Sprintf("0x%X", uintptr(id))
}

// PlatformIDs returns the list of available platforms on the system.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetPlatformIDs.html
func PlatformIDs() ([]PlatformID, error) {
	requiredCount := C.cl_uint(0)
	status := C.clGetPlatformIDs(0, nil, &requiredCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	if requiredCount == 0 {
		return nil, nil
	}
	ids := make([]PlatformID, requiredCount)
	returnedCount := C.cl_uint(0)
	status = C.clGetPlatformIDs(requiredCount, (*C.cl_platform_id)(unsafe.Pointer(&ids[0])), &returnedCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	usedCount := returnedCount
	if usedCount > requiredCount {
		usedCount = requiredCount
	}
	return ids[:usedCount], nil
}

// PlatformInfoName identifies properties of a platform, which can be queried with PlatformInfo().
type PlatformInfoName C.cl_platform_info

const (
	// PlatformNameInfo refers to a human-readable string that identifies the platform.
	//
	// Returned type: string
	PlatformNameInfo PlatformInfoName = C.CL_PLATFORM_NAME
	// PlatformVendorInfo refers to a human-readable string that identifies the vendor of the platform.
	//
	// Returned type: string
	PlatformVendorInfo PlatformInfoName = C.CL_PLATFORM_VENDOR
	// PlatformProfileInfo refers to the profile name supported by the implementation.
	// The profile name returned can be one of the following strings:
	//
	// "FULL_PROFILE" - if the implementation supports the OpenCL specification (functionality defined as part of the
	// core specification and does not require any extensions to be supported).
	//
	// "EMBEDDED_PROFILE" - if the implementation supports the OpenCL embedded profile. The embedded profile is defined
	// to be a subset for each version of OpenCL.
	//
	// Returned type: string
	PlatformProfileInfo PlatformInfoName = C.CL_PLATFORM_PROFILE
	// PlatformVersionInfo refers to the OpenCL version supported by the implementation.
	// This version string has the following format:
	//
	// OpenCL<space><major_version.minor_version><space><platform-specific information>
	//
	// Returned type: string
	PlatformVersionInfo PlatformInfoName = C.CL_PLATFORM_VERSION
	// PlatformNumericVersionInfo refers to the detailed (major, minor, patch) version supported by the platform.
	// The major and minor version numbers returned must match those returned via PlatformVersionInfo.
	//
	// Returned type: Version
	// Since: 3.0
	PlatformNumericVersionInfo PlatformInfoName = C.CL_PLATFORM_NUMERIC_VERSION
	// PlatformExtensionsInfo refers to a space separated list of extension names (the extension names themselves do not
	// contain any spaces) supported by the platform. Each extension that is supported by all devices associated with
	// this platform must be reported here.
	//
	// Returned type: string
	PlatformExtensionsInfo PlatformInfoName = C.CL_PLATFORM_EXTENSIONS
	// PlatformExtensionsWithVersionInfo refers to an array of description (name and version) structures that lists all
	// the extensions supported by the platform. The same extension name must not be reported more than once.
	// The list of extensions reported must match the list reported via PlatformExtensionsInfo.
	//
	// Returned type: []NameVersion
	// Since: 3.0
	PlatformExtensionsWithVersionInfo PlatformInfoName = C.CL_PLATFORM_EXTENSIONS_WITH_VERSION
	// PlatformHostTimerResolutionInfo refers to the resolution of the host timer in nanoseconds as used by
	// DeviceAndHostTimer() and HostTimer().
	// This value must be 0 for devices that do not support device and host timer synchronization.
	//
	// Returned type: uint64
	// Since: 2.1
	PlatformHostTimerResolutionInfo PlatformInfoName = C.CL_PLATFORM_HOST_TIMER_RESOLUTION
)

// PlatformInfo queries information about an OpenCL platform.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use PlatformInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetPlatformInfo.html
func PlatformInfo(id PlatformID, paramName PlatformInfoName, param HostMemory) (uintptr, error) {
	sizeReturn := C.size_t(0)
	paramPtr := ResolvePointer(param, false, "param")
	status := C.clGetPlatformInfo(
		id.handle(),
		C.cl_platform_info(paramName),
		sizeOf(param),
		paramPtr,
		&sizeReturn)
	runtime.KeepAlive(paramPtr)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// PlatformInfoString is a convenience method for PlatformInfo() to query information values that are string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func PlatformInfoString(id PlatformID, paramName PlatformInfoName) (string, error) {
	return queryString(func(param HostMemory) (uintptr, error) {
		return PlatformInfo(id, paramName, param)
	})
}

// ExtensionFunctionAddressForPlatform returns the address of the extension function named by functionName
// for a given platform.
//
// The pointer returned should be cast to a C-function pointer type matching the extension function's definition
// defined in the appropriate extension specification and header file.
//
// A return value of nil indicates that the specified function does not exist for the implementation or
// platform is not a valid platform.
// A non-nil return value for ExtensionFunctionAddressForPlatform() does not guarantee that an extension function
// is actually supported by the platform. The application must also make a corresponding query using
// PlatformInfo(platform, PlatformExtensionsInfo, ...) or DeviceInfo(device, DeviceExtensionsInfo, ...) to determine
// if an extension is supported by the OpenCL implementation.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetExtensionFunctionAddressForPlatform.html
func ExtensionFunctionAddressForPlatform(id PlatformID, functionName string) unsafe.Pointer {
	rawName := C.CString(functionName)
	defer C.free(unsafe.Pointer(rawName))
	return C.clGetExtensionFunctionAddressForPlatform(id.handle(), rawName)
}

// UnloadPlatformCompiler allows the implementation to release the resources allocated by the OpenCL compiler for
// a platform.
//
// This function allows the implementation to release the resources allocated by the OpenCL compiler for platform.
// This is a hint from the application and does not guarantee that the compiler will not be used in the future or
// that the compiler will actually be unloaded by the implementation.
// Calls to BuildProgram(), CompileProgram(), or LinkProgram() after UnloadPlatformCompiler() will reload the compiler,
// if necessary, to build the appropriate program executable.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clUnloadPlatformCompiler.html
func UnloadPlatformCompiler(id PlatformID) error {
	status := C.clUnloadPlatformCompiler(id.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}
