package cl30

// #include <stdlib.h>
import "C"
import "unsafe"

// HostPointer identifies a pointer in host-space.
type HostPointer interface {
	// Pointer returns the raw pointer value.
	Pointer() unsafe.Pointer
}

// StaticHostPointer is a marker interface to identify pointers that will not be moved by the runtime.
type StaticHostPointer interface {
	HostPointer
	// IsStatic marks that the value of the pointer will not change during runtime.
	IsStatic()
}

// ResolvePointer is used to retrieve the unsafe pointer value of a host pointer.
// In case a static pointer value is required and the given pointer does not cast to StaticHostPointer,
// the function panics.
func ResolvePointer(ptr HostPointer, requireStatic bool, paramName string) unsafe.Pointer {
	if ptr == nil {
		return nil
	}
	raw := ptr.Pointer()
	if raw == nil {
		return nil
	}
	if !requireStatic {
		return raw
	}
	if _, isStatic := ptr.(StaticHostPointer); !isStatic {
		// If you receive a panic from this line, it means a pointer that potentially can change its value
		// was used in a function call which requires the value to be static. This is typically the case for
		// functions that run asynchronously while keeping the raw pointer value in C-land. It is possible that
		// the Go runtime moves values around in memory, and their addresses can change. Using such a pointer
		// may end up in crashes at unexpected and arbitrary times.
		panic("A pointer was provided that is not explicitly static. Parameter=" + paramName)
	}
	return raw
}

// HostMemory represents a range of bytes in the host-accessible memory.
type HostMemory interface {
	HostPointer
	// Size returns the number of bytes this range represents.
	Size() int
}

// Null returns a HostMemory instance that represents the null-pointer.
// It resolves to nil, has a size of zero, and is fixed.
func Null() HostMemory {
	var null *FixedHostMemory
	return null
}

func sizeOf(mem HostMemory) C.size_t {
	if mem == nil {
		return 0
	}
	return C.size_t(mem.Size())
}

// HostMemoryBytes returns a slice wrapped over the low-level pointer. Use this function for convenience.
func HostMemoryBytes(mem HostMemory) []byte {
	if mem == nil {
		return nil
	}
	return unsafe.Slice((*byte)(mem.Pointer()), mem.Size())
}

// FixedHostMemory is a memory range "fixed" at an address.
// It is created by using the low-level allocation function, which will not change its address during runtime.
type FixedHostMemory struct {
	raw  unsafe.Pointer
	size int
}

// AllocFixedHostMemory allocates a fixed memory range of given size, in bytes.
// Call Free() when you no longer need the memory block.
func AllocFixedHostMemory(size int) *FixedHostMemory {
	raw := C.calloc((C.size_t)(size), 1)
	return &FixedHostMemory{
		raw:  raw,
		size: size,
	}
}

// Free releases the underlying memory buffer. Call this function to avoid memory leaks, and call it only
// after no more references to the raw pointer are in use.
func (mem *FixedHostMemory) Free() {
	if (mem == nil) || (mem.raw == nil) {
		return
	}
	C.free(mem.raw)
	mem.size = 0
	mem.raw = nil
}

// Size returns the number of bytes this range represents.
func (mem *FixedHostMemory) Size() int {
	if mem == nil {
		return 0
	}
	return mem.size
}

// Pointer returns the address to the first byte of the range.
func (mem *FixedHostMemory) Pointer() unsafe.Pointer {
	if mem == nil {
		return nil
	}
	return mem.raw
}

// IsStatic marks that this allocated memory has a static pointer that will not change during runtime.
func (mem *FixedHostMemory) IsStatic() {}

type runtimeHostMemory struct {
	ptr  unsafe.Pointer
	size int
}

// Size returns the number of bytes this range represents.
func (mem *runtimeHostMemory) Size() int {
	if mem == nil {
		return 0
	}
	return mem.size
}

// Pointer returns the address to the first byte of the range.
func (mem *runtimeHostMemory) Pointer() unsafe.Pointer {
	if mem == nil {
		return nil
	}
	return mem.ptr
}

// HostValueOf returns a HostMemory instance that represents a copy of the given value.
// Use this to pass in a copy of a primitive or a structure to a function taking a HostMemory.
func HostValueOf[T any](v T) HostMemory {
	return HostReferenceOf(&v)
}

// HostReferenceOf returns a HostMemory instance that represents the memory location of the given value.
// Use this to pass in a pointer of a Go type to a function, where the call will return into this given type.
func HostReferenceOf[T any](v *T) HostMemory {
	if v == nil {
		return Null()
	}
	return &runtimeHostMemory{
		ptr:  unsafe.Pointer(v),
		size: int(unsafe.Sizeof(*v)),
	}
}

// HostVectorOf returns a HostMemory instance that represents the memory location of the given slice.
// Use this to pass in a pointer of a Go type to a function, where the call will return into this given slice.
func HostVectorOf[T any](v []T) HostMemory {
	if len(v) == 0 {
		return Null()
	}
	return &runtimeHostMemory{
		ptr:  unsafe.Pointer(&v[0]),
		size: int(unsafe.Sizeof(v[0])) * len(v),
	}
}
