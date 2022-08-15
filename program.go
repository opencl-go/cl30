package cl30

// #include "api.h"
// extern cl_int cl30BuildProgram(cl_program program,
//    cl_uint numDevices, cl_device_id *devices,
//    char *options, uintptr_t *userData);
// extern cl_int cl30CompileProgram(cl_program program,
//    cl_uint numDevices, cl_device_id *devices,
//    char *options,
//    cl_uint numInputHeaders, cl_program *headers, char **includeNames,
//    uintptr_t *userData);
// extern cl_program cl30LinkProgram(cl_context context,
//    cl_uint numDevices, cl_device_id *devices,
//    char *options,
//    cl_uint numInputPrograms, cl_program *programs,
//    uintptr_t *userData,
//    cl_int *errReturn);
import "C"
import (
	"fmt"
	"unsafe"
)

// Program objects contain executable code for the OpenCL runtime.
type Program uintptr

func (program Program) handle() C.cl_program {
	return *(*C.cl_program)(unsafe.Pointer(&program))
}

// String provides a readable presentation of the program identifier.
// It is based on the numerical value of the underlying pointer.
func (program Program) String() string {
	return fmt.Sprintf("0x%X", uintptr(program))
}

// CreateProgramWithSource creates a program object for a context, and loads source code specified by text strings
// into the program object.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateProgramWithSource.html
func CreateProgramWithSource(context Context, sources []string) (Program, error) {
	rawSources := make([]*C.char, len(sources))
	for i := 0; i < len(sources); i++ {
		rawSources[i] = C.CString(sources[i])
	}
	defer func() {
		for _, rawSource := range rawSources {
			C.free(unsafe.Pointer(rawSource))
		}
	}()
	var status C.cl_int
	program := C.clCreateProgramWithSource(
		context.handle(),
		C.cl_uint(len(rawSources)),
		(**C.char)(unsafe.Pointer(&rawSources[0])),
		nil,
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Program(*((*uintptr)(unsafe.Pointer(&program)))), nil
}

// CreateProgramWithIl creates a program object for a context, and loads the intermediate language (IL) into the
// program object.
//
// The intermediate language pointed to by il will be loaded into the program object. The devices associated with
// the program object are the devices associated with context.
//
// Since: 2.1
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateProgramWithIL.html
func CreateProgramWithIl(context Context, il []byte) (Program, error) {
	var rawIl unsafe.Pointer
	if len(il) > 0 {
		rawIl = unsafe.Pointer(&il[0])
	}
	var status C.cl_int
	program := C.clCreateProgramWithIL(
		context.handle(),
		rawIl,
		C.size_t(len(il)),
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Program(*((*uintptr)(unsafe.Pointer(&program)))), nil
}

// CreateProgramWithBinary creates a program object for a context, and loads binary bits into the program object.
//
// The returned slice of errors represents the load-status per device.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateProgramWithBinary.html
func CreateProgramWithBinary(context Context, devices []DeviceID, binaries [][]byte) (Program, []error, error) {
	rawBinaries := make([]*C.uchar, len(binaries))
	binaryLengths := make([]C.size_t, len(binaries))
	for i := 0; i < len(binaries); i++ {
		rawBinaries[i] = (*C.uchar)(unsafe.Pointer(&binaries[i][0]))
		binaryLengths[i] = C.size_t(len(binaries[i]))
	}
	binaryStatus := make([]C.cl_int, len(devices))
	var status C.cl_int
	program := C.clCreateProgramWithBinary(
		context.handle(),
		C.cl_uint(len(devices)),
		(*C.cl_device_id)(unsafe.Pointer(&devices[0])),
		(*C.size_t)(unsafe.Pointer(&binaryLengths[0])),
		(**C.uchar)(unsafe.Pointer(&rawBinaries[0])),
		(*C.cl_int)(unsafe.Pointer(&binaryStatus[0])),
		&status)
	binaryErr := make([]error, len(devices))
	for i := 0; i < len(devices); i++ {
		if binaryStatus[i] != C.CL_SUCCESS {
			binaryErr[i] = StatusError(binaryStatus[i])
		}
	}
	if status != C.CL_SUCCESS {
		return 0, binaryErr, StatusError(status)
	}
	return Program(*((*uintptr)(unsafe.Pointer(&program)))), binaryErr, nil
}

// CreateProgramWithBuiltInKernels creates a program object for a context, and loads the information related to the
// built-in kernels into a program object.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateProgramWithBuiltInKernels.html
func CreateProgramWithBuiltInKernels(context Context, devices []DeviceID, kernelNames string) (Program, error) {
	rawKernelNames := C.CString(kernelNames)
	defer C.free(unsafe.Pointer(rawKernelNames))
	var status C.cl_int
	program := C.clCreateProgramWithBuiltInKernels(
		context.handle(),
		C.cl_uint(len(devices)),
		(*C.cl_device_id)(unsafe.Pointer(&devices[0])),
		rawKernelNames,
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Program(*((*uintptr)(unsafe.Pointer(&program)))), nil
}

// RetainProgram increments the program reference count.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clRetainProgram.html
func RetainProgram(program Program) error {
	status := C.clRetainProgram(program.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ReleaseProgram decrements the program reference count.
//
// The program object is deleted after all kernel objects associated with program have been deleted and
// the program reference count becomes zero.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clReleaseProgram.html
func ReleaseProgram(program Program) error {
	status := C.clReleaseProgram(program.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// BuildProgram builds (compiles and links) a program executable from the program source or binary.
//
// The notification routine is a callback function that an application can register and which will be called when
// the program executable has been built (successfully or unsuccessfully). If callback is not nil, BuildProgram()
// does not need to wait for the build to complete and can return immediately once the build operation can begin.
// Any state changes of the program object that result from calling BuildProgram() (e.g. build status or log) will
// be observable from this callback function. The build operation can begin if the context, program whose sources
// are being compiled and linked, list of devices and build options specified are all valid and appropriate host
// and device resources needed to perform the build are available.
// If callback is nil, BuildProgram() does not return until the build has completed.
// This callback function may be called asynchronously by the OpenCL implementation. It is the applications
// responsibility to ensure that the callback function is thread-safe.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clBuildProgram.html
func BuildProgram(program Program, devices []DeviceID, options string, callback func()) error {
	rawOptions := C.CString(options)
	defer C.free(unsafe.Pointer(rawOptions))
	var rawDevices unsafe.Pointer
	if len(devices) > 0 {
		rawDevices = unsafe.Pointer(&devices[0])
	}
	var callbackUserData userData
	if callback != nil {
		var err error
		callbackUserData, err = userDataFor(callback)
		if err != nil {
			return err
		}
	}
	status := C.cl30BuildProgram(
		program.handle(),
		C.cl_uint(len(devices)),
		(*C.cl_device_id)(rawDevices),
		rawOptions,
		callbackUserData.ptr)
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return StatusError(status)
	}
	return nil
}

//export cl30GoProgramBuildCallback
func cl30GoProgramBuildCallback(_ Program, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func())
	callbackUserData.Delete()
	callback()
}

// SetProgramSpecializationConstant sets a constant for a program created from intermediate language.
//
// The specialization value will be used by subsequent calls to BuildProgram() until another call to
// SetProgramSpecializationConstant() changes it.
//
// Since: 2.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetProgramSpecializationConstant.html
func SetProgramSpecializationConstant(program Program, id uint32, size uintptr, value unsafe.Pointer) error {
	status := C.clSetProgramSpecializationConstant(
		program.handle(),
		C.cl_uint(id),
		C.size_t(size),
		value)
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// IncludeHeader is a named program to be used with CompileProgram().
type IncludeHeader struct {
	Name    string
	Program Program
}

// CompileProgram compiles a program's source for all the devices or a specific device(s) in the OpenCL context
// associated with a program.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCompileProgram.html
func CompileProgram(program Program, devices []DeviceID, options string, headers []IncludeHeader, callback func()) error {
	rawOptions := C.CString(options)
	defer C.free(unsafe.Pointer(rawOptions))
	var rawDevices unsafe.Pointer
	if len(devices) > 0 {
		rawDevices = unsafe.Pointer(&devices[0])
	}
	var callbackUserData userData
	if callback != nil {
		var err error
		callbackUserData, err = userDataFor(callback)
		if err != nil {
			return err
		}
	}
	var rawHeaderProgramsPtr unsafe.Pointer
	var rawHeaderNamesPtr unsafe.Pointer
	rawHeaderPrograms := make([]Program, len(headers))
	rawHeaderNames := make([]*C.char, len(headers))
	for i := 0; i < len(headers); i++ {
		rawHeaderPrograms[i] = headers[i].Program
		rawHeaderNames[i] = C.CString(headers[i].Name)
	}
	defer func() {
		for _, rawHeaderName := range rawHeaderNames {
			C.free(unsafe.Pointer(rawHeaderName))
		}
	}()
	if len(headers) > 0 {
		rawHeaderProgramsPtr = unsafe.Pointer(&rawHeaderPrograms[0])
		rawHeaderNamesPtr = unsafe.Pointer(&rawHeaderNames[0])
	}
	status := C.cl30CompileProgram(
		program.handle(),
		C.cl_uint(len(devices)),
		(*C.cl_device_id)(rawDevices),
		rawOptions,
		C.cl_uint(len(headers)),
		(*C.cl_program)(rawHeaderProgramsPtr),
		(**C.char)(rawHeaderNamesPtr),
		callbackUserData.ptr)
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return StatusError(status)
	}
	return nil
}

//export cl30GoProgramCompileCallback
func cl30GoProgramCompileCallback(_ Program, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func())
	callbackUserData.Delete()
	callback()
}

// LinkProgram links a set of compiled program objects and libraries for all the devices or a specific device(s)
// in the OpenCL context and creates a library or executable.
//
// The notification routine is a callback function that an application can register and which will be called when
// the program executable has been built (successfully or unsuccessfully).
// If callback is not nil, LinkProgram() does not have to wait until the linker to complete and can return
// if the linking operation can begin.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clLinkProgram.html
func LinkProgram(context Context, devices []DeviceID, options string, programs []Program, callback func(Program)) (Program, error) {
	rawOptions := C.CString(options)
	defer C.free(unsafe.Pointer(rawOptions))
	var rawDevices unsafe.Pointer
	if len(devices) > 0 {
		rawDevices = unsafe.Pointer(&devices[0])
	}
	var callbackUserData userData
	if callback != nil {
		var err error
		callbackUserData, err = userDataFor(callback)
		if err != nil {
			return 0, err
		}
	}
	var status C.cl_int
	program := C.cl30LinkProgram(
		context.handle(),
		C.cl_uint(len(devices)),
		(*C.cl_device_id)(rawDevices),
		rawOptions,
		C.cl_uint(len(programs)),
		(*C.cl_program)(unsafe.Pointer(&programs[0])),
		callbackUserData.ptr,
		&status)
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return 0, StatusError(status)
	}
	return Program(*((*uintptr)(unsafe.Pointer(&program)))), nil
}

//export cl30GoProgramLinkCallback
func cl30GoProgramLinkCallback(program Program, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func(Program))
	callbackUserData.Delete()
	callback(program)
}

// ProgramBuildInfoName identifies properties of a program build, which can be queried with ProgramBuildInfo().
type ProgramBuildInfoName C.cl_program_build_info

const (
	// ProgramBuildStatusInfo returns the build, compile, or link status, whichever was performed last on the
	// specified program object for device.
	//
	// Returned type: BuildStatus
	ProgramBuildStatusInfo ProgramBuildInfoName = C.CL_PROGRAM_BUILD_STATUS
	// ProgramBuildOptionsInfo return the build, compile, or link options specified by the options argument in
	// BuildProgram(), CompileProgram(), or LinkProgram(), whichever was performed last on the specified program
	// object for device.
	//
	// If build status of the specified program for device is BuildNoneStatus, an empty string is returned.
	//
	// Returned type: string
	ProgramBuildOptionsInfo ProgramBuildInfoName = C.CL_PROGRAM_BUILD_OPTIONS
	// ProgramBuildLogInfo returns the build, compile, or link log for BuildProgram(), CompileProgram(), or
	// LinkProgram(), whichever was performed last on program for device.
	//
	// If build status of the specified program for device is BuildNoneStatus, an empty string is returned.
	//
	// Returned type: string
	ProgramBuildLogInfo ProgramBuildInfoName = C.CL_PROGRAM_BUILD_LOG
	// ProgramBinaryTypeInfo returns the program binary type for device.
	//
	// Returned type: ProgramBinaryType
	// Since: 1.2
	ProgramBinaryTypeInfo ProgramBuildInfoName = C.CL_PROGRAM_BINARY_TYPE
	// ProgramBuildGlobalVariableTotalSizeInfo returns the total amount of storage, in bytes, used by program
	// variables in the global address space.
	//
	// Returned type: uintptr
	// Since: 2.0
	ProgramBuildGlobalVariableTotalSizeInfo ProgramBuildInfoName = C.CL_PROGRAM_BUILD_GLOBAL_VARIABLE_TOTAL_SIZE
)

// BuildStatus describes the build, compile, or link status of a program.
type BuildStatus C.cl_build_status

const (
	// BuildNoneStatus is the build status returned if no BuildProgram(), CompileProgram(), or LinkProgram() has been
	// performed on the specified program object for device.
	BuildNoneStatus BuildStatus = C.CL_BUILD_NONE
	// BuildSuccessStatus is the build status returned if BuildProgram(), CompileProgram(), or LinkProgram() -
	// whichever was performed last on the specified program object for device - was successful.
	BuildSuccessStatus BuildStatus = C.CL_BUILD_SUCCESS
	// BuildErrorStatus is the build status returned if BuildProgram(), CompileProgram(), or LinkProgram() -
	// whichever was performed last on the specified program object for device - generated an error.
	BuildErrorStatus BuildStatus = C.CL_BUILD_ERROR
	// BuildInProgressStatus is the build status returned if BuildProgram(), CompileProgram(), or LinkProgram() -
	// whichever was performed last on the specified program object for device - has not finished.
	BuildInProgressStatus BuildStatus = C.CL_BUILD_IN_PROGRESS
)

// ProgramBinaryType identifies the program binary type for devices.
type ProgramBinaryType C.cl_program_binary_type

const (
	// ProgramBinaryTypeNone is set if there is no binary associated with the specified program object for device.
	ProgramBinaryTypeNone ProgramBinaryType = C.CL_PROGRAM_BINARY_TYPE_NONE
	// ProgramBinaryTypeCompiledObject is a compiled binary associated with device.
	// This is the case when the specified program object was created using CreateProgramWithSource() and compiled
	// using CompileProgram(), or when a compiled binary was loaded using CreateProgramWithBinary().
	ProgramBinaryTypeCompiledObject ProgramBinaryType = C.CL_PROGRAM_BINARY_TYPE_COMPILED_OBJECT
	// ProgramBinaryTypeLibrary is a library binary associated with device.
	// This is the case when the specified program object was linked by LinkProgram() using the -create-library link
	// option, or when a compiled library binary was loaded using CreateProgramWithBinary().
	ProgramBinaryTypeLibrary ProgramBinaryType = C.CL_PROGRAM_BINARY_TYPE_LIBRARY
	// ProgramBinaryTypeExecutable is an executable binary associated with device.
	// This is the case when the specified program object was linked by LinkProgram() without the -create-library link
	// option, or when an executable binary was built using BuildProgram().
	ProgramBinaryTypeExecutable ProgramBinaryType = C.CL_PROGRAM_BINARY_TYPE_EXECUTABLE
)

// ProgramBuildInfo returns build information for each device in the program object.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use ProgramBuildInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetProgramBuildInfo.html
func ProgramBuildInfo(program Program, device DeviceID, paramName ProgramBuildInfoName, paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetProgramBuildInfo(
		program.handle(),
		device.handle(),
		C.cl_program_build_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// ProgramBuildInfoString is a convenience method for ProgramBuildInfo() to query information values that are
// string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func ProgramBuildInfoString(program Program, device DeviceID, paramName ProgramBuildInfoName) (string, error) {
	return queryString(func(paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
		return ProgramBuildInfo(program, device, paramName, paramSize, paramValue)
	})
}

// ProgramInfoName identifies properties of a program, which can be queried with ProgramInfo().
type ProgramInfoName C.cl_program_info

const (
	// ProgramReferenceCountInfo returns the program reference count.
	//
	// Note: The reference count returned should be considered immediately stale. It is unsuitable for
	// general use in applications. This feature is provided for identifying memory leaks.
	//
	// Returned type: Uint
	ProgramReferenceCountInfo ProgramInfoName = C.CL_PROGRAM_REFERENCE_COUNT
	// ProgramContextInfo returns the context specified when the program object is created.
	//
	// Returned type: Context
	ProgramContextInfo ProgramInfoName = C.CL_PROGRAM_CONTEXT
	// ProgramNumDevicesInfo returns the number of devices associated with program.
	//
	// Returned type: Uint
	ProgramNumDevicesInfo ProgramInfoName = C.CL_PROGRAM_NUM_DEVICES
	// ProgramDevicesInfo returns the list of devices associated with the program object. This can be the
	// devices associated with context on which the program object has been created or can be a subset of devices
	// that are specified when a program object is created using CreateProgramWithBinary().
	//
	// Returned type: []DeviceID
	ProgramDevicesInfo ProgramInfoName = C.CL_PROGRAM_DEVICES
	// ProgramSourceInfo returns the program source code specified by CreateProgramWithSource().
	// The source string returned is a concatenation of all source strings specified to CreateProgramWithSource().
	//
	// If program is created using CreateProgramWithBinary(), CreateProgramWithIl(), or
	// CreateProgramWithBuiltInKernels(), an empty string or the appropriate program source code is returned depending
	// on whether the program source code is stored in the binary.
	//
	// Returned type: string
	ProgramSourceInfo ProgramInfoName = C.CL_PROGRAM_SOURCE
	// ProgramBinarySizesInfo returns an array that contains the size in bytes of the program binary
	// (could be an executable binary, compiled binary, or library binary) for each device associated with program.
	// The size of the array is the number of devices associated with program. If a binary is not available for a
	// device(s), a size of zero is returned.
	//
	// If program is created using CreateProgramWithBuiltInKernels(), the implementation may return zero in any
	// entries of the returned array.
	//
	// Returned type: []uintptr
	ProgramBinarySizesInfo ProgramInfoName = C.CL_PROGRAM_BINARY_SIZES
	// ProgramBinariesInfo returns the program binaries (could be an executable binary, compiled binary, or
	// library binary) for all devices associated with program. For each device in program, the binary returned can
	// be the binary specified for the device when program is created with CreateProgramWithBinary() or it can be the
	// executable binary generated by BuildProgram() or LinkProgram().
	// If program is created with CreateProgramWithSource() or CreateProgramWithIl(), the binary returned is the
	// binary generated by BuildProgram(), CompileProgram(), or LinkProgram(). The bits returned can be an
	// implementation-specific intermediate representation (a.k.a. IR) or device specific executable bits or both.
	// The decision on which information is returned in the binary is up to the OpenCL implementation.
	//
	// The paramValue points to an array of N pointers allocated by the caller, where N is the number of devices
	// associated with program. The buffer sizes needed to allocate the memory that these N pointers refer to can
	// be queried using the ProgramBinarySizesInfo query.
	//
	// Each entry in this array is used by the implementation as the location in memory where to copy the program
	// binary for a specific device, if there is a binary available. To find out which device the program binary
	// in the array refers to, use the ProgramDevicesInfo query to get the list of devices. There is a one-to-one
	// correspondence between the array of N pointers returned by ProgramBinariesInfo and array of devices returned
	// by ProgramDevicesInfo.
	//
	// Returned type: []unsafe.Pointer (pointing to byte arrays)
	ProgramBinariesInfo ProgramInfoName = C.CL_PROGRAM_BINARIES
	// ProgramNumKernelsInfo returns the number of kernels declared in program that can be created with CreateKernel().
	// This information is only available after a successful program executable has been built for at least one device
	// in the list of devices associated with program.
	//
	// Returned type: uintptr
	// Since: 1.2
	ProgramNumKernelsInfo ProgramInfoName = C.CL_PROGRAM_NUM_KERNELS
	// ProgramKernelNamesInfo returns a semi-colon separated list of kernel names in program that can be created
	// with CreateKernel(). This information is only available after a successful program executable has been built
	// for at least one device in the list of devices associated with program.
	//
	// Returned type: string
	// Since: 1.2
	ProgramKernelNamesInfo ProgramInfoName = C.CL_PROGRAM_KERNEL_NAMES
	// ProgramIlInfo returns the program intermediate language (IL) for programs created with CreateProgramWithIl().
	//
	// If program is created with CreateProgramWithSource(), CreateProgramWithBinary(), or
	// CreateProgramWithBuiltInKernels() the memory pointed to by paramValue will be unchanged and
	// the returned size be 0.
	//
	// Returned type: []byte
	// Since: 2.1
	ProgramIlInfo ProgramInfoName = C.CL_PROGRAM_IL
	// ProgramScopeGlobalCtorsPresentInfo indicates that the program object contains non-trivial constructor(s)
	// that will be executed by runtime before any kernel from the program is executed.
	// This information is only available after a successful program executable has been built for at least one
	// device in the list of devices associated with program.
	//
	// Querying ProgramScopeGlobalCtorsPresentInfo may unconditionally return False if no devices associated
	// with program support constructors for program scope global variables. Support for constructors and destructors
	// for program scope global variables is required only for OpenCL 2.2 devices.
	//
	// Returned type: Bool
	// Since: 2.2
	ProgramScopeGlobalCtorsPresentInfo ProgramInfoName = C.CL_PROGRAM_SCOPE_GLOBAL_CTORS_PRESENT
	// ProgramScopeGlobalDtorsPresentInfo indicates that the program object contains non-trivial destructor(s)
	// that will be executed by runtime when program is destroyed.
	// This information is only available after a successful program executable has been built for at least one
	// device in the list of devices associated with program.
	//
	// Querying ProgramScopeGlobalDtorsPresentInfo may unconditionally return False if no devices associated
	// with program support destructors for program scope global variables. Support for constructors and destructors
	// for program scope global variables is required only for OpenCL 2.2 devices.
	//
	// Returned type: Bool
	// Since: 2.2
	ProgramScopeGlobalDtorsPresentInfo ProgramInfoName = C.CL_PROGRAM_SCOPE_GLOBAL_DTORS_PRESENT
)

// ProgramInfo returns information of the program object.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use ProgramInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetProgramInfo.html
func ProgramInfo(program Program, paramName ProgramInfoName, paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetProgramInfo(
		program.handle(),
		C.cl_program_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// ProgramInfoString is a convenience method for ProgramInfo() to query information values that are string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func ProgramInfoString(program Program, paramName ProgramInfoName) (string, error) {
	return queryString(func(paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
		return ProgramInfo(program, paramName, paramSize, paramValue)
	})
}
