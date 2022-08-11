package cl30

import (
	"unsafe"
)

// #include "api.h"
// extern cl_int cl30ExtTerminateContextKHR(void *fn, cl_context context);
import "C"

// ExtensionTerminateContextKhr represents the functionality provided by the "cl_khr_terminate_context" extension.
// Load the extension with LoadExtensionTerminateContextKhr().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/cl_khr_terminate_context.html
// Extension: KhrTerminateContextExtensionName
type ExtensionTerminateContextKhr struct {
	clTerminateContextKhr unsafe.Pointer
}

// LoadExtensionTerminateContextKhr loads the required functions for the extension and returns an instance
// to ExtensionTerminateContextKhr if possible.
//
// Extension: KhrTerminateContextExtensionName
func LoadExtensionTerminateContextKhr(id PlatformID) (*ExtensionTerminateContextKhr, error) {
	clTerminateContextKhr := ExtensionFunctionAddressForPlatform(id, "clTerminateContextKHR")
	if clTerminateContextKhr == nil {
		return nil, ErrExtensionNotAvailable
	}
	return &ExtensionTerminateContextKhr{clTerminateContextKhr: clTerminateContextKhr}, nil
}

// TerminateContext terminates all pending work associated with the context and renders all data owned by the context
// invalid. It is the responsibility of the application to release all objects associated with the context being
// terminated.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clTerminateContextKHR.html
// Extension: KhrTerminateContextExtensionName
func (ext *ExtensionTerminateContextKhr) TerminateContext(context Context) error {
	if (ext == nil) || (ext.clTerminateContextKhr == nil) {
		return ErrExtensionNotLoaded
	}
	status := C.cl30ExtTerminateContextKHR(ext.clTerminateContextKhr, context.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

const (
	// KhrTerminateContextExtensionName is the official name of the extension
	// handled by ExtensionTerminateContextKhr.
	KhrTerminateContextExtensionName = "cl_khr_terminate_context"

	// ErrContextTerminatedKhr is returned for all operations within a context after the context has been terminated.
	//
	// Extension: KhrTerminateContextExtensionName
	ErrContextTerminatedKhr StatusError = C.CL_CONTEXT_TERMINATED_KHR

	// DeviceTerminateCapabilityKhrInfo represents the capability about a device to which degree it supports
	// termination of their contexts.
	//
	// Info value type: DeviceTerminateCapabilityKhrFlags
	// Extension: KhrTerminateContextExtensionName
	DeviceTerminateCapabilityKhrInfo DeviceInfoName = C.CL_DEVICE_TERMINATE_CAPABILITY_KHR

	// ContextTerminateKhrProperty represents a context property that requests to enable context termination.
	//
	// Use WithTermination() for convenience.
	//
	// Property value type: Bool
	// Extension: KhrTerminateContextExtensionName
	ContextTerminateKhrProperty uintptr = C.CL_CONTEXT_TERMINATE_KHR
)

// DeviceTerminateCapabilityKhrFlags describe the termination capability of the OpenCL device.
//
// Extension: KhrTerminateContextExtensionName
type DeviceTerminateCapabilityKhrFlags C.cl_bitfield

const (
	// DeviceTerminateCapabilityKhrContext indicates that context termination is supported.
	//
	// Note: The constant is an assumption. Refer to https://github.com/KhronosGroup/OpenCL-Docs/issues/813 for details.
	//
	// Extension: KhrTerminateContextExtensionName
	DeviceTerminateCapabilityKhrContext DeviceTerminateCapabilityKhrFlags = 1 << 0
)

// WithTermination is a convenience function to create a valid ContextTerminateKhrProperty.
// Use it in combination with CreateContext() or CreateContextFromType().
//
// Extension: KhrTerminateContextExtensionName
func WithTermination(enabled bool) ContextProperty {
	return ContextProperty{ContextTerminateKhrProperty, uintptr(BoolFrom(enabled))}
}
