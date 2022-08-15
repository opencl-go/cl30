package cl30

// #include "api.h"
// extern cl_context cl30CreateContext(cl_context_properties *properties,
//     cl_uint numDevices, cl_device_id *devices,
//     uintptr_t *userData,
//     cl_int *errcodeReturn);
// extern cl_context cl30CreateContextFromType(cl_context_properties *properties,
//     cl_device_type deviceType,
//     uintptr_t *userData,
//     cl_int *errcodeReturn);
// extern cl_int cl30SetContextDestructorCallback(cl_context context, uintptr_t *userData);
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// A Context is used by the OpenCL runtime for managing objects such as command-queues, memory,
// program and kernel objects, and for executing kernels on one or more devices specified in the context.
type Context uintptr

func (c Context) handle() C.cl_context {
	return *((*C.cl_context)(unsafe.Pointer(&c)))
}

// String provides a readable presentation of the platform identifier.
// It is based on the numerical value of the underlying pointer.
func (c Context) String() string {
	return fmt.Sprintf("0x%X", uintptr(c))
}

const (
	// ContextPlatformProperty specifies the platform to use.
	//
	// Use OnPlatform() for convenience.
	//
	// Property value type: PlatformID
	ContextPlatformProperty uintptr = C.CL_CONTEXT_PLATFORM
	// ContextInteropUserSyncProperty specifies whether the user is responsible for synchronization between OpenCL and
	// other APIs. Please refer to the specific sections in the OpenCL Extension Specification that describe sharing
	// with other APIs for restrictions on using this flag.
	//
	// Use WithInteropUserSync() for convenience.
	//
	// Property value type: Bool
	// Since: 1.2
	ContextInteropUserSyncProperty uintptr = C.CL_CONTEXT_INTEROP_USER_SYNC
)

// ContextProperty is one entry of properties which are taken into account when creating context objects.
type ContextProperty []uintptr

// OnPlatform is a convenience function to create a valid ContextPlatformProperty.
// Use it in combination with CreateContext() or CreateContextFromType().
func OnPlatform(id PlatformID) ContextProperty {
	return ContextProperty{ContextPlatformProperty, uintptr(id)}
}

// WithInteropUserSync is a convenience function to create a valid ContextInteropUserSyncProperty.
// Use it in combination with CreateContext() or CreateContextFromType().
func WithInteropUserSync(value bool) ContextProperty {
	return ContextProperty{ContextInteropUserSyncProperty, uintptr(BoolFrom(value))}
}

// CreateContext creates an OpenCL context for the specified devices.
//
// The callback is an optional receiver for any errors that happen during creation or execution of the context.
// It is possible that one registered callback is re-used for multiple context objects. Synchronization and separation
// between the contexts is left to the user.
// It is possible to call ContextErrorCallback.Release() while any context the callback is registered with is still
// active.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateContext.html
func CreateContext(deviceIds []DeviceID, callback *ContextErrorCallback, properties ...ContextProperty) (Context, error) {
	var rawPropertyList []uintptr
	for _, property := range properties {
		rawPropertyList = append(rawPropertyList, property...)
	}
	var rawProperties unsafe.Pointer
	if len(properties) > 0 {
		rawPropertyList = append(rawPropertyList, 0)
		rawProperties = unsafe.Pointer(&rawPropertyList[0])
	}
	var rawDeviceIds unsafe.Pointer
	if len(deviceIds) > 0 {
		rawDeviceIds = unsafe.Pointer(&deviceIds[0])
	}
	callbackKey := (*C.uintptr_t)(nil)
	if callback != nil {
		callbackKey = callback.userData.ptr
	}
	var status C.cl_int
	context := C.cl30CreateContext(
		(*C.cl_context_properties)(rawProperties),
		C.uint(len(deviceIds)),
		(*C.cl_device_id)(rawDeviceIds),
		callbackKey,
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Context(*((*uintptr)(unsafe.Pointer(&context)))), nil
}

// CreateContextFromType creates an OpenCL context for devices that match the given device type.
// The context does not reference any sub-devices that may have been created from these devices.
//
// The callback is an optional receiver for any errors that happen during creation or execution of the context.
// It is possible that one registered callback is re-used for multiple context objects. Synchronization and separation
// between the contexts is left to the user.
// It is possible to call ContextErrorCallback.Release() while any context the callback is registered with is still
// active.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateContextFromType.html
func CreateContextFromType(deviceType DeviceTypeFlags, callback *ContextErrorCallback, properties ...ContextProperty) (Context, error) {
	var rawPropertyList []uintptr
	for _, property := range properties {
		rawPropertyList = append(rawPropertyList, property...)
	}
	var rawProperties unsafe.Pointer
	if len(properties) > 0 {
		rawPropertyList = append(rawPropertyList, 0)
		rawProperties = unsafe.Pointer(&rawPropertyList[0])
	}
	callbackKey := (*C.uintptr_t)(nil)
	if callback != nil {
		callbackKey = callback.userData.ptr
	}
	var status C.cl_int
	context := C.cl30CreateContextFromType(
		(*C.cl_context_properties)(rawProperties),
		C.cl_device_type(deviceType),
		callbackKey,
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Context(*((*uintptr)(unsafe.Pointer(&context)))), nil
}

// ContextErrorHandler is informed about an error that occurred within the processing of a context.
type ContextErrorHandler interface {
	// Handle receives the information on the event. The private information is an opaque detail, specific
	// to the event, which may help during further analysis.
	Handle(errorInfo string, privateInfo []byte)
}

// ContextErrorHandlerFunc is a convenience type for ContextErrorHandler. This function type implements
// the interface and forwards the call to itself.
type ContextErrorHandlerFunc func(errorInfo string, privateInfo []byte)

// Handle calls the function itself.
func (handler ContextErrorHandlerFunc) Handle(errorInfo string, privateInfo []byte) {
	handler(errorInfo, privateInfo)
}

// ContextErrorCallback is a registered callback that can be used to receive error messages from contexts.
// Create and register a new callback with NewContextErrorCallback().
// The callback is a globally registered resource that must be released with Release() when it is no longer needed.
type ContextErrorCallback struct {
	userData userData
	handler  ContextErrorHandler
}

// NewContextErrorCallback creates and registers a new callback.
//
// As this is a globally registered resource, registration may fail if memory is exhausted.
//
// The provided handler can be called from other threads from within the OpenCL runtime.
func NewContextErrorCallback(handler ContextErrorHandler) (*ContextErrorCallback, error) {
	handlerUserData, err := userDataFor(handler)
	if err != nil {
		return nil, err
	}
	cb := &ContextErrorCallback{
		userData: handlerUserData,
		handler:  handler,
	}
	contextErrorCallbackMutex.Lock()
	defer contextErrorCallbackMutex.Unlock()
	contextErrorCallbacksByPtr[handlerUserData.ptr] = cb
	return cb, nil
}

// Release removes the registered callback from the system. When this function returns, the assigned
// handler will no longer be called.
//
// In case you release the error callback before the associated context is destroyed,
// there is a slight chance for a later, newly created error callback to be called for that older context.
// This can happen if the allocated low-level memory block that holds the Go handle receives the same pointer as the
// previous callback had.
func (cb *ContextErrorCallback) Release() {
	if (cb == nil) || (cb.userData.ptr == nil) {
		return
	}
	contextErrorCallbackMutex.Lock()
	defer contextErrorCallbackMutex.Unlock()
	delete(contextErrorCallbacksByPtr, cb.userData.ptr)
	cb.userData.Delete()
}

var (
	contextErrorCallbackMutex  = sync.RWMutex{}
	contextErrorCallbacksByPtr = map[*C.uintptr_t]*ContextErrorCallback{}
)

//export cl30GoContextErrorCallback
func cl30GoContextErrorCallback(errorInfo *C.char, privateInfoPtr *C.uint8_t, privateInfoLen C.size_t, key *C.uintptr_t) {
	contextErrorCallbackMutex.RLock()
	defer contextErrorCallbackMutex.RUnlock()
	cb, known := contextErrorCallbacksByPtr[key]
	if !known {
		return
	}
	privateInfo := unsafe.Slice((*byte)(privateInfoPtr), uintptr(privateInfoLen))
	cb.handler.Handle(C.GoString(errorInfo), privateInfo)
}

// RetainContext increments the context reference count.
//
// CreateContext() and CreateContextFromType() perform an implicit retain. This is very helpful for 3rd party
// libraries, which typically get a context passed to them by the application. However, it is possible that the
// application may delete the context without informing the library. Allowing functions to attach to (i.e. retain)
// and release a context solves the problem of a context being used by a library no longer being valid.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clRetainContext.html
func RetainContext(context Context) error {
	status := C.clRetainContext(context.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ReleaseContext decrements the context reference count.
//
// After the reference count becomes zero and all the objects attached to context (such as memory objects,
// command-queues) are released, the context is deleted.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clReleaseContext.html
func ReleaseContext(context Context) error {
	status := C.clReleaseContext(context.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ContextInfoName identifies properties of a context, which can be queried with ContextInfo().
type ContextInfoName C.cl_context_info

const (
	// ContextReferenceCountInfo returns the context reference count.
	//
	// Note: The reference count returned should be considered immediately stale. It is unsuitable for
	// general use in applications. This feature is provided for identifying memory leaks.
	//
	// Returned type: Uint
	ContextReferenceCountInfo ContextInfoName = C.CL_CONTEXT_REFERENCE_COUNT
	// ContextDevicesInfo returns the list of devices and sub-devices in context.
	//
	// Returned type: []DeviceID
	ContextDevicesInfo ContextInfoName = C.CL_CONTEXT_DEVICES
	// ContextNumDevicesInfo returns the number of devices in context.
	//
	// Returned type: Uint
	// Since: 1.1
	ContextNumDevicesInfo ContextInfoName = C.CL_CONTEXT_NUM_DEVICES
	// ContextPropertiesInfo returns the properties argument specified in CreateContext() or CreateContextFromType().
	//
	// The returned list is the concatenated list of all the properties provided at creation.
	// This list also includes the terminating zero value.
	//
	// Returned type: []uintptr
	ContextPropertiesInfo ContextInfoName = C.CL_CONTEXT_PROPERTIES
)

// ContextInfo queries information about a context.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use ContextInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetContextInfo.html
func ContextInfo(context Context, paramName ContextInfoName, paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetContextInfo(
		context.handle(),
		C.cl_context_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// ContextInfoString is a convenience method for ContextInfo() to query information values that are string-based.
//
// This function does not verify the queried information is indeed of type string. It assumes the information is
// a NUL terminated raw string and will extract the bytes as characters before that.
func ContextInfoString(context Context, paramName ContextInfoName) (string, error) {
	return queryString(func(paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
		return ContextInfo(context, paramName, paramSize, paramValue)
	})
}

// SetContextDestructorCallback registers a destructor callback function with a context.
//
// Each call to SetContextDestructorCallback() registers the specified callback function on a destructor callback
// stack associated with context.
// The registered callback functions are called in the reverse order in which they were registered.
//
// If a context callback function was specified when context was created, it will not be called after any
// context destructor callback is called. Therefore, the context destructor callback provides a mechanism for
// an application to safely re-use or release the context callback function when context was created.
//
// Since: 3.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetContextDestructorCallback.html
func SetContextDestructorCallback(context Context, callback func()) error {
	callbackUserData, err := userDataFor(callback)
	if err != nil {
		return err
	}
	status := C.cl30SetContextDestructorCallback(context.handle(), callbackUserData.ptr)
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return StatusError(status)
	}
	return nil
}

//export cl30GoContextDestructorCallback
func cl30GoContextDestructorCallback(_ Context, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func())
	callbackUserData.Delete()
	callback()
}
