package cl30

// #include "api.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// PipeProperty is one entry of properties which are taken into account when creating pipes.
type PipeProperty []uintptr

// CreatePipe creates a pipe object.
//
// For the flags parameter, only MemReadWriteFlag and MemHostNoAccessFlag can be specified when creating a pipe object.
// If the value specified for flags is 0, the default is used which is MemReadWriteFlag | MemHostNoAccessFlag.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreatePipe.html
func CreatePipe(context Context, flags MemFlags, packetSize, maxPackets uint32, properties ...PipeProperty) (MemObject, error) {
	var rawPropertyList []uintptr
	for _, property := range properties {
		rawPropertyList = append(rawPropertyList, property...)
	}
	var rawProperties unsafe.Pointer
	if len(properties) > 0 {
		rawPropertyList = append(rawPropertyList, 0)
		rawProperties = unsafe.Pointer(&rawPropertyList[0])
	}
	var status C.cl_int
	pipe := C.clCreatePipe(
		context.handle(),
		C.cl_mem_flags(flags),
		C.cl_uint(packetSize),
		C.cl_uint(maxPackets),
		(*C.cl_pipe_properties)(rawProperties),
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return MemObject(*((*uintptr)(unsafe.Pointer(&pipe)))), nil
}

// PipeInfoName identifies properties of a pipe, which can be queried with PipeInfo().
type PipeInfoName C.cl_pipe_info

const (
	// PipePacketSizeInfo returns pipe packet size specified when pipe is created with CreatePipe().
	//
	// Returned type: uint32
	// Since: 2.0
	PipePacketSizeInfo PipeInfoName = C.CL_PIPE_PACKET_SIZE
	// PipeMaxPacketsInfo returns maximum number of packets specified when pipe is created with CreatePipe().
	//
	// Returned type: uint32
	// Since: 2.0
	PipeMaxPacketsInfo PipeInfoName = C.CL_PIPE_MAX_PACKETS
	// PipePropertiesInfo returns the properties argument specified in CreatePipe().
	//
	// Returned type: []uintptr
	// Since: 3.0
	PipePropertiesInfo PipeInfoName = C.CL_PIPE_PROPERTIES
)

// PipeInfo queries information specific to a pipe object.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetPipeInfo.html
func PipeInfo(pipe MemObject, paramName PipeInfoName, param HostMemory) (uintptr, error) {
	sizeReturn := C.size_t(0)
	paramPtr := ResolvePointer(param, false, "param")
	status := C.clGetPipeInfo(
		pipe.handle(),
		C.cl_pipe_info(paramName),
		sizeOf(param),
		paramPtr,
		&sizeReturn)
	runtime.KeepAlive(paramPtr)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}
