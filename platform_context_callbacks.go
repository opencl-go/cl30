package cl30

import (
	"runtime/cgo"
	"sync"
	"sync/atomic"
	"unsafe"
)

// #include "api.h"
// extern CL_CALLBACK void cl30CContextErrorCallback(char const *errorInfo,
//	                                                 void const *privateInfoPtr, size_t privateInfoLen,
//                                                   intptr_t userData);
// extern CL_CALLBACK void cl30CContextDestructorCallback(cl_context context, void *userData);
import "C"

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
	key     uintptr
	handler ContextErrorHandler
}

// ErrContextErrorCallbackKeySpaceExhausted is returned from NewContextErrorCallback() in case
// no more callback instances can be registered.
const ErrContextErrorCallbackKeySpaceExhausted WrapperError = "key space exhausted"

// NewContextErrorCallback creates and registers a new callback.
//
// As this is a globally registered resource, registration may fail if the internal lookup key space is
// exhausted. This key space is based on uintptr, so it varies by system when this is exhausted.
//
// The provided handler can be called from other threads from within the OpenCL runtime.
func NewContextErrorCallback(handler ContextErrorHandler) (*ContextErrorCallback, error) {
	key := atomic.AddUintptr(&contextErrorCallbackCounter, 1)
	if key == 0 {
		return nil, ErrContextErrorCallbackKeySpaceExhausted
	}
	cb := &ContextErrorCallback{
		key:     key,
		handler: handler,
	}
	contextErrorCallbackMutex.Lock()
	defer contextErrorCallbackMutex.Unlock()
	contextErrorCallbacksByKey[key] = cb
	return cb, nil
}

// Release removes the registered callback from the system. When this function returns, the assigned
// handler will no longer be called.
func (cb *ContextErrorCallback) Release() {
	if (cb == nil) || (cb.key == 0) {
		return
	}
	contextErrorCallbackMutex.Lock()
	defer contextErrorCallbackMutex.Unlock()
	delete(contextErrorCallbacksByKey, cb.key)
	cb.key = 0
}

var (
	contextErrorCallbackCounter uintptr
	contextErrorCallbackMutex   = sync.RWMutex{}
	contextErrorCallbacksByKey  = map[uintptr]*ContextErrorCallback{}
)

func cContextErrorCallbackFunc() unsafe.Pointer {
	return C.cl30CContextErrorCallback
}

//export cl30GoContextErrorCallback
func cl30GoContextErrorCallback(errorInfo *C.char, privateInfoPtr *C.uint8_t, privateInfoLen C.size_t, key C.intptr_t) {
	contextErrorCallbackMutex.RLock()
	defer contextErrorCallbackMutex.RUnlock()
	cb, known := contextErrorCallbacksByKey[uintptr(key)]
	if !known {
		return
	}
	privateInfo := unsafe.Slice((*byte)(privateInfoPtr), uintptr(privateInfoLen))
	cb.handler.Handle(C.GoString(errorInfo), privateInfo)
}

func cContextDestructorCallbackFunc() unsafe.Pointer {
	return C.cl30CContextDestructorCallback
}

//export cl30GoContextDestructorCallback
func cl30GoContextDestructorCallback(context Context, userData unsafe.Pointer) {
	handle := *(*cgo.Handle)(userData)
	callback := handle.Value().(func())
	handle.Delete()
	callback()
}
