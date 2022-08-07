package cl30

// #include "api.h"
import "C"
import "fmt"

// StatusError represents an error based on a status value from an OpenCL call.
type StatusError C.cl_int

// Error returns the string presentation of the numeric value.
// A name lookup is not performed as errors can be extended through extensions, making a consistent presentation
// difficult.
func (err StatusError) Error() string {
	return fmt.Sprintf("%d", int(err))
}

// This block contains common error constants.
const (
	ErrDeviceNotFound                     StatusError = C.CL_DEVICE_NOT_FOUND
	ErrDeviceNotAvailable                 StatusError = C.CL_DEVICE_NOT_AVAILABLE
	ErrCompilerNotAvailable               StatusError = C.CL_COMPILER_NOT_AVAILABLE
	ErrMemObjectAllocationFailure         StatusError = C.CL_MEM_OBJECT_ALLOCATION_FAILURE
	ErrOutOfResources                     StatusError = C.CL_OUT_OF_RESOURCES
	ErrOutOfHostMemory                    StatusError = C.CL_OUT_OF_HOST_MEMORY
	ErrProfilingInfoNotAvailable          StatusError = C.CL_PROFILING_INFO_NOT_AVAILABLE
	ErrMemCopyOverlap                     StatusError = C.CL_MEM_COPY_OVERLAP
	ErrImageFormatMismatch                StatusError = C.CL_IMAGE_FORMAT_MISMATCH
	ErrImageFormatNotSupported            StatusError = C.CL_IMAGE_FORMAT_NOT_SUPPORTED
	ErrBuildProgramFailure                StatusError = C.CL_BUILD_PROGRAM_FAILURE
	ErrMapFailure                         StatusError = C.CL_MAP_FAILURE
	ErrMisalignedSubBufferOffset          StatusError = C.CL_MISALIGNED_SUB_BUFFER_OFFSET
	ErrExecStatusErrorForEventsInWaitList StatusError = C.CL_EXEC_STATUS_ERROR_FOR_EVENTS_IN_WAIT_LIST
	ErrCompileProgramFailure              StatusError = C.CL_COMPILE_PROGRAM_FAILURE
	ErrLinkerNotAvailable                 StatusError = C.CL_LINKER_NOT_AVAILABLE
	ErrLinkProgramFailure                 StatusError = C.CL_LINK_PROGRAM_FAILURE
	ErrDevicePartitionFailed              StatusError = C.CL_DEVICE_PARTITION_FAILED
	ErrKernelArgInfoNotAvailable          StatusError = C.CL_KERNEL_ARG_INFO_NOT_AVAILABLE
	ErrInvalidValue                       StatusError = C.CL_INVALID_VALUE
	ErrInvalidDeviceType                  StatusError = C.CL_INVALID_DEVICE_TYPE
	ErrInvalidPlatform                    StatusError = C.CL_INVALID_PLATFORM
	ErrInvalidDevice                      StatusError = C.CL_INVALID_DEVICE
	ErrInvalidContext                     StatusError = C.CL_INVALID_CONTEXT
	ErrInvalidQueueProperties             StatusError = C.CL_INVALID_QUEUE_PROPERTIES
	ErrInvalidCommandQueue                StatusError = C.CL_INVALID_COMMAND_QUEUE
	ErrInvalidHostPtr                     StatusError = C.CL_INVALID_HOST_PTR
	ErrInvalidMemObject                   StatusError = C.CL_INVALID_MEM_OBJECT
	ErrINvalidImageFormatDescriptor       StatusError = C.CL_INVALID_IMAGE_FORMAT_DESCRIPTOR
	ErrInvalidImageSize                   StatusError = C.CL_INVALID_IMAGE_SIZE
	ErrInvalidSampler                     StatusError = C.CL_INVALID_SAMPLER
	ErrInvalidBinary                      StatusError = C.CL_INVALID_BINARY
	ErrInvalidBuildOptions                StatusError = C.CL_INVALID_BUILD_OPTIONS
	ErrInvalidProgram                     StatusError = C.CL_INVALID_PROGRAM
	ErrInvalidProgramExecutable           StatusError = C.CL_INVALID_PROGRAM_EXECUTABLE
	ErrInvalidKernelName                  StatusError = C.CL_INVALID_KERNEL_NAME
	ErrInvalidKernelDefinition            StatusError = C.CL_INVALID_KERNEL_DEFINITION
	ErrInvalidKernel                      StatusError = C.CL_INVALID_KERNEL
	ErrInvalidArgIndex                    StatusError = C.CL_INVALID_ARG_INDEX
	ErrInvalidArgValue                    StatusError = C.CL_INVALID_ARG_VALUE
	ErrInvalidArgSize                     StatusError = C.CL_INVALID_ARG_SIZE
	ErrInvalidKernelArgs                  StatusError = C.CL_INVALID_KERNEL_ARGS
	ErrInvalidWorkDimension               StatusError = C.CL_INVALID_WORK_DIMENSION
	ErrInvalidWorkGroupSize               StatusError = C.CL_INVALID_WORK_GROUP_SIZE
	ErrInvalidWorkItemSize                StatusError = C.CL_INVALID_WORK_ITEM_SIZE
	ErrInvalidGlobalOffset                StatusError = C.CL_INVALID_GLOBAL_OFFSET
	ErrInvalidEventWaitList               StatusError = C.CL_INVALID_EVENT_WAIT_LIST
	ErrInvalidEvent                       StatusError = C.CL_INVALID_EVENT
	ErrInvalidOperation                   StatusError = C.CL_INVALID_OPERATION
	ErrInvalidGlObject                    StatusError = C.CL_INVALID_GL_OBJECT
	ErrInvalidBufferSize                  StatusError = C.CL_INVALID_BUFFER_SIZE
	ErrInvalidMipLevel                    StatusError = C.CL_INVALID_MIP_LEVEL
	ErrInvalidGlobalWorkSize              StatusError = C.CL_INVALID_GLOBAL_WORK_SIZE
	ErrInvalidProperty                    StatusError = C.CL_INVALID_PROPERTY
	ErrInvalidImageDescriptor             StatusError = C.CL_INVALID_IMAGE_DESCRIPTOR
	ErrInvalidCompilerOptions             StatusError = C.CL_INVALID_COMPILER_OPTIONS
	ErrInvalidLinkerOptions               StatusError = C.CL_INVALID_LINKER_OPTIONS
	ErrInvalidDevicePartitionCount        StatusError = C.CL_INVALID_DEVICE_PARTITION_COUNT
	ErrInvalidPipeSize                    StatusError = C.CL_INVALID_PIPE_SIZE
	ErrInvalidDeviceQueue                 StatusError = C.CL_INVALID_DEVICE_QUEUE
	ErrInvalidSpecID                      StatusError = C.CL_INVALID_SPEC_ID
	ErrMaxSizeRestrictionExceeded         StatusError = C.CL_MAX_SIZE_RESTRICTION_EXCEEDED
)

// WrapperError represents a basic error that occurs within the wrapper.
type WrapperError string

// String returns the error text.
func (err WrapperError) Error() string {
	return string(err)
}

const (
	// ErrExtensionNotAvailable is returned in case an extension loader function could not complete its setup.
	ErrExtensionNotAvailable WrapperError = "extension not available"
	// ErrExtensionNotLoaded is returned in case an extension was used without being loaded.
	ErrExtensionNotLoaded WrapperError = "extension not loaded"
	// ErrDataSizeLimitExceeded is returned by convenience functions that query information, typically string-loading
	// functions. The error occurs if the required size would be beyond a (feasible) limit.
	ErrDataSizeLimitExceeded WrapperError = "data size limit exceeded"
	// ErrOutOfMemory is returned by wrapper functions that need to allocate memory.
	ErrOutOfMemory WrapperError = "out of memory"
)
