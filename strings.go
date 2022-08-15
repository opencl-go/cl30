package cl30

// #include "api.h"
import "C"
import (
	"unsafe"
)

// queryString extracts a string with the help of a load function.
// The load function shall return the required number of bytes for the string, including the terminating NUL byte.
// The load function is called twice, once with zero/nil to query the needed size, then a second time to retrieve
// the value.
func queryString(load func(paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error)) (string, error) {
	requiredSize, err := load(0, nil)
	if err != nil {
		return "", err
	}
	if requiredSize > 1024*1024*10 {
		return "", ErrDataSizeLimitExceeded
	}
	if requiredSize == 0 {
		return "", nil
	}
	raw := C.calloc(C.size_t(requiredSize), 1)
	if raw == nil {
		return "", ErrOutOfMemory
	}
	defer C.free(raw)
	returnedSize, err := load(requiredSize, raw)
	if err != nil {
		return "", err
	}
	// The returned size may be different from the originally reported value. Avoid using more than possible.
	if returnedSize > requiredSize {
		returnedSize = requiredSize
	}
	// The returned size should include the terminating 0x00 byte. Avoid an underflow when skipping it in next step.
	if returnedSize == 0 {
		returnedSize = 1
	}
	usedSize := returnedSize - 1
	return C.GoStringN((*C.char)(raw), C.int(usedSize)), nil
}
