package cl30

import (
	"sync"
	"unsafe"
)

// #include "api.h"
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

//export cl30GoContextDestructorCallback
func cl30GoContextDestructorCallback(_ Context, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func())
	callbackUserData.Delete()
	callback()
}
