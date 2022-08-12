package cl30

// #cgo CFLAGS: -DCL_USE_DEPRECATED_OPENCL_1_2_APIS
// #cgo CXXFLAGS: -DCL_USE_DEPRECATED_OPENCL_1_2_APIS
// #cgo CPPFLAGS: -DCL_USE_DEPRECATED_OPENCL_1_2_APIS
// #include "api.h"
import "C"
import "unsafe"

// CreateCommandQueue creates a command-queue on a specific device.
//
// Deprecated: 1.2; Use CreateCommandQueueWithProperties() instead.
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateCommandQueue.html
func CreateCommandQueue(context Context, deviceID DeviceID, properties CommandQueuePropertiesFlags) (CommandQueue, error) {
	var status C.cl_int
	commandQueue := C.clCreateCommandQueue(
		context.handle(),
		deviceID.handle(),
		C.cl_command_queue_properties(properties),
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return CommandQueue(*((*uintptr)(unsafe.Pointer(&commandQueue)))), nil
}
