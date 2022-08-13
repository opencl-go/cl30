package cl30

// #include "api.h"
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// userData contains a dynamically allocated memory that holds a cgo.Handle .
// This type is necessary to safely transport the value of a handle across.
//
// The example for cgo.Handle showcases to transport the pointer to the handle value - yet the value
// the pointer points to can (and will) go out of scope for callbacks for which this userData is intended.
// As the handle value itself is not allowed to be directly cast into an unsafe.Pointer, the value needs
// to be stored in an allocated memory block. Incidentally, this memory block has then a pointer that can
// be used as actual C-land user data.
type userData struct {
	ptr *C.uintptr_t
}

func userDataFor(v any) (userData, error) {
	ptr := (*C.uintptr_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.uintptr_t(0)))))
	if ptr == nil {
		return userData{}, ErrOutOfMemory
	}
	h := cgo.NewHandle(v)
	*ptr = C.uintptr_t(h)
	return userData{ptr: ptr}, nil
}

func userDataFrom(ptr *C.uintptr_t) userData {
	return userData{ptr: ptr}
}

func (data userData) Value() any {
	h := cgo.Handle(*data.ptr)
	return h.Value()
}

func (data userData) Delete() {
	if data.ptr == nil {
		return
	}
	h := cgo.Handle(*data.ptr)
	h.Delete()
	C.free(unsafe.Pointer(data.ptr))
	data.ptr = nil
}
