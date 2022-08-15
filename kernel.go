package cl30

// #include "api.h"
// extern cl_int cl30EnqueueNativeKernel(cl_command_queue commandQueue,
//    void *args, size_t argsSize,
//    cl_uint numMemObjects, cl_mem *memList, void const *argsMemLoc,
//    cl_uint waitListCount, cl_event const *waitList,
//    cl_event *event);
import "C"
import (
	"fmt"
	"unsafe"
)

// Kernel object references a particular __kernel function and its arguments for execution.
type Kernel uintptr

func (kernel Kernel) handle() C.cl_kernel {
	return *(*C.cl_kernel)(unsafe.Pointer(&kernel))
}

// String provides a readable presentation of the kernel identifier.
// It is based on the numerical value of the underlying pointer.
func (kernel Kernel) String() string {
	return fmt.Sprintf("0x%X", uintptr(kernel))
}

// CreateKernel creates a kernel object.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateKernel.html
func CreateKernel(program Program, name string) (Kernel, error) {
	rawName := C.CString(name)
	defer C.free(unsafe.Pointer(rawName))
	var status C.cl_int
	kernel := C.clCreateKernel(program.handle(), rawName, &status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Kernel(*((*uintptr)(unsafe.Pointer(&kernel)))), nil
}

// CreateKernelsInProgram creates kernel objects for all kernel functions in a program object.
//
// Kernel objects are not created for any __kernel functions in program that do not have the same function
// definition across all devices for which a program executable has been successfully built.
//
// Kernel objects can only be created once you have a program object with a valid program source or binary loaded
// into the program object and the program executable has been successfully built for one or more devices associated
// with program.
// No changes to the program executable are allowed while there are kernel objects associated with a program object.
// This means that calls to BuildProgram() and CompileProgram() return ErrInvalidOperation if there are kernel
// objects attached to a program object.
//
// The OpenCL context associated with program will be the context associated with kernel.
// The list of devices associated with program are the devices associated with kernel.
// Devices associated with a program object for which a valid program executable has been built can be used to
// execute kernels declared in the program object.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateKernelsInProgram.html
func CreateKernelsInProgram(program Program) ([]Kernel, error) {
	var requiredCount C.cl_uint
	status := C.clCreateKernelsInProgram(program.handle(), 0, nil, &requiredCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	if requiredCount == 0 {
		return nil, nil
	}
	kernels := make([]Kernel, int(requiredCount))
	var returnedCount C.cl_uint
	status = C.clCreateKernelsInProgram(
		program.handle(),
		requiredCount,
		(*C.cl_kernel)(unsafe.Pointer(&kernels[0])),
		&returnedCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	return kernels[:int(returnedCount)], nil
}

// CloneKernel makes a shallow copy of the kernel object.
//
// Cloning is used to make a shallow copy of the kernel object, its arguments and any information passed to the
// kernel object using SetKernelExecInfo(). If the kernel object was ready to be enqueued before copying it,
// the clone of the kernel object is ready to enqueue.
//
// The returned kernel object is an exact copy of the source kernel, with one caveat: the reference count on the
// returned kernel object is set as if it had been returned by CreateKernel().
// The reference count of the source kernel will not be changed.
//
// The resulting kernel will be in the same state as if CreateKernel() is called to create the resultant
// kernel with the same arguments as those used to create the source kernel, the latest call to
// SetKernelArg() or SetKernelArgSVMPointer() for each argument index applied to kernel and the last call to
// SetKernelExecInfo() for each value of the param name parameter are applied to the new kernel object.
//
// All arguments of the new kernel object must be intact, and it may be correctly used in the same situations
// as kernel except those that assume a pre-existing reference count. Setting arguments on the new kernel object
// will not affect the source kernel except insofar as the argument points to a shared underlying entity and in
// that situation behavior is as if two kernel objects had been created and the same argument applied to each.
// Only the data stored in the kernel object is copied; data referenced by the kernels arguments are not copied.
// For example, if a buffer or pointer argument is set on a kernel object, the pointer is copied but the underlying
// memory allocation is not.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCloneKernel.html
func CloneKernel(kernel Kernel) (Kernel, error) {
	var status C.cl_int
	kernelCopy := C.clCloneKernel(kernel.handle(), &status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Kernel(*((*uintptr)(unsafe.Pointer(&kernelCopy)))), nil
}

// RetainKernel increments the kernel reference count.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clRetainKernel.html
func RetainKernel(kernel Kernel) error {
	status := C.clRetainKernel(kernel.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ReleaseKernel decrements the kernel reference count.
//
// The kernel object is deleted once the number of instances that are retained to kernel become zero and the kernel
// object is no longer needed by any enqueued commands that use kernel.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clReleaseKernel.html
func ReleaseKernel(kernel Kernel) error {
	status := C.clReleaseKernel(kernel.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// SetKernelArg sets the argument value for a specific argument of a kernel.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetKernelArg.html
func SetKernelArg(kernel Kernel, index uint32, size uintptr, value unsafe.Pointer) error {
	status := C.clSetKernelArg(
		kernel.handle(),
		C.cl_uint(index),
		C.size_t(size),
		value)
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// SetKernelArgSvmPointer sets an SVM pointer as the argument value for a specific argument of a kernel.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetKernelArgSVMPointer.html
func SetKernelArgSvmPointer(kernel Kernel, index uint32, value unsafe.Pointer) error {
	status := C.clSetKernelArgSVMPointer(kernel.handle(), C.cl_uint(index), value)
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// KernelExecInfoName describes an extra parameter beyond arguments for a kernel.
type KernelExecInfoName C.cl_kernel_exec_info

const (
	// KernelExecInfoSvmPtrs defines the SVM pointers that must reference locations contained entirely within buffers
	// that are passed to kernel as arguments, or that are passed through the execution information.
	//
	// Non-argument SVM buffers must be specified by passing pointers to those buffers via SetKernelExecInfo()
	// for coarse-grain and fine-grain buffer SVM allocations but not for fine-grain system SVM allocations.
	//
	// Required type: []unsafe.Pointer
	// Since: 2.0
	KernelExecInfoSvmPtrs KernelExecInfoName = C.CL_KERNEL_EXEC_INFO_SVM_PTRS
	// KernelExecInfoSvmFineGrainSystem indicates whether the kernel uses pointers that are fine-grain system SVM
	// allocations. These fine-grain system SVM pointers may be passed as arguments or defined in SVM buffers that
	// are passed as arguments to kernel.
	//
	// Required type: Bool
	// Since: 2.0
	KernelExecInfoSvmFineGrainSystem KernelExecInfoName = C.CL_KERNEL_EXEC_INFO_SVM_FINE_GRAIN_SYSTEM
)

// SetKernelExecInfo passes additional information other than argument values to a kernel.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetKernelExecInfo.html
func SetKernelExecInfo(kernel Kernel, paramName KernelExecInfoName, paramSize uintptr, paramValue unsafe.Pointer) error {
	status := C.clSetKernelExecInfo(
		kernel.handle(),
		C.cl_kernel_exec_info(paramName),
		C.size_t(paramSize),
		paramValue)
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// KernelInfoName identifies properties of a kernel, which can be queried with KernelInfo().
type KernelInfoName C.cl_kernel_info

const (
	// KernelFunctionNameInfo returns the kernel function name.
	//
	// Returned type: string
	KernelFunctionNameInfo KernelInfoName = C.CL_KERNEL_FUNCTION_NAME
	// KernelNumArgsInfo returns the number of arguments to kernel.
	//
	// Returned type: Uint
	KernelNumArgsInfo KernelInfoName = C.CL_KERNEL_NUM_ARGS
	// KernelReferenceCountInfo returns the kernel reference count.
	//
	// Note: The reference count returned should be considered immediately stale. It is unsuitable for general
	// use in applications. This feature is provided for identifying memory leaks.
	//
	// Returned type: Uint
	KernelReferenceCountInfo KernelInfoName = C.CL_KERNEL_REFERENCE_COUNT
	// KernelContextInfo returns the context associated with kernel.
	//
	// Returned type: Context
	KernelContextInfo KernelInfoName = C.CL_KERNEL_CONTEXT
	// KernelProgramInfo returns the program object associated with kernel.
	//
	// Returned type: Program
	KernelProgramInfo KernelInfoName = C.CL_KERNEL_PROGRAM
	// KernelAttributesInfo returns any attributes specified using the __attribute__ OpenCL C qualifier
	// (or using an OpenCL C++ qualifier syntax [[]] ) with the kernel function declaration in the program source.
	//
	// Returned type: string
	// Since: 1.2
	KernelAttributesInfo KernelInfoName = C.CL_KERNEL_ATTRIBUTES
)

// KernelInfo returns information about the kernel object.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use KernelInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetKernelInfo.html
func KernelInfo(kernel Kernel, paramName KernelInfoName, paramSize uint, paramValue unsafe.Pointer) (uint, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetKernelInfo(
		kernel.handle(),
		C.cl_kernel_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uint(sizeReturn), nil
}

// KernelInfoString is a convenience method for KernelInfo() to query information values that are
// string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func KernelInfoString(kernel Kernel, paramName KernelInfoName) (string, error) {
	return queryString(func(paramSize uint, paramValue unsafe.Pointer) (uint, error) {
		return KernelInfo(kernel, paramName, paramSize, paramValue)
	})
}

// KernelWorkGroupInfoName identifies properties of a kernel work group, which can be queried with KernelWorkGroupInfo().
type KernelWorkGroupInfoName C.cl_kernel_work_group_info

const (
	// KernelWorkGroupSizeInfo provides a mechanism for the application to query the maximum work-group size
	// that can be used to execute the kernel on a specific device given by device.
	// The OpenCL implementation uses the resource requirements of the kernel (register usage etc.) to determine what
	// this work-group size should be.
	//
	// As a result and unlike DeviceMaxWorkGroupSizeInfo this value may vary from one kernel to another as well as
	// one device to another.
	//
	// KernelCompileWorkGroupSizeInfo will be less than or equal to DeviceMaxWorkGroupSizeInfo for a given kernel object.
	//
	// Returned type: uintptr
	KernelWorkGroupSizeInfo KernelWorkGroupInfoName = C.CL_KERNEL_WORK_GROUP_SIZE
	// KernelCompileWorkGroupSizeInfo returns the work-group size specified in the kernel source or IL.
	//
	// If the work-group size is not specified in the kernel source or IL, (0, 0, 0) is returned.
	//
	// Returned type: [3]uintptr
	KernelCompileWorkGroupSizeInfo KernelWorkGroupInfoName = C.CL_KERNEL_COMPILE_WORK_GROUP_SIZE
	// KernelLocalMemSizeInfo returns the amount of local memory in bytes being used by a kernel.
	// This includes local memory that may be needed by an implementation to execute the kernel, variables declared
	// inside the kernel with the __local address qualifier and local memory to be allocated for arguments to the
	// kernel declared as pointers with the __local address qualifier and whose size is specified with SetKernelArg().
	//
	// If the local memory size, for any pointer argument to the kernel declared with the __local address qualifier,
	// is not specified, its size is assumed to be 0.
	//
	// Returned type: uint64
	KernelLocalMemSizeInfo KernelWorkGroupInfoName = C.CL_KERNEL_LOCAL_MEM_SIZE
	// KernelPreferredWorkGroupSizeMultipleInfo returns the preferred multiple of work-group size for launch.
	// This is a performance hint. Specifying a work-group size that is not a multiple of the value returned by this
	// query as the value of the local work size argument to EnqueueNDRangeKernel() will not fail to enqueue the kernel
	// for execution unless the work-group size specified is larger than the device maximum.
	//
	// Returned type: uintptr
	KernelPreferredWorkGroupSizeMultipleInfo KernelWorkGroupInfoName = C.CL_KERNEL_PREFERRED_WORK_GROUP_SIZE_MULTIPLE
	// KernelPrivateMemSizeInfo returns the minimum amount of private memory, in bytes, used by each work-item in
	// the kernel. This value may include any private memory needed by an implementation to execute the kernel,
	// including that used by the language built-ins and variable declared inside the kernel with the __private qualifier.
	//
	// Returned type: uint64
	KernelPrivateMemSizeInfo KernelWorkGroupInfoName = C.CL_KERNEL_PRIVATE_MEM_SIZE
	// KernelGlobalWorkSizeInfo provides a mechanism for the application to query the maximum global size that can be
	// used to execute a kernel (i.e. globalWorkSize argument to EnqueueNDRangeKernel()) on a custom device given by
	// device or a built-in kernel on an OpenCL device given by device.
	//
	// If device is not a custom device and kernel is not a built-in kernel, GetKernelWorkGroupInfo() returns the
	// error ErrInvalidValue.
	//
	// Returned type: [3]uintptr
	// Since: 1.2
	KernelGlobalWorkSizeInfo KernelWorkGroupInfoName = C.CL_KERNEL_GLOBAL_WORK_SIZE
)

// KernelWorkGroupInfo returns information about the kernel object that may be specific to a device.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetKernelWorkGroupInfo.html
func KernelWorkGroupInfo(kernel Kernel, device DeviceID, paramName KernelWorkGroupInfoName, paramSize uint, paramValue unsafe.Pointer) (uint, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetKernelWorkGroupInfo(
		kernel.handle(),
		device.handle(),
		C.cl_kernel_work_group_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uint(sizeReturn), nil
}

// KernelArgInfoName identifies properties of a kernel argument, which can be queried with KernelArgInfo().
type KernelArgInfoName C.cl_kernel_arg_info

const (
	// KernelArgAddressQualifierInfo returns the address qualifier specified for the argument.
	//
	// Returned type: KernelArgAddressQualifier
	// Since: 1.2
	KernelArgAddressQualifierInfo KernelArgInfoName = C.CL_KERNEL_ARG_ADDRESS_QUALIFIER
	// KernelArgAccessQualifierInfo returns the access qualifier specified for the argument.
	//
	// Returned type: KernelArgAccessQualifier
	// Since: 1.2
	KernelArgAccessQualifierInfo KernelArgInfoName = C.CL_KERNEL_ARG_ACCESS_QUALIFIER
	// KernelArgTypeNameInfo returns the type name specified for the argument.
	// The type name returned will be the argument type name as it was declared with any whitespace removed.
	// If argument type name is an unsigned scalar type (i.e. unsigned char, unsigned short, unsigned int,
	// unsigned long), uchar, ushort, uint and ulong will be returned.
	// The argument type name returned does not include any type qualifiers.
	//
	// Returned type: string
	// Since: 1.2
	KernelArgTypeNameInfo KernelArgInfoName = C.CL_KERNEL_ARG_TYPE_NAME
	// KernelArgTypeQualifierInfo returns a bitfield describing one or more type qualifiers specified for the argument.
	//
	// Returned type: KernelArgTypeQualifier
	// Since: 1.2
	KernelArgTypeQualifierInfo KernelArgInfoName = C.CL_KERNEL_ARG_TYPE_QUALIFIER
	// KernelArgNameInfo returns the name specified for the argument.
	//
	// Returned type: string
	// Since: 1.2
	KernelArgNameInfo KernelArgInfoName = C.CL_KERNEL_ARG_NAME
)

// KernelArgAddressQualifier describes the address qualifier for a kernel argument.
type KernelArgAddressQualifier C.cl_kernel_arg_address_qualifier

// List of possible KernelArgAddressQualifier values.
const (
	KernelArgAddressGlobal   KernelArgAddressQualifier = C.CL_KERNEL_ARG_ADDRESS_GLOBAL
	KernelArgAddressLocal    KernelArgAddressQualifier = C.CL_KERNEL_ARG_ADDRESS_LOCAL
	KernelArgAddressConstant KernelArgAddressQualifier = C.CL_KERNEL_ARG_ADDRESS_CONSTANT
	KernelArgAddressPrivate  KernelArgAddressQualifier = C.CL_KERNEL_ARG_ADDRESS_PRIVATE
)

// KernelArgAccessQualifier describes the access qualifier for a kernel argument.
type KernelArgAccessQualifier C.cl_kernel_arg_access_qualifier

// List of possible KernelArgAccessQualifier values.
const (
	KernelArgAccessReadOnly  KernelArgAccessQualifier = C.CL_KERNEL_ARG_ACCESS_READ_ONLY
	KernelArgAccessWriteOnly KernelArgAccessQualifier = C.CL_KERNEL_ARG_ACCESS_WRITE_ONLY
	KernelArgAccessReadWrite KernelArgAccessQualifier = C.CL_KERNEL_ARG_ACCESS_READ_WRITE
	KernelArgAccessNone      KernelArgAccessQualifier = C.CL_KERNEL_ARG_ACCESS_NONE
)

// KernelArgTypeQualifier describes the type for a kernel argument.
type KernelArgTypeQualifier C.cl_kernel_arg_type_qualifier

// List of possible KernelArgTypeQualifier values.
const (
	KernelArgTypeNone     KernelArgTypeQualifier = C.CL_KERNEL_ARG_TYPE_NONE
	KernelArgTypeConst    KernelArgTypeQualifier = C.CL_KERNEL_ARG_TYPE_CONST
	KernelArgTypeRestrict KernelArgTypeQualifier = C.CL_KERNEL_ARG_TYPE_RESTRICT
	KernelArgTypeVolatile KernelArgTypeQualifier = C.CL_KERNEL_ARG_TYPE_VOLATILE
	KernelArgTypePipe     KernelArgTypeQualifier = C.CL_KERNEL_ARG_TYPE_PIPE
)

// KernelArgInfo returns information about the arguments of a kernel.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.For convenience, use KernelArgInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetKernelArgInfo.html
func KernelArgInfo(kernel Kernel, index uint32, paramName KernelArgInfoName, paramSize uint, paramValue unsafe.Pointer) (uint, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetKernelArgInfo(
		kernel.handle(),
		C.cl_uint(index),
		C.cl_kernel_work_group_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uint(sizeReturn), nil
}

// KernelArgInfoString is a convenience method for KernelArgInfo() to query information values that are
// string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func KernelArgInfoString(kernel Kernel, index uint32, paramName KernelArgInfoName) (string, error) {
	return queryString(func(paramSize uint, paramValue unsafe.Pointer) (uint, error) {
		return KernelArgInfo(kernel, index, paramName, paramSize, paramValue)
	})
}

// KernelSubGroupInfoName identifies properties of a kernel, which can be queried with KernelSubGroupInfo().
type KernelSubGroupInfoName C.cl_kernel_sub_group_info

const (
	// KernelMaxSubGroupSizeForNDRangeInfo returns the maximum sub-group size for this kernel.
	// All sub-groups must be the same size, while the last subgroup in any work-group (i.e. the subgroup with
	// the maximum index) could be the same or smaller size.
	//
	// The input value must be an array of uintptr values corresponding to the local work size parameter of the
	// intended dispatch. The number of dimensions in the ND-range will be inferred from the value specified for
	// input size.
	//
	// Input type: []uintptr
	// Returned type: uintptr
	// Since: 2.1
	// Extension: cl_khr_subgroups
	KernelMaxSubGroupSizeForNDRangeInfo KernelSubGroupInfoName = C.CL_KERNEL_MAX_SUB_GROUP_SIZE_FOR_NDRANGE
	// KernelSubGroupCountForNDRangeInfo returns the number of sub-groups that will be present in each work-group for
	// a given local work size. All workgroups, apart from the last work-group in each dimension in the presence of
	// non-uniform work-group sizes, will have the same number of sub-groups.
	//
	// The input value must be an array of uintptr values corresponding to the local work size parameter of the
	// intended dispatch. The number of dimensions in the ND-range will be inferred from the value specified
	// for input size.
	//
	// Input type: []uintptr
	// Returned type: uintptr
	// Since: 2.1
	// Extension: cl_khr_subgroups
	KernelSubGroupCountForNDRangeInfo KernelSubGroupInfoName = C.CL_KERNEL_SUB_GROUP_COUNT_FOR_NDRANGE
	// KernelLocalSizeForSubGroupCountInfo returns the local size that will generate the requested number of
	// sub-groups for the kernel. The output array must be an array of uintptr values corresponding to the local
	// size parameter. Any returned work-group will have one dimension. Other dimensions inferred from the value
	// specified for param_value_size will be filled with the value 1. The returned value will produce an exact number
	// of sub-groups and result in no partial groups for an executing kernel except in the case where the last
	// work-group in a dimension has a size different from that of the other groups.
	// If no work-group size can accommodate the requested number of sub-groups, 0 will be returned in each element
	// of the return array.
	//
	// Input type: uintptr
	// Returned type: []uintptr
	// Since: 2.1
	// Extension: cl_khr_subgroups
	KernelLocalSizeForSubGroupCountInfo KernelSubGroupInfoName = C.CL_KERNEL_LOCAL_SIZE_FOR_SUB_GROUP_COUNT
	// KernelMaxNumSubGroupsInfo provides a mechanism for the application to query the maximum number of sub-groups
	// that may make up each work-group to execute a kernel on a specific device given by device.
	// The OpenCL implementation uses the resource requirements of the kernel (register usage etc.) to determine
	// what this work-group size should be. The returned value may be used to compute a work-group size to enqueue
	// the kernel with to give a round number of sub-groups for an enqueue.
	//
	// Input type: (ignored)
	// Returned type: uintptr
	// Since: 2.1
	// Extension: cl_khr_subgroups
	KernelMaxNumSubGroupsInfo KernelSubGroupInfoName = C.CL_KERNEL_MAX_NUM_SUB_GROUPS
	// KernelCompileNumSubGroupsInfo returns the number of sub-groups per work-group specified in the kernel source
	// or IL. If the sub-group count is not specified then 0 is returned.
	//
	// Input type: (ignored)
	// Returned type: uintptr
	// Since: 2.1
	// Extension: cl_khr_subgroups
	KernelCompileNumSubGroupsInfo KernelSubGroupInfoName = C.CL_KERNEL_COMPILE_NUM_SUB_GROUPS
)

// KernelSubGroupInfo returns information about the kernel object.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.
//
// Since: 2.1
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetKernelSubGroupInfo.html
func KernelSubGroupInfo(kernel Kernel, device DeviceID, paramName KernelSubGroupInfoName,
	inputSize uint, inputValue unsafe.Pointer,
	paramSize uint, paramValue unsafe.Pointer) (uint, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetKernelSubGroupInfo(
		kernel.handle(),
		device.handle(),
		C.cl_kernel_sub_group_info(paramName),
		C.size_t(inputSize),
		inputValue,
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uint(sizeReturn), nil
}

// WorkDimension describes the parameters within one dimension of a work group.
type WorkDimension struct {
	GlobalOffset uintptr
	GlobalSize   uintptr
	LocalSize    uintptr
}

// EnqueueNDRangeKernel enqueues a command to execute a kernel on a device.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueNDRangeKernel.html
func EnqueueNDRangeKernel(commandQueue CommandQueue, kernel Kernel, workDimensions []WorkDimension, waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	globalWorkOffsets := make([]uintptr, len(workDimensions))
	globalWorkSizes := make([]uintptr, len(workDimensions))
	localWorkSizes := make([]uintptr, len(workDimensions))
	for i, dimension := range workDimensions {
		globalWorkOffsets[i] = dimension.GlobalOffset
		globalWorkSizes[i] = dimension.GlobalSize
		localWorkSizes[i] = dimension.LocalSize
	}
	status := C.clEnqueueNDRangeKernel(
		commandQueue.handle(),
		kernel.handle(),
		C.cl_uint(len(workDimensions)),
		(*C.size_t)(unsafe.Pointer(&globalWorkOffsets[0])),
		(*C.size_t)(unsafe.Pointer(&globalWorkSizes[0])),
		(*C.size_t)(unsafe.Pointer(&localWorkSizes[0])),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueNativeKernel enqueues a command to execute a native Go function not compiled using the OpenCL compiler.
//
// The provided callback function will receive pointers to global memory that represents the provided MemObject
// entries.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueNativeKernel.html
func EnqueueNativeKernel(commandQueue CommandQueue, callback func([]unsafe.Pointer), memObjects []MemObject, waitList []Event, event *Event) error {
	callbackUserData, err := userDataFor(func(argBasePtr unsafe.Pointer) {
		argMovePtr := argBasePtr
		memPtr := make([]unsafe.Pointer, len(memObjects))
		for i := 0; i < len(memObjects); i++ {
			memPtr[i] = unsafe.Pointer(*(**uintptr)(argMovePtr))
			argMovePtr = unsafe.Add(argMovePtr, unsafe.Sizeof(uintptr(0)))
		}
		callback(memPtr)
	})
	if err != nil {
		return err
	}
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	rawArgs := make([]uintptr, len(memObjects)+1)
	rawArgs[0] = uintptr(unsafe.Pointer(callbackUserData.ptr))
	var rawArgsMemLocs []uintptr
	var rawArgsPtr unsafe.Pointer
	var rawMemObjectsPtr unsafe.Pointer
	var rawArgsMemLocsPtr unsafe.Pointer
	if len(memObjects) > 0 {
		rawMemObjectsPtr = unsafe.Pointer(&memObjects[0])
		rawArgsMemLocs = make([]uintptr, len(memObjects))
		for i := 0; i < len(memObjects); i++ {
			rawArgsMemLocs[i] = uintptr(unsafe.Pointer(&rawArgs[1+i]))
		}
		rawArgsMemLocsPtr = unsafe.Pointer(&rawArgsMemLocs[0])
	}
	rawArgsPtr = unsafe.Pointer(&rawArgs[0])
	status := C.cl30EnqueueNativeKernel(
		commandQueue.handle(),
		rawArgsPtr,
		C.size_t(uintptr(len(rawArgs))*unsafe.Sizeof(uintptr(0))),
		C.cl_uint(len(memObjects)),
		(*C.cl_mem)(rawMemObjectsPtr),
		rawArgsMemLocsPtr,
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return StatusError(status)
	}
	return nil
}

//export cl30GoKernelNativeCallback
func cl30GoKernelNativeCallback(args unsafe.Pointer) {
	callbackUserData := userDataFrom(*(**C.uintptr_t)(args))
	callback := callbackUserData.Value().(func(unsafe.Pointer))
	callbackUserData.Delete()
	callback(unsafe.Add(args, unsafe.Sizeof(uintptr(0))))
}
