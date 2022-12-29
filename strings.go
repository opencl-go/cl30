package cl30

// #include "api.h"
import "C"

// queryString extracts a string with the help of a load function.
// The load function shall return the required number of bytes for the string, including the terminating NUL byte.
// The load function is called twice, once with zero/nil to query the needed size, then a second time to retrieve
// the value.
func queryString(load func(param HostMemory) (uintptr, error)) (string, error) {
	requiredSize, err := load(Null())
	if err != nil {
		return "", err
	}
	if requiredSize > 1024*1024*10 {
		return "", ErrDataSizeLimitExceeded
	}
	if requiredSize == 0 {
		return "", nil
	}
	raw := AllocFixedHostMemory(int(requiredSize))
	if raw == nil {
		return "", ErrOutOfMemory
	}
	defer raw.Free()
	returnedSize, err := load(raw)
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
	return C.GoStringN((*C.char)(raw.Pointer()), C.int(usedSize)), nil
}
