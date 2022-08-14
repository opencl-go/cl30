package cl30

// #cgo CFLAGS: -DCL_USE_DEPRECATED_OPENCL_1_2_APIS -DCL_USE_DEPRECATED_OPENCL_2_2_APIS
// #cgo CXXFLAGS: -DCL_USE_DEPRECATED_OPENCL_1_2_APIS -DCL_USE_DEPRECATED_OPENCL_2_2_APIS
// #cgo CPPFLAGS: -DCL_USE_DEPRECATED_OPENCL_1_2_APIS -DCL_USE_DEPRECATED_OPENCL_2_2_APIS
// #include "api.h"
// extern cl_int cl30SetProgramReleaseCallback(cl_program program, uintptr_t *userData);
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

// CreateSampler creates a sampler object.
//
// Deprecated: 1.2; Use CreateSamplerWithProperties() instead.
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateSampler.html
func CreateSampler(context Context, normalizedCoords bool, addressingMode SamplerAddressingMode, filterMode SamplerFilterMode) (Sampler, error) {
	var status C.cl_int
	sampler := C.clCreateSampler(
		context.handle(),
		C.cl_bool(BoolFrom(normalizedCoords)),
		C.cl_addressing_mode(addressingMode),
		C.cl_filter_mode(filterMode),
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Sampler(*((*uintptr)(unsafe.Pointer(&sampler)))), nil
}

// SetProgramReleaseCallback registers a destructor callback function with a program object.
//
// Each call to SetProgramReleaseCallback() registers the specified callback function on a callback stack associated
// with program. The registered callback functions are called in the reverse order in which they were registered.
// The registered callback functions are called after destructors (if any) for program scope global variables (if any)
// are called and before the program object is deleted.
// This provides a mechanism for an application to be notified when destructors for program scope global variables
// are complete.
//
// SetProgramReleaseCallback() may unconditionally return an error if no devices in the context associated with
// program support destructors for program scope global variables.
// Support for constructors and destructors for program scope global variables is required only for OpenCL 2.2 devices.
//
// Since: 2.2
// Deprecated: 2.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetProgramReleaseCallback.html
func SetProgramReleaseCallback(program Program, callback func()) error {
	callbackUserData, err := userDataFor(callback)
	if err != nil {
		return err
	}
	status := C.cl30SetProgramReleaseCallback(program.handle(), callbackUserData.ptr)
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return StatusError(status)
	}
	return nil
}

//export cl30GoProgramReleaseCallback
func cl30GoProgramReleaseCallback(_ Program, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func())
	callbackUserData.Delete()
	callback()
}
