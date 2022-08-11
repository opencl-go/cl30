package cl30

// #include "api.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// DeviceID references an OpenCL device of the system.
type DeviceID uintptr

func (id DeviceID) handle() C.cl_device_id {
	return *(*C.cl_device_id)(unsafe.Pointer(&id))
}

// String provides a readable presentation of the device identifier.
// It is based on the numerical value of the underlying pointer.
func (id DeviceID) String() string {
	return fmt.Sprintf("0x%X", uintptr(id))
}

// DeviceTypeFlags is a bitfield that identifies the type of OpenCL device.
// It can be used for DeviceIDs() to filter for the requested devices.
type DeviceTypeFlags C.cl_device_type

const (
	// DeviceTypeCPU is an OpenCL device similar to a traditional CPU (Central Processing Unit).
	// The host processor that executes OpenCL host code may also be considered a CPU OpenCL device.
	DeviceTypeCPU DeviceTypeFlags = C.CL_DEVICE_TYPE_CPU
	// DeviceTypeDefault is the default OpenCL device in the platform.
	// The default OpenCL device must not be a DeviceTypeCustom device.
	DeviceTypeDefault DeviceTypeFlags = C.CL_DEVICE_TYPE_DEFAULT
	// DeviceTypeGpu is an OpenCL device similar to a GPU (Graphics Processing Unit).
	// Many systems include a dedicated processor for graphics or rendering that may be considered a GPU OpenCL device.
	DeviceTypeGpu DeviceTypeFlags = C.CL_DEVICE_TYPE_GPU
	// DeviceTypeAccelerator are dedicated devices that may accelerate OpenCL programs, such as FPGAs
	// (Field Programmable Gate Arrays), DSPs (Digital Signal Processors), or AI (Artificial Intelligence) processors.
	DeviceTypeAccelerator DeviceTypeFlags = C.CL_DEVICE_TYPE_ACCELERATOR
	// DeviceTypeCustom are specialized devices that implement some OpenCL runtime APIs but do not support
	// all required OpenCL functionality.
	//
	// Since: 1.2
	DeviceTypeCustom DeviceTypeFlags = C.CL_DEVICE_TYPE_CUSTOM
	// DeviceTypeAll identifies all OpenCL devices available in the platform, except for DeviceTypeCustom devices.
	DeviceTypeAll DeviceTypeFlags = C.CL_DEVICE_TYPE_ALL
)

// DeviceIDs queries devices available on a platform.
//
// The deviceType is a bitfield that identifies the type of OpenCL device. The deviceType can be used to query
// specific OpenCL devices or all OpenCL devices available.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetDeviceIDs.html
func DeviceIDs(platformID PlatformID, deviceType DeviceTypeFlags) ([]DeviceID, error) {
	count := C.cl_uint(0)
	status := C.clGetDeviceIDs(platformID.handle(), C.cl_device_type(deviceType), 0, nil, &count)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	if count == 0 {
		return nil, nil
	}
	ids := make([]DeviceID, count)
	status = C.clGetDeviceIDs(platformID.handle(), C.cl_device_type(deviceType), count, (*C.cl_device_id)(unsafe.Pointer(&ids[0])), &count)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	return ids[:count], nil
}

// DeviceInfoName identifies properties of a device, which can be queried with DeviceInfo().
type DeviceInfoName C.cl_device_info

const (
	// DeviceAddressBits specify the default compute device address space size of the global address space specified
	// as an unsigned integer value in bits. Currently supported values are 32 or 64 bits.
	//
	// Returned type: uint
	DeviceAddressBits DeviceInfoName = C.CL_DEVICE_ADDRESS_BITS
	// DeviceAtomicFenceCapabilities describes the various memory orders and scopes that the device supports for
	// atomic fence operations. This is a bit-field that has the same set of possible values as described for
	// DeviceAtomicMemoryCapabilities.
	//
	// Because atomic scopes are hierarchical, a device that supports a wide scope must also support all narrower
	// scopes, except for the work-item scope, which is a special case.
	//
	// Returned type: DeviceAtomicCapabilitiesFlags
	// Since: 3.0
	DeviceAtomicFenceCapabilities DeviceInfoName = C.CL_DEVICE_ATOMIC_FENCE_CAPABILITIES
	// DeviceAtomicMemoryCapabilities describes the various memory orders and scopes that the device supports for
	// atomic memory operations.
	//
	// Because atomic scopes are hierarchical, a device that supports a wide scope must also support all narrower
	// scopes, except for the work-item scope, which is a special case.
	//
	// Returned type: DeviceAtomicCapabilitiesFlags
	// Since: 3.0
	DeviceAtomicMemoryCapabilities DeviceInfoName = C.CL_DEVICE_ATOMIC_MEMORY_CAPABILITIES
	// DeviceAvailable is true if the device is available and false otherwise. A device is considered to be
	// available if the device can be expected to successfully execute commands enqueued to the device.
	//
	// Returned type: Bool
	DeviceAvailable DeviceInfoName = C.CL_DEVICE_AVAILABLE
	// DeviceBuiltInKernels represents a semicolon separated list of built-in kernels supported by the device.
	// An empty string is returned if no built-in kernels are supported by the device.
	//
	// Returned type: string
	// Since: 1.2
	DeviceBuiltInKernels DeviceInfoName = C.CL_DEVICE_BUILT_IN_KERNELS
	// DeviceCompilerAvailable is False if the implementation does not have a compiler available to compile the
	// program source. It is True if the compiler is available.
	//
	// This can be False for the embedded platform profile only.
	//
	// Returned type: Bool
	DeviceCompilerAvailable DeviceInfoName = C.CL_DEVICE_COMPILER_AVAILABLE
	// DeviceDeviceEnqueueCapabilities describes device-side enqueue capabilities of the device. This is a bit-field.
	//
	// If DeviceQueueReplaceableDefault is set, DeviceQueueSupported must also be set.
	//
	// Devices that set DeviceQueueSupported for DeviceDeviceEnqueueCapabilities must also return True for
	// DeviceGenericAddressSpaceSupport.
	//
	// Returned type: DeviceDeviceEnqueueCapabilitiesFlags
	// Since: 3.0
	DeviceDeviceEnqueueCapabilities DeviceInfoName = C.CL_DEVICE_DEVICE_ENQUEUE_CAPABILITIES
	// DeviceDoubleFpConfig describes double precision floating-point capability of the OpenCL device.
	// This is a bit-field.
	//
	// Returned type: DeviceFpConfigFlags
	// Since: 1.2
	DeviceDoubleFpConfig DeviceInfoName = C.CL_DEVICE_DOUBLE_FP_CONFIG
	// DeviceEndianLittle is True if the OpenCL device is a little endian device and False otherwise.
	//
	// Returned type: Bool
	DeviceEndianLittle DeviceInfoName = C.CL_DEVICE_ENDIAN_LITTLE
	// DeviceErrorCorrectionSupport is True if the device implements error correction for all accesses to compute
	// device memory (global and constant). Is False if the device does not implement such error correction.
	//
	// Returned type: Bool
	DeviceErrorCorrectionSupport DeviceInfoName = C.CL_DEVICE_ERROR_CORRECTION_SUPPORT
	// DeviceExecutionCapabilities describes the execution capabilities of the device. This is a bit-field.
	//
	// The mandated minimum capability is ExecKernel.
	//
	// Returned type: DeviceExecCapabilitiesFlags
	DeviceExecutionCapabilities DeviceInfoName = C.CL_DEVICE_EXECUTION_CAPABILITIES
	// DeviceExtensions returns a space separated list of extension names (the extension names themselves do not
	// contain any spaces) supported by the device. The list of extension names may include Khronos approved
	// extension names and vendor specified extension names.
	//
	// Returned type: string
	DeviceExtensions DeviceInfoName = C.CL_DEVICE_EXTENSIONS
	// DeviceExtensionsWithVersion returns an array of description (name and version) structures. The same
	// extension name must not be reported more than once. The list of extensions reported must match the list
	// reported via DeviceExtensions.
	//
	// Returned type: []NameVersion
	// Since: 3.0
	DeviceExtensionsWithVersion DeviceInfoName = C.CL_DEVICE_EXTENSIONS_WITH_VERSION
	// DeviceGenericAddressSpaceSupport is True if the device supports the generic address space and its
	// associated built-in functions, and False otherwise.
	//
	// Returned type: Bool
	// Since: 3.0
	DeviceGenericAddressSpaceSupport DeviceInfoName = C.CL_DEVICE_GENERIC_ADDRESS_SPACE_SUPPORT
	// DeviceGlobalMemCacheSize returns the size of global memory cache in bytes.
	//
	// Returned type: Ulong
	DeviceGlobalMemCacheSize DeviceInfoName = C.CL_DEVICE_GLOBAL_MEM_CACHE_SIZE
	// DeviceGlobalMemCacheType represents the type of global memory cache supported.
	//
	// Returned type: DeviceMemCacheTypeEnum
	DeviceGlobalMemCacheType DeviceInfoName = C.CL_DEVICE_GLOBAL_MEM_CACHE_TYPE
	// DeviceGlobalMemCachelineSize is the size of global memory cache line in bytes.
	//
	// Returned type: Uint
	DeviceGlobalMemCachelineSize DeviceInfoName = C.CL_DEVICE_GLOBAL_MEM_CACHELINE_SIZE
	// DeviceGlobalMemSize is the size of global device memory in bytes.
	//
	// Returned type: Ulong
	DeviceGlobalMemSize DeviceInfoName = C.CL_DEVICE_GLOBAL_MEM_SIZE
	// DeviceGlobalVariablePreferredTotalSize is the maximum preferred total size, in bytes, of all program variables
	// in the global address space. This is a performance hint. An implementation may place such variables in storage
	// with optimized device access. This query returns the capacity of such storage. The minimum value is 0.
	//
	// Returned type: uintptr
	// Since: 2.0
	DeviceGlobalVariablePreferredTotalSize DeviceInfoName = C.CL_DEVICE_GLOBAL_VARIABLE_PREFERRED_TOTAL_SIZE
	// DeviceIlVersion represents the intermediate languages that can be supported by CreateProgramWithIl for this
	// device. Returns a space-separated list of IL version strings of the form
	// <IL_Prefix>_<Major_Version>.<Minor_Version>.
	//
	// For an OpenCL 2.1 or 2.2 device, SPIR-V is a required IL prefix.
	//
	// If the device does not support intermediate language programs, the value must be "" (an empty string).
	//
	// Returned type: string
	// Since: 2.1
	// Extension: cl_khr_il_program
	DeviceIlVersion DeviceInfoName = C.CL_DEVICE_IL_VERSION
	// DeviceIlsWithVersion returns an array of descriptions (name and version) for all supported intermediate
	// languages. Intermediate languages with the same name may be reported more than once but each name and
	// major/minor version combination may only be reported once. The list of intermediate languages reported must
	// match the list reported via DeviceIlVersion.
	//
	// For an OpenCL 2.1 or 2.2 device, at least one version of SPIR-V must be reported.
	//
	// Returned type: []NameVersion
	// Since: 3.0
	// Extension: cl_khr_il_program
	DeviceIlsWithVersion DeviceInfoName = C.CL_DEVICE_ILS_WITH_VERSION
	// DeviceImage2dMaxHeight is the maximum height of 2D image in pixels.
	// The minimum value is 16384 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	DeviceImage2dMaxHeight DeviceInfoName = C.CL_DEVICE_IMAGE2D_MAX_HEIGHT
	// DeviceImage2dMaxWidth is the maximum width of 2D image or 1D image not created from a buffer object in pixels.
	// The minimum value is 16384 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	DeviceImage2dMaxWidth DeviceInfoName = C.CL_DEVICE_IMAGE2D_MAX_WIDTH
	// DeviceImage3dMaxDepth is the maximum depth of 3D image in pixels.
	// The minimum value is 2048 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	DeviceImage3dMaxDepth DeviceInfoName = C.CL_DEVICE_IMAGE3D_MAX_DEPTH
	// DeviceImage3dMaxHeight is the maximum height of 3D image in pixels.
	// The minimum value is 2048 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	DeviceImage3dMaxHeight DeviceInfoName = C.CL_DEVICE_IMAGE3D_MAX_HEIGHT
	// DeviceImage3dMaxWidth is the maximum width of 3D image in pixels.
	// The minimum value is 2048 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	DeviceImage3dMaxWidth DeviceInfoName = C.CL_DEVICE_IMAGE3D_MAX_WIDTH
	// DeviceImageBaseAddressAlignment specifies the minimum alignment in pixels of the host_ptr specified to
	// CreateBuffer() or CreateBufferWithProperties() when a 2D image is created from a buffer which was created
	// using MemUseHostPtr. The value returned must be a power of 2.
	//
	// This value must be 0 for devices that do not support 2D images created from a buffer.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceImageBaseAddressAlignment DeviceInfoName = C.CL_DEVICE_IMAGE_BASE_ADDRESS_ALIGNMENT
	// DeviceImageMaxArraySize is the maximum number of images in a 1D or 2D image array.
	// The minimum value is 2048 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	// Since: 1.2
	DeviceImageMaxArraySize DeviceInfoName = C.CL_DEVICE_IMAGE_MAX_ARRAY_SIZE
	// DeviceImageMaxBufferSize is the maximum number of pixels for a 1D image created from a buffer object.
	// The minimum value is 65536 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: uintptr
	// Since: 1.2
	DeviceImageMaxBufferSize DeviceInfoName = C.CL_DEVICE_IMAGE_MAX_BUFFER_SIZE
	// DeviceImagePitchAlignment is the row pitch alignment size in pixels for 2D images created from a buffer.
	// The value returned must be a power of 2.
	// This value must be 0 for devices that do not support 2D images created from a buffer.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceImagePitchAlignment DeviceInfoName = C.CL_DEVICE_IMAGE_PITCH_ALIGNMENT
	// DeviceImageSupport is True if images are supported by the OpenCL device and False otherwise.
	DeviceImageSupport DeviceInfoName = C.CL_DEVICE_IMAGE_SUPPORT
	// DeviceLatestConformanceVersionPassed returns the latest version of the conformance test suite that this device
	// has fully passed in accordance with the official conformance process.
	//
	// Returned type: string
	// Since: 3.0
	DeviceLatestConformanceVersionPassed DeviceInfoName = C.CL_DEVICE_LATEST_CONFORMANCE_VERSION_PASSED
	// DeviceLinkerAvailable is False if the implementation does not have a linker available.
	// It is True if the linker is available.
	//
	// This can be False for the embedded platform profile only.
	// This must be True if DeviceCompilerAvailable is True.
	//
	// Returned type: Bool
	// Since: 1.2
	DeviceLinkerAvailable DeviceInfoName = C.CL_DEVICE_LINKER_AVAILABLE
	// DeviceLocalMemSize is the size of local memory region in bytes. The minimum value is 32 KB for devices
	// that are not of type DeviceTypeCustom.
	//
	// Returned type: Ulong
	DeviceLocalMemSize DeviceInfoName = C.CL_DEVICE_LOCAL_MEM_SIZE
	// DeviceLocalMemType is the type of local memory supported.
	// This can be set to DeviceLocalMemTypeLocal implying dedicated local memory storage such as SRAM, or
	// DeviceLocalMemTypeGlobal.
	//
	// For custom devices, DeviceLocalMemTypeNone can also be returned indicating no local memory support.
	//
	// Returned type: DeviceLocalMemTypeEnum
	DeviceLocalMemType DeviceInfoName = C.CL_DEVICE_LOCAL_MEM_TYPE
	// DeviceMaxClockFrequency is the clock frequency of the device in MHz. The meaning of this value is
	// implementation-defined. For devices with multiple clock domains, the clock frequency for any of the clock
	// domains may be returned. For devices that dynamically change frequency for power or thermal reasons, the
	// returned clock frequency may be any valid frequency.
	//
	// Note: This definition is missing before version 2.2.
	//
	// Returned type: Uint
	// Deprecated: by 2.2
	DeviceMaxClockFrequency DeviceInfoName = C.CL_DEVICE_MAX_CLOCK_FREQUENCY
	// DeviceMaxComputeUnits refers to the number of parallel compute units on the OpenCL device.
	// A work-group executes on a single compute unit. The minimum value is 1.
	//
	// Returned type: Uint
	DeviceMaxComputeUnits DeviceInfoName = C.CL_DEVICE_MAX_COMPUTE_UNITS
	// DeviceMaxConstantArgs is the maximum number of arguments declared with the __constant qualifier in a kernel.
	// The minimum value is 8 for devices that are not of type DeviceTypeCustom.
	//
	// Returned type: Uint
	DeviceMaxConstantArgs DeviceInfoName = C.CL_DEVICE_MAX_CONSTANT_ARGS
	// DeviceMaxConstantBufferSize is the maximum size in bytes of a constant buffer allocation. The minimum value
	// is 64 KB for devices that are not of type DeviceTypeCustom.
	//
	// Returned type: Ulong
	DeviceMaxConstantBufferSize DeviceInfoName = C.CL_DEVICE_MAX_CONSTANT_BUFFER_SIZE
	// DeviceMaxGlobalVariableSize is the maximum number of bytes of storage that may be allocated for any single
	// variable in program scope or inside a function in an OpenCL kernel language declared in the global address space.
	//
	// The minimum value is 64 KB if the device supports program scope global variables, and must be 0 for devices
	// that do not support program scope global variables.
	//
	// Returned type: uintptr
	// Since: 2.0
	DeviceMaxGlobalVariableSize DeviceInfoName = C.CL_DEVICE_MAX_GLOBAL_VARIABLE_SIZE
	// DeviceMaxMemAllocSize is the maximum size of memory object allocation in bytes. The minimum value is
	// max(min(1024 * 1024 * 1024, 1/4th of DeviceGlobalMemSize), 32 * 1024 * 1024)
	// for devices that are not of type DeviceTypeCustom.
	//
	// Returned type: Ulong
	DeviceMaxMemAllocSize DeviceInfoName = C.CL_DEVICE_MAX_MEM_ALLOC_SIZE
	// DeviceMaxNumSubGroups is the maximum number of subgroups in a work-group that a device is capable of executing
	// on a single compute unit, for any given kernel-instance running on the device.
	//
	// The minimum value is 1 if the device supports subgroups, and must be 0 for devices that do not support subgroups.
	//
	// Returned type: Uint
	// Since: 2.1
	DeviceMaxNumSubGroups DeviceInfoName = C.CL_DEVICE_MAX_NUM_SUB_GROUPS
	// DeviceMaxOnDeviceEvents is the maximum number of events in use by a device queue. These refer to events
	// returned by the enqueue_ built-in functions to a device queue or user events returned by the create_user_event
	// built-in function that have not been released.
	//
	// The minimum value is 1024 for devices supporting on-device queues, and must be 0 for devices that do not
	// support on-device queues.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceMaxOnDeviceEvents DeviceInfoName = C.CL_DEVICE_MAX_ON_DEVICE_EVENTS
	// DeviceMaxOnDeviceQueues is the maximum number of device queues that can be created for this device in a
	// single context.
	//
	// The minimum value is 1 for devices supporting on-device queues, and must be 0 for devices that do not
	// support on-device queues.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceMaxOnDeviceQueues DeviceInfoName = C.CL_DEVICE_MAX_ON_DEVICE_QUEUES
	// DeviceMaxParameterSize is the maximum size in bytes of all arguments that can be passed to a kernel.
	//
	// The minimum value is 1024 for devices that are not of type DeviceTypeCustom. For this minimum value,
	// only a maximum of 128 arguments can be passed to a kernel
	//
	// Returned type: uintptr
	DeviceMaxParameterSize DeviceInfoName = C.CL_DEVICE_MAX_PARAMETER_SIZE
	// DeviceMaxPipeArgs is the maximum number of pipe objects that can be passed as arguments to a kernel.
	// The minimum value is 16 for devices supporting pipes, and must be 0 for devices that do not support pipes.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceMaxPipeArgs DeviceInfoName = C.CL_DEVICE_MAX_PIPE_ARGS
	// DeviceMaxReadImageArgs is the maximum number of image objects arguments of a kernel declared with the read_only
	// qualifier. The minimum value is 128 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: Uint
	DeviceMaxReadImageArgs DeviceInfoName = C.CL_DEVICE_MAX_READ_IMAGE_ARGS
	// DeviceMaxReadWriteImageArgs is the maximum number of image objects arguments of a kernel declared with the
	// write_only or read_write qualifier.
	//
	// The minimum value is 64 if the device supports read-write images arguments, and must be 0 for devices that
	// do not support read-write images.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceMaxReadWriteImageArgs DeviceInfoName = C.CL_DEVICE_MAX_READ_WRITE_IMAGE_ARGS
	// DeviceMaxSamplers is the maximum number of samplers that can be used in a kernel.
	// The minimum value is 16 if DeviceImageSupport is True, the value is 0 otherwise.
	//
	// Returned type: Uint
	DeviceMaxSamplers DeviceInfoName = C.CL_DEVICE_MAX_SAMPLERS
	// DeviceMaxWorkGroupSize is the maximum number of work-items in a work-group that a device is capable of
	// executing on a single compute unit, for any given kernel-instance running on the device.
	// The minimum value is 1. The returned value is an upper limit and will not necessarily maximize performance.
	// This maximum may be larger than supported by a specific kernel.
	//
	// Returned type: uintptr
	DeviceMaxWorkGroupSize DeviceInfoName = C.CL_DEVICE_MAX_WORK_GROUP_SIZE
	// DeviceMaxWorkItemDimensions is the maximum dimensions that specify the global and local work-item IDs used by
	// the data parallel execution model. The minimum value is 3 for devices that are not of type DeviceTypeCustom.
	//
	// Return type: Uint
	DeviceMaxWorkItemDimensions DeviceInfoName = C.CL_DEVICE_MAX_WORK_ITEM_DIMENSIONS
	// DeviceMaxWorkItemSizes is the maximum number of work-items that can be specified in each dimension of the
	// work-group to EnqueueNDRangeKernel().
	// Returns N uintptr entries, where N is the value returned by the query for DeviceMaxWorkItemDimensions.
	// The minimum value is (1, 1, 1) for devices that are not of type DeviceTypeCustom.
	//
	// Returned type: []uintptr
	DeviceMaxWorkItemSizes DeviceInfoName = C.CL_DEVICE_MAX_WORK_ITEM_SIZES
	// DeviceMaxWriteImageArgs is the maximum number of image objects arguments of a kernel declared with the
	// write_only qualifier. The minimum value is 64 if DeviceImageSupport is True, the value is 0 otherwise.
	DeviceMaxWriteImageArgs DeviceInfoName = C.CL_DEVICE_MAX_WRITE_IMAGE_ARGS
	// DeviceMemBaseAddrAlign is the alignment requirement (in bits) for sub-buffer offsets. The minimum value is
	// the size (in bits) of the largest OpenCL built-in data type supported by the device
	// (long16 in FULL profile, long16 or int16 in EMBEDDED profile) for devices that are not of type DeviceTypeCustom.
	//
	// Returned type: Uint
	DeviceMemBaseAddrAlign DeviceInfoName = C.CL_DEVICE_MEM_BASE_ADDR_ALIGN
	// DeviceName refers to a human-readable string that identifies the device.
	//
	// Returned type: string
	DeviceName DeviceInfoName = C.CL_DEVICE_NAME
	// DeviceNativeVectorWidthChar returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// Returned type: Uint
	// Since: 1.1
	DeviceNativeVectorWidthChar DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_CHAR
	// DeviceNativeVectorWidthDouble returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// If double precision is not supported, this value must be 0.
	//
	// Returned type: Uint
	// Since: 1.1
	DeviceNativeVectorWidthDouble DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_DOUBLE
	// DeviceNativeVectorWidthFloat returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// Returned type: Uint
	// Since: 1.1
	DeviceNativeVectorWidthFloat DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_FLOAT
	// DeviceNativeVectorWidthHalf returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// If the cl_khr_fp16 extension is not supported, this value must be 0.
	//
	// Returned type: Uint
	// Since: 1.1
	// Extension: cl_khr_fp16
	DeviceNativeVectorWidthHalf DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_HALF
	// DeviceNativeVectorWidthInt returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// Returned type: Uint
	// Since: 1.1
	DeviceNativeVectorWidthInt DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_INT
	// DeviceNativeVectorWidthLong returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// Returned type: Uint
	// Since: 1.1
	DeviceNativeVectorWidthLong DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_LONG
	// DeviceNativeVectorWidthShort returns the native ISA vector width.
	// The vector width is defined as the number of scalar elements that can be stored in the vector.
	//
	// Returned type: Uint
	// Since: 1.1
	DeviceNativeVectorWidthShort DeviceInfoName = C.CL_DEVICE_NATIVE_VECTOR_WIDTH_SHORT
	// DeviceNonUniformWorkGroupSupport is True if the device supports non-uniform work-groups, and False otherwise.
	//
	// Returned type: Bool
	// Since: 3.0
	DeviceNonUniformWorkGroupSupport DeviceInfoName = C.CL_DEVICE_NON_UNIFORM_WORK_GROUP_SUPPORT
	// DeviceOpenClCAllVersions returns an array of name, version descriptions listing all the versions of OpenCL C
	// supported by the compiler for the device. In each returned description structure, the name field is required
	// to be "OpenCL C". The list may include both newer non-backwards compatible OpenCL C versions, such as
	// OpenCL C 3.0, and older OpenCL C versions with mandatory backwards compatibility.
	// The version returned by DeviceOpenClCVersion is required to be present in the list.
	//
	// For devices that support compilation from OpenCL C source:
	//
	// Because OpenCL 3.0 is backwards compatible with OpenCL C 1.2, and OpenCL C 1.2 is backwards compatible
	// with OpenCL C 1.1 and OpenCL C 1.0, support for at least OpenCL C 3.0, OpenCL C 1.2, OpenCL C 1.1, and
	// OpenCL C 1.0 is required for an OpenCL 3.0 device.
	//
	// Support for OpenCL C 2.0, OpenCL C 1.2, OpenCL C 1.1, and OpenCL C 1.0 is required for
	// an OpenCL 2.0, OpenCL 2.1, or OpenCL 2.2 device.
	//
	// Support for OpenCL C 1.2, OpenCL C 1.1, and OpenCL C 1.0 is required for an OpenCL 1.2 device.
	//
	// Support for OpenCL C 1.1 and OpenCL C 1.0 is required for an OpenCL 1.1 device.
	//
	// Support for at least OpenCL C 1.0 is required for an OpenCL 1.0 device.
	//
	// For devices that do not support compilation from OpenCL C source, this query may return an empty array.
	//
	// Returned type: []NameVersion
	// Since: 3.0
	DeviceOpenClCAllVersions DeviceInfoName = C.CL_DEVICE_OPENCL_C_ALL_VERSIONS
	// DeviceOpenClCFeatures returns an array of optional OpenCL C features supported by the compiler for
	// the device alongside the OpenCL C version that introduced the feature macro. For example, if a compiler
	// supports an OpenCL C 3.0 feature, the returned name will be the full name of the OpenCL C feature macro,
	// and the returned version will be 3.0.0.
	//
	// For devices that do not support compilation from OpenCL C source, this query may return an empty array.
	//
	// Returned type: []NameVersion
	// Since: 3.0
	DeviceOpenClCFeatures DeviceInfoName = C.CL_DEVICE_OPENCL_C_FEATURES
	// DeviceOpenClCVersion returns the highest fully backwards compatible OpenCL C version supported by the
	// compiler for the device. For devices supporting compilation from OpenCL C source, this will return
	// a version string with the following format:
	//
	// OpenCL<space>C<space><major_version.minor_version><space><vendor-specific information>
	//
	// Returned type: string
	// Since: 1.1
	// Deprecated: 3.0; This query has been superseded by the DeviceOpenClCAllVersions query,
	// which returns a set of OpenCL C versions supported by a device.
	DeviceOpenClCVersion DeviceInfoName = C.CL_DEVICE_OPENCL_C_VERSION
	// DeviceParentDevice returns the DeviceID of the parent device to which this sub-device belongs.
	// If device is a root-level device, a zero value is returned.
	//
	// Returned type: DeviceID
	// Since: 1.2
	DeviceParentDevice DeviceInfoName = C.CL_DEVICE_PARENT_DEVICE
	// DevicePartitionAffinityDomain returns the list of supported affinity domains for partitioning the device
	// using DevicePartitionByAffinityDomain. This is a bit-field.
	// If the device does not support any affinity domains, a value of 0 will be returned.
	//
	// Returned type: DeviceAffinityDomainFlags
	// Since: 1.2
	DevicePartitionAffinityDomain DeviceInfoName = C.CL_DEVICE_PARTITION_AFFINITY_DOMAIN
	// DevicePartitionMaxSubDevices returns the maximum number of sub-devices that can be created when
	// a device is partitioned.
	// The value returned cannot exceed DeviceMaxComputeUnits.
	//
	// Returned type: Uint
	// Since: 1.2
	DevicePartitionMaxSubDevices DeviceInfoName = C.CL_DEVICE_PARTITION_MAX_SUB_DEVICES
	// DevicePartitionProperties returns the list of partition types supported by device.
	// This is an array of uintptr values drawn from the list of DevicePartitionEqually, DevicePartitionByCounts,
	// and DevicePartitionByAffinityDomain.
	// If the device cannot be partitioned (i.e. there is no partitioning scheme supported by the device that will
	// return at least two subdevices), a value of 0 will be returned.
	//
	// Returned type: []uintptr
	// Since: 1.2
	DevicePartitionProperties DeviceInfoName = C.CL_DEVICE_PARTITION_PROPERTIES
	// DevicePartitionType returns the properties argument specified in CreateSubDevices() if device is a sub-device.
	// In the case where the properties argument to CreateSubDevices() is DevicePartitionByAffinityDomain,
	// DeviceAffinityDomainNextPartitionable, the affinity domain used to perform the partition will be returned.
	// This can be one of the following values:
	//
	// DeviceAffinityDomainNuma
	// DeviceAffinityDomainL4Cache
	// DeviceAffinityDomainL3Cache
	// DeviceAffinityDomainL2Cache
	// DeviceAffinityDomainL1Cache
	//
	// Otherwise the implementation may either return a parameter size of 0 i.e. there is no partition type
	// associated with device or can return a property value of 0 (where 0 is used to terminate the
	// partition property list) in the memory that the result value points to.
	//
	// Returned type: []uintptr
	// Since: 1.2
	DevicePartitionType DeviceInfoName = C.CL_DEVICE_PARTITION_TYPE
	// DevicePipeMaxActiveReservations is the maximum number of reservations that can be active for a pipe per
	// work-item in a kernel. A work-group reservation is counted as one reservation per work-item.
	// The minimum value is 1 for devices supporting pipes, and must be 0 for devices that do not support pipes.
	//
	// Returned type: Uint
	// Since: 2.0
	DevicePipeMaxActiveReservations DeviceInfoName = C.CL_DEVICE_PIPE_MAX_ACTIVE_RESERVATIONS
	// DevicePipeMaxPacketSize is the maximum size of pipe packet in bytes.
	// The minimum value is 1024 bytes if the device supports pipes, and must be 0 for devices that do not
	// support pipes.
	//
	// Returned type: Uint
	// Since: 2.0
	DevicePipeMaxPacketSize DeviceInfoName = C.CL_DEVICE_PIPE_MAX_PACKET_SIZE
	// DevicePipeSupport is True if the device supports pipes, and False otherwise.
	// Devices that return True for DevicePipeSupport must also return True for DeviceGenericAddressSpaceSupport.
	//
	// Returned type: Bool
	// Since: 3.0
	DevicePipeSupport DeviceInfoName = C.CL_DEVICE_PIPE_SUPPORT
	// DevicePlatform returns the platform associated with this device.
	//
	// Returned type: PlatformID
	DevicePlatform DeviceInfoName = C.CL_DEVICE_PLATFORM
	// DevicePreferredGlobalAtomicAlignment returns the value representing the preferred alignment in bytes for
	// OpenCL 2.0 atomic types to global memory. This query can return 0 which indicates that the preferred
	// alignment is aligned to the natural size of the type.
	//
	// Returned type: Uint
	// Since: 2.0
	DevicePreferredGlobalAtomicAlignment DeviceInfoName = C.CL_DEVICE_PREFERRED_GLOBAL_ATOMIC_ALIGNMENT
	// DevicePreferredInteropUserSync is True if the devices preference is for the user to be responsible for
	// synchronization, when sharing memory objects between OpenCL and other APIs such as DirectX,
	// False if the device / implementation has a performant path for performing synchronization of memory object
	// shared between OpenCL and other APIs such as DirectX.
	//
	// Returned type: Bool
	// Since: 1.2
	DevicePreferredInteropUserSync DeviceInfoName = C.CL_DEVICE_PREFERRED_INTEROP_USER_SYNC
	// DevicePreferredLocalAtomicAlignment returns the value representing the preferred alignment in bytes for
	// OpenCL 2.0 atomic types to local memory. This query can return 0 which indicates that the preferred
	// alignment is aligned to the natural size of the type.
	//
	// Returned type: Uint
	// Since: 2.0
	DevicePreferredLocalAtomicAlignment DeviceInfoName = C.CL_DEVICE_PREFERRED_LOCAL_ATOMIC_ALIGNMENT
	// DevicePreferredPlatformAtomicAlignment returns the value representing the preferred alignment in bytes for
	// OpenCL 2.0 fine-grained SVM atomic types. This query can return 0 which indicates that the preferred
	// alignment is aligned to the natural size of the type.
	//
	// Returned type: Uint
	// Since: 2.0
	DevicePreferredPlatformAtomicAlignment DeviceInfoName = C.CL_DEVICE_PREFERRED_PLATFORM_ATOMIC_ALIGNMENT
	// DevicePreferredVectorWidthChar is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	//
	// Returned type: Uint
	DevicePreferredVectorWidthChar DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_CHAR
	// DevicePreferredVectorWidthDouble is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	// If double precision is not supported, this value must be 0.
	//
	// Returned type: Uint
	DevicePreferredVectorWidthDouble DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_DOUBLE
	// DevicePreferredVectorWidthFloat is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	//
	// Returned type: Uint
	DevicePreferredVectorWidthFloat DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_FLOAT
	// DevicePreferredVectorWidthHalf is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	// If the cl_khr_fp16 extension is not supported, this value must be 0.
	//
	// Returned type: Uint
	// Since: 1.1
	// Extension: cl_khr_fp16
	DevicePreferredVectorWidthHalf DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_HALF
	// DevicePreferredVectorWidthInt is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	//
	// Returned type: Uint
	DevicePreferredVectorWidthInt DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_INT
	// DevicePreferredVectorWidthLong is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	//
	// Returned type: Uint
	DevicePreferredVectorWidthLong DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_LONG
	// DevicePreferredVectorWidthShort is the preferred native vector width size for built-in scalar types that
	// can be put into vectors. The vector width is defined as the number of scalar elements that can be stored
	// in the vector.
	//
	// Returned type: Uint
	DevicePreferredVectorWidthShort DeviceInfoName = C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_SHORT
	// DevicePrintfBufferSize is the maximum size in bytes of the internal buffer that holds the output of printf
	// calls from a kernel. The minimum value for the FULL profile is 1 MB.
	//
	// Returned type: uintptr
	// Since: 1.2
	DevicePrintfBufferSize DeviceInfoName = C.CL_DEVICE_PRINTF_BUFFER_SIZE
	// DeviceProfile is the OpenCL profile string. Returns the profile name supported by the device.
	// The profile name returned can be one of the following strings:
	//
	// "FULL_PROFILE" - if the device supports the OpenCL specification (functionality defined as part of the core
	// specification and does not require any extensions to be supported).
	//
	// "EMBEDDED_PROFILE" - if the device supports the OpenCL embedded profile.
	//
	// Returned type: string
	DeviceProfile DeviceInfoName = C.CL_DEVICE_PROFILE
	// DeviceProfilingTimerResolution describes the resolution of device timer. This is measured in nanoseconds.
	//
	// Returned type: uintptr
	DeviceProfilingTimerResolution DeviceInfoName = C.CL_DEVICE_PROFILING_TIMER_RESOLUTION
	// DeviceQueueOnDeviceMaxSize is the maximum size of the device queue in bytes.
	// The minimum value is 256 KB for the full profile and 64 KB for the embedded profile for devices supporting
	// on-device queues, and must be 0 for devices that do not support on-device queues.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceQueueOnDeviceMaxSize DeviceInfoName = C.CL_DEVICE_QUEUE_ON_DEVICE_MAX_SIZE
	// DeviceQueueOnDevicePreferredSize is the preferred size of the device queue, in bytes. Applications should
	// use this size for the device queue to ensure good performance.
	// The minimum value is 16 KB for devices supporting on-device queues, and must be 0 for devices that
	// do not support on-device queues.
	//
	// Returned type: Uint
	// Since: 2.0
	DeviceQueueOnDevicePreferredSize DeviceInfoName = C.CL_DEVICE_QUEUE_ON_DEVICE_PREFERRED_SIZE
	// DeviceQueueOnDeviceProperties describes the on device command-queue properties supported by the device.
	// This is a bit-field.
	//
	// Returned type: CommandQueuePropertiesFlags
	// Since: 2.0
	DeviceQueueOnDeviceProperties DeviceInfoName = C.CL_DEVICE_QUEUE_ON_DEVICE_PROPERTIES
	// DeviceQueueOnHostProperties describes the on host command-queue properties supported by the device.
	// This is a bit-field
	//
	// Returned type: CommandQueuePropertiesFlags
	// Since: 2.0
	DeviceQueueOnHostProperties DeviceInfoName = C.CL_DEVICE_QUEUE_ON_HOST_PROPERTIES
	// DeviceReferenceCount returns the device reference count. If the device is a root-level device,
	// a reference count of one is returned.
	//
	// Note: The reference count returned should be considered immediately stale. It is unsuitable for general
	// use in applications. This feature is provided for identifying memory leaks.
	//
	// Returned type: Uint
	// Since: 1.2
	DeviceReferenceCount DeviceInfoName = C.CL_DEVICE_REFERENCE_COUNT
	// DeviceSingleFpConfig describes single precision floating-point capability of the device. This is a bit-field.
	//
	// Returned type: DeviceFpConfigFlags
	DeviceSingleFpConfig DeviceInfoName = C.CL_DEVICE_SINGLE_FP_CONFIG
	// DeviceSubGroupIndependentForwardProgress is True if this device supports independent forward progress of
	// subgroups, False otherwise.
	// This value must be True for devices that support the cl_khr_subgroups extension, and must return False for
	// devices that do not support subgroups.
	//
	// Returned type: Bool
	// Since: 2.1
	// Extension: cl_khr_subgroups
	DeviceSubGroupIndependentForwardProgress DeviceInfoName = C.CL_DEVICE_SUB_GROUP_INDEPENDENT_FORWARD_PROGRESS
	// DeviceSvmCapabilities describes the various shared virtual memory (SVM) memory allocation types the
	// device supports. This is a bit-field.
	//
	// Returned type: DeviceSvmCapabilitiesFlags
	// Since: 2.0
	DeviceSvmCapabilities DeviceInfoName = C.CL_DEVICE_SVM_CAPABILITIES
	// DeviceType is the type or types of the OpenCL device.
	//
	// Returned type: DeviceTypeFlags
	DeviceType DeviceInfoName = C.CL_DEVICE_TYPE
	// DeviceVendor refers to a human-readable string that identifies the vendor of the device.
	//
	// Returned type: string
	DeviceVendor DeviceInfoName = C.CL_DEVICE_VENDOR
	// DeviceVendorID is a unique device vendor identifier.
	//
	// Returned type: Uint
	DeviceVendorID DeviceInfoName = C.CL_DEVICE_VENDOR_ID
	// DeviceVersion is an OpenCL version string. Returns the OpenCL version supported by the device.
	// This version string has the following format:
	//
	// OpenCL<space><major_version.minor_version><space><vendor-specific information>
	//
	// Returned type: string
	DeviceVersion DeviceInfoName = C.CL_DEVICE_VERSION
	// DeviceWorkGroupCollectiveFunctionsSupport is True if the device supports work-group collective functions
	// e.g. work_group_broadcast, work_group_reduce, and work_group_scan, and False otherwise.
	//
	// Returned type: Bool
	// Since: 3.0
	DeviceWorkGroupCollectiveFunctionsSupport DeviceInfoName = C.CL_DEVICE_WORK_GROUP_COLLECTIVE_FUNCTIONS_SUPPORT
	// DriverVersion specifies the OpenCL software driver version string. Follows a vendor-specific format.
	//
	// Returned type: string
	DriverVersion DeviceInfoName = C.CL_DRIVER_VERSION
)

// DeviceAtomicCapabilitiesFlags are used to determine the DeviceAtomicFenceCapabilities
// and DeviceAtomicMemoryCapabilities with DeviceInfo().
type DeviceAtomicCapabilitiesFlags C.cl_device_atomic_capabilities

const (
	// DeviceAtomicOrderRelaxed identifies support for the relaxed memory order.
	//
	// Since: 3.0
	DeviceAtomicOrderRelaxed DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_ORDER_RELAXED
	// DeviceAtomicOrderAcqRel identifies support for the "acquire", "release", and "acquire-release" memory orders.
	//
	// Since: 3.0
	DeviceAtomicOrderAcqRel DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_ORDER_ACQ_REL
	// DeviceAtomicOrderSeqCst identifies support for the sequentially consistent memory order.
	//
	// Since: 3.0
	DeviceAtomicOrderSeqCst DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_ORDER_SEQ_CST
	// DeviceAtomicScopeWorkItem identifies support for memory ordering constraints that apply to a single work-item.
	//
	// Note that this flag does not provide meaning for atomic memory operations, but only for atomic fence operations
	// in certain circumstances, refer to the Memory Scope section of the OpenCL C specification.
	//
	// Since: 3.0
	DeviceAtomicScopeWorkItem DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_SCOPE_WORK_ITEM
	// DeviceAtomicScopeWorkGroup identifies support for memory ordering constraints that apply to all work-items
	// in a work-group.
	//
	// Since: 3.0
	DeviceAtomicScopeWorkGroup DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_SCOPE_WORK_GROUP
	// DeviceAtomicScopeDevice identifies support for memory ordering constraints that apply to all work-items
	// executing on the device.
	//
	// Since: 3.0
	DeviceAtomicScopeDevice DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_SCOPE_DEVICE
	// DeviceAtomicScopeAllDevices identifies support for memory ordering constraints that apply to all work-items
	// executing across all devices that can share SVM memory with each other and the host process.
	//
	// Since: 3.0
	DeviceAtomicScopeAllDevices DeviceAtomicCapabilitiesFlags = C.CL_DEVICE_ATOMIC_SCOPE_ALL_DEVICES
)

// DeviceDeviceEnqueueCapabilitiesFlags are used to determine the DeviceDeviceEnqueueCapabilities with DeviceInfo().
type DeviceDeviceEnqueueCapabilitiesFlags C.cl_device_device_enqueue_capabilities

const (
	// DeviceQueueSupported identifies that the device supports device-side enqueue and on-device queues.
	//
	// Since: 3.0
	DeviceQueueSupported DeviceDeviceEnqueueCapabilitiesFlags = C.CL_DEVICE_QUEUE_SUPPORTED
	// DeviceQueueReplaceableDefault identifies that the device supports a replaceable default on-device queue.
	//
	// Since: 3.0
	DeviceQueueReplaceableDefault DeviceDeviceEnqueueCapabilitiesFlags = C.CL_DEVICE_QUEUE_REPLACEABLE_DEFAULT
)

// DeviceFpConfigFlags are used to determine the DeviceSingleFpConfig and DeviceDoubleFpConfig with DeviceInfo().
type DeviceFpConfigFlags C.cl_device_fp_config

const (
	// FpDenorm identifies denorms are supported.
	FpDenorm DeviceFpConfigFlags = C.CL_FP_DENORM
	// FpInfNan identifies INF and quiet NaNs are supported.
	FpInfNan DeviceFpConfigFlags = C.CL_FP_INF_NAN
	// FpRoundToNearest identifies round to nearest even rounding mode supported.
	FpRoundToNearest DeviceFpConfigFlags = C.CL_FP_ROUND_TO_NEAREST
	// FpRoundToZero identifies round to zero rounding mode supported.
	FpRoundToZero DeviceFpConfigFlags = C.CL_FP_ROUND_TO_ZERO
	// FpRoundToInf identifies round to positive and negative infinity rounding modes supported.
	FpRoundToInf DeviceFpConfigFlags = C.CL_FP_ROUND_TO_INF
	// FpFma identifies IEEE754-2008 fused multiply-add is supported.
	FpFma DeviceFpConfigFlags = C.CL_FP_FMA
	// FpSoftFloat identifies basic floating-point operations (such as addition, subtraction, multiplication)
	// are implemented in software.
	//
	// Since: 1.1
	FpSoftFloat DeviceFpConfigFlags = C.CL_FP_SOFT_FLOAT
	// FpCorrectlyRoundedDivideSqrt identifies divide and sqrt are correctly rounded as defined by the IEEE754
	// specification.
	//
	// Since: 1.2
	FpCorrectlyRoundedDivideSqrt DeviceFpConfigFlags = C.CL_FP_CORRECTLY_ROUNDED_DIVIDE_SQRT
)

// DeviceExecCapabilitiesFlags are used to determine the DeviceExecutionCapabilities with DeviceInfo().
type DeviceExecCapabilitiesFlags C.cl_device_exec_capabilities

const (
	// ExecKernel identifies that the OpenCL device can execute OpenCL kernels.
	ExecKernel DeviceExecCapabilitiesFlags = C.CL_EXEC_KERNEL
	// ExecNativeKernel identifies that the OpenCL device can execute native kernels.
	ExecNativeKernel DeviceExecCapabilitiesFlags = C.CL_EXEC_NATIVE_KERNEL
)

// DeviceMemCacheTypeEnum is used to determine the DeviceGlobalMemCacheType with DeviceInfo().
type DeviceMemCacheTypeEnum C.cl_device_mem_cache_type

// These are the possible values for DeviceMemCacheTypeEnum. They are slightly reworded compared to the original
// OpenCL API to avoid potential name/type clashes in the future.
const (
	DeviceMemCacheNone      DeviceMemCacheTypeEnum = C.CL_NONE
	DeviceMemCacheReadOnly  DeviceMemCacheTypeEnum = C.CL_READ_ONLY_CACHE
	DeviceMemCacheReadWrite DeviceMemCacheTypeEnum = C.CL_READ_WRITE_CACHE
)

// DeviceLocalMemTypeEnum is used to determine the DeviceLocalMemType with DeviceInfo().
type DeviceLocalMemTypeEnum C.cl_device_local_mem_type

// These are the possible values for DeviceLocalMemTypeEnum. They are slightly reworded compared to the original
// OpenCL API to avoid potential name/type clashes in the future.
const (
	DeviceLocalMemTypeNone   DeviceLocalMemTypeEnum = C.CL_NONE
	DeviceLocalMemTypeLocal  DeviceLocalMemTypeEnum = C.CL_LOCAL
	DeviceLocalMemTypeGlobal DeviceLocalMemTypeEnum = C.CL_GLOBAL
)

// DeviceSvmCapabilitiesFlags is used to determine DeviceSvmCapabilities with DeviceInfo().
type DeviceSvmCapabilitiesFlags C.cl_device_svm_capabilities

const (
	// DeviceSvmCoarseGrainBuffer identifies support for coarse-grain buffer sharing using SvmAlloc().
	// Memory consistency is guaranteed at synchronization points and the host must use calls to EnqueueMapBuffer()
	// and EnqueueUnmapMemObject().
	//
	// Since: 2.0
	DeviceSvmCoarseGrainBuffer DeviceSvmCapabilitiesFlags = C.CL_DEVICE_SVM_COARSE_GRAIN_BUFFER
	// DeviceSvmFineGrainBuffer identifies support for fine-grain buffer sharing using SvmAlloc().
	// Memory consistency is guaranteed at synchronization points without need for EnqueueMapBuffer() and
	// EnqueueUnmapMemObject().
	//
	// Since: 2.0
	DeviceSvmFineGrainBuffer DeviceSvmCapabilitiesFlags = C.CL_DEVICE_SVM_FINE_GRAIN_BUFFER
	// DeviceSvmFineGrainSystem identifies support for sharing the host’s entire virtual memory including memory
	// allocated using malloc. Memory consistency is guaranteed at synchronization points.
	//
	// Since: 2.0
	DeviceSvmFineGrainSystem DeviceSvmCapabilitiesFlags = C.CL_DEVICE_SVM_FINE_GRAIN_SYSTEM
	// DeviceSvmAtomics identifies support for the OpenCL 2.0 atomic operations that provide memory consistency
	// across the host and all OpenCL devices supporting fine-grain SVM allocations.
	//
	// Since: 2.0
	DeviceSvmAtomics DeviceSvmCapabilitiesFlags = C.CL_DEVICE_SVM_ATOMICS
)

// DeviceInfo queries specific information about a device.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use DeviceInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetDeviceInfo.html
func DeviceInfo(id DeviceID, paramName DeviceInfoName, paramSize uint, paramValue unsafe.Pointer) (uint, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetDeviceInfo(
		id.handle(),
		C.cl_device_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uint(sizeReturn), nil
}

// DeviceInfoString is a convenience method for DeviceInfo() to query information values that are string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func DeviceInfoString(id DeviceID, paramName DeviceInfoName) (string, error) {
	return queryString(func(paramSize uint, paramValue unsafe.Pointer) (uint, error) {
		return DeviceInfo(id, paramName, paramSize, paramValue)
	})
}

// DeviceAndHostTimer returns a reasonably synchronized pair of timestamps from the device timer and the host timer
// as seen by device.
//
// The resolution of the device timer may be queried via DeviceInfo() and the flag DeviceProfilingTimerResolution.
// The resolution of the host timer may be queried via PlatformInfo() and the flag PlatformHostTimerResolution.
//
// Since: 2.1
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetDeviceAndHostTimer.html
func DeviceAndHostTimer(id DeviceID) (device uint64, host uint64, err error) {
	status := C.clGetDeviceAndHostTimer(id.handle(), (*C.cl_ulong)(&device), (*C.cl_ulong)(&host))
	if status != C.CL_SUCCESS {
		return 0, 0, StatusError(status)
	}
	return
}

// HostTimer returns the current value of the host clock as seen by device.
// This value is in the same timebase as the host timestamp returned from DeviceAndHostTimer().
//
// Since: 2.1
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetHostTimer.html
func HostTimer(id DeviceID) (uint64, error) {
	var host uint64
	status := C.clGetHostTimer(id.handle(), (*C.cl_ulong)(&host))
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return host, nil
}

const (
	// DevicePartitionEqually requests to split the aggregate device into as many smaller aggregate devices as
	// can be created, each containing N compute units. The value N is passed as the value accompanying this property.
	// If N does not divide evenly into DeviceMaxComputeUnits, then the remaining compute units are not used.
	//
	// Use PartitionedEqually() for convenience.
	//
	// Property value type: Uint
	// Since: 1.2
	DevicePartitionEqually uintptr = C.CL_DEVICE_PARTITION_EQUALLY
	// DevicePartitionByCounts is followed by a list of compute unit counts terminated with 0 or
	// DevicePartitionByCountsListEnd. For each non-zero count M in the list, a sub-device is created with M compute
	// units in it.
	//
	// The number of non-zero count entries in the list may not exceed DevicePartitionMaxSubDevices.
	//
	// The total number of compute units specified may not exceed DeviceMaxComputeUnits.
	//
	// Use PartitionedByCounts() for convenience.
	//
	// Property value type: Uint
	// Since: 1.2
	DevicePartitionByCounts uintptr = C.CL_DEVICE_PARTITION_BY_COUNTS
	// DevicePartitionByCountsListEnd terminates the property value list started by DevicePartitionByCounts.
	//
	// Since: 1.2
	DevicePartitionByCountsListEnd uintptr = C.CL_DEVICE_PARTITION_BY_COUNTS_LIST_END
	// DevicePartitionByAffinityDomain splits the device into smaller aggregate devices containing one or more
	// compute units that all share part of a cache hierarchy. The value accompanying this property may be drawn
	// from the constants of DeviceAffinityDomainFlags.
	//
	// Use PartitionedByAffinityDomain() for convenience.
	//
	// Property value type: DeviceAffinityDomainFlags
	// Since: 1.2
	DevicePartitionByAffinityDomain uintptr = C.CL_DEVICE_PARTITION_BY_AFFINITY_DOMAIN
)

// DevicePartitionProperty is one entry of properties which are taken into account when creating sub-devices
// with CreateSubDevices().
type DevicePartitionProperty []uintptr

// PartitionedEqually is a convenience function to handle the DevicePartitionEqually property.
// Use it in combination with CreateSubDevices().
func PartitionedEqually(units uint) DevicePartitionProperty {
	return DevicePartitionProperty{DevicePartitionEqually, uintptr(units)}
}

// PartitionedByCounts is a convenience function to handle the DevicePartitionByCounts property.
// Use it in combination with CreateSubDevices().
func PartitionedByCounts(units []uint) DevicePartitionProperty {
	values := make(DevicePartitionProperty, 0, len(units)+2)
	values = append(values, DevicePartitionByCounts)
	for _, unit := range units {
		values = append(values, uintptr(unit))
	}
	values = append(values, DevicePartitionByCountsListEnd)
	return values
}

// DeviceAffinityDomainFlags describe how sub-devices are partitioned according to their cache hierarchy.
type DeviceAffinityDomainFlags C.cl_device_affinity_domain

const (
	// DeviceAffinityDomainNuma splits the device into sub-devices comprised of compute units that share a NUMA node.
	//
	// Since: 1.2
	DeviceAffinityDomainNuma DeviceAffinityDomainFlags = C.CL_DEVICE_AFFINITY_DOMAIN_NUMA
	// DeviceAffinityDomainL4Cache splits the device into sub-devices comprised of compute units that share a
	// level 4 data cache.
	//
	// Since: 1.2
	DeviceAffinityDomainL4Cache DeviceAffinityDomainFlags = C.CL_DEVICE_AFFINITY_DOMAIN_L4_CACHE
	// DeviceAffinityDomainL3Cache splits the device into sub-devices comprised of compute units that share a
	// level 3 data cache.
	//
	// Since: 1.2
	DeviceAffinityDomainL3Cache DeviceAffinityDomainFlags = C.CL_DEVICE_AFFINITY_DOMAIN_L3_CACHE
	// DeviceAffinityDomainL2Cache splits the device into sub-devices comprised of compute units that share a
	// level 2 data cache.
	//
	// Since: 1.2
	DeviceAffinityDomainL2Cache DeviceAffinityDomainFlags = C.CL_DEVICE_AFFINITY_DOMAIN_L2_CACHE
	// DeviceAffinityDomainL1Cache splits the device into sub-devices comprised of compute units that share a
	// level 1 data cache.
	//
	// Since: 1.2
	DeviceAffinityDomainL1Cache DeviceAffinityDomainFlags = C.CL_DEVICE_AFFINITY_DOMAIN_L1_CACHE
	// DeviceAffinityDomainNextPartitionable splits the device along the next partitionable affinity domain.
	// The implementation shall find the first level along which the device or sub-device may be further subdivided
	// in the order NUMA, L4, L3, L2, L1, and partition the device into sub-devices comprised of compute units that
	// share memory subsystems at this level.
	//
	// Determine what happened by calling DeviceInfo() with DevicePartitionType on the sub-devices.
	//
	// Since: 1.2
	DeviceAffinityDomainNextPartitionable DeviceAffinityDomainFlags = C.CL_DEVICE_AFFINITY_DOMAIN_NEXT_PARTITIONABLE
)

// PartitionedByAffinityDomain is a convenience function to handle the DevicePartitionByAffinityDomain property.
// Use it in combination with CreateSubDevices().
func PartitionedByAffinityDomain(domain DeviceAffinityDomainFlags) DevicePartitionProperty {
	return DevicePartitionProperty{DevicePartitionByAffinityDomain, uintptr(domain)}
}

// CreateSubDevices creates an array of sub-devices that each reference a non-intersecting set of compute units within
// the device identified by id, according to the partition scheme given by properties.
// Only one of the available partitioning schemes can be specified in properties.
//
// The output sub-devices may be used in every way that the root (or parent) device can be used, including
// creating contexts, building programs, further calls to CreateSubDevices(), and creating command-queues.
// When a command-queue is created against a sub-device, the commands enqueued on the queue are executed only
// on the sub-device.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateSubDevices.html
func CreateSubDevices(id DeviceID, properties ...DevicePartitionProperty) ([]DeviceID, error) {
	var rawPropertyList []uintptr
	for _, property := range properties {
		rawPropertyList = append(rawPropertyList, property...)
	}
	var rawProperties unsafe.Pointer
	if len(properties) > 0 {
		rawPropertyList = append(rawPropertyList, 0)
		rawProperties = unsafe.Pointer(&rawPropertyList[0])
	}

	requiredCount := C.cl_uint(0)
	status := C.clCreateSubDevices(
		id.handle(),
		(*C.cl_device_partition_property)(rawProperties),
		0, nil,
		&requiredCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	ids := make([]DeviceID, requiredCount)
	reportedCount := C.cl_uint(0)
	status = C.clCreateSubDevices(
		id.handle(),
		(*C.cl_device_partition_property)(rawProperties),
		requiredCount,
		(*C.cl_device_id)(unsafe.Pointer(&ids[0])),
		&reportedCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	return ids[:reportedCount], nil
}

// RetainDevice increments the device reference count if device is a valid sub-device created by a call to
// CreateSubDevices(). If id refers to a root level device, meaning a DeviceID returned by DeviceIDs(), the device
// reference count remains unchanged.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clRetainDevice.html
func RetainDevice(id DeviceID) error {
	status := C.clRetainDevice(id.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ReleaseDevice decrements the device reference count if device is a valid sub-device created by a call to
// CreateSubDevices(). If id refers to a root level device, meaning a DeviceID returned by DeviceIDs(), the device
// reference count remains unchanged.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clReleaseDevice.html
func ReleaseDevice(id DeviceID) error {
	status := C.clReleaseDevice(id.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}
