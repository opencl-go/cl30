package cl30

// #include "api.h"
// extern cl_int cl30EnqueueSVMFree(cl_command_queue commandQueue,
//    cl_uint svmPointerCount, void *svmPointers,
//    uintptr_t *userData,
//    cl_uint waitListCount, cl_event const *waitList,
//    cl_event *event);
// extern cl_int cl30EnqueueSVMMigrateMem(
//    cl_command_queue commandQueue,
//    cl_uint svmPointerCount, void *svmPointers,
//    size_t *sizes, cl_mem_migration_flags flags,
//    cl_uint waitListCount, cl_event const *waitList,
//    cl_event *event);
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// SvmMemFlags describe properties of a shared virtual memory (SVM) buffer.
type SvmMemFlags C.cl_mem_flags

// SvmBuffer is a shared virtual memory (SVM) buffer.
type SvmBuffer struct {
	size int
	ptr  unsafe.Pointer
}

// Size returns the number of bytes the buffer contains.
func (buf *SvmBuffer) Size() int {
	return buf.size
}

// Pointer returns the raw pointer value to the buffer. This value is only usable for host copy operations
// while the buffer is mapped.
func (buf *SvmBuffer) Pointer() unsafe.Pointer {
	return buf.ptr
}

// IsStatic is a marker to indicate the pointer is a static value.
func (buf *SvmBuffer) IsStatic() {}

// SvmAlloc allocates a shared virtual memory (SVM) buffer that can be shared by the host and all devices in an OpenCL
// context that support shared virtual memory.
//
// For flags, potential values are MemReadWriteFlag, MemWriteOnlyFlag, MemReadOnlyFlag, MemSvmAtomicsFlag,
// MemSvmFineGrainBufferFlag.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSVMAlloc.html
func SvmAlloc(context Context, flags SvmMemFlags, size int, alignment uint32) (*SvmBuffer, error) {
	ptr := C.clSVMAlloc(
		context.handle(),
		C.cl_svm_mem_flags(flags),
		C.size_t(size),
		C.cl_uint(alignment))
	if ptr == nil {
		return nil, ErrOutOfMemory
	}
	buf := &SvmBuffer{
		size: size,
		ptr:  ptr,
	}
	return buf, nil
}

// SvmFree frees a shared virtual memory buffer allocated using SvmAlloc().
//
// SvmFree does not wait for previously enqueued commands that may be using ptr to finish before freeing the memory.
// It is the responsibility of the application to make sure that enqueued commands that use ptr have finished before
// freeing the memory.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSVMFree.html
func SvmFree(context Context, ptr *SvmBuffer) {
	C.clSVMFree(context.handle(), ptr.ptr)
}

// EnqueueSvmFree enqueues a command to free shared virtual memory allocated using SvmAlloc() or a shared system
// memory pointer.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueSVMFree.html
func EnqueueSvmFree(commandQueue CommandQueue, ptrs []HostPointer, callback func(CommandQueue, []unsafe.Pointer), waitList []Event, event *Event) error {
	var callbackUserData userData
	if callback != nil {
		var err error
		callbackUserData, err = userDataFor(callback)
		if err != nil {
			return err
		}
	}
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	rawPointers := make([]unsafe.Pointer, len(ptrs))
	ptrAddresses := make([]uintptr, len(ptrs))
	for i, ptr := range ptrs {
		rawPointers[i] = ResolvePointer(ptr, true, fmt.Sprintf("ptrs[%d]", i))
		ptrAddresses[i] = uintptr(rawPointers[i])
	}
	status := C.cl30EnqueueSVMFree(
		commandQueue.handle(),
		C.cl_uint(len(ptrs)),
		unsafe.Pointer(&ptrAddresses[0]),
		callbackUserData.ptr,
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	for _, rawPtr := range rawPointers {
		runtime.KeepAlive(rawPtr)
	}
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

//export cl30GoSvmFreeCallback
func cl30GoSvmFreeCallback(commandQueue CommandQueue, svmPointerCount C.cl_uint, svmPointers unsafe.Pointer, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func(CommandQueue, []unsafe.Pointer))
	callbackUserData.Delete()
	ptrs := unsafe.Slice((*unsafe.Pointer)(svmPointers), int(svmPointerCount))
	callback(commandQueue, ptrs)
}

// EnqueueSvmMemcpy enqueues a command to do a memcpy operation.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueSVMMemcpy.html
func EnqueueSvmMemcpy(commandQueue CommandQueue, blocking bool, dstPtr, srcPtr HostPointer, size int,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	rawDstPtr := ResolvePointer(dstPtr, !blocking, "dstPtr")
	rawSrcPtr := ResolvePointer(srcPtr, !blocking, "srcPtr")
	status := C.clEnqueueSVMMemcpy(
		commandQueue.handle(),
		C.cl_bool(BoolFrom(blocking)),
		rawDstPtr,
		rawSrcPtr,
		C.size_t(size),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	runtime.KeepAlive(rawDstPtr)
	runtime.KeepAlive(rawSrcPtr)
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueSvmMemFill enqueues a command to fill a region in memory with a pattern of a given pattern size.
//
// The pattern must be a scalar or vector integer or floating-point data type supported by OpenCL.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueSVMMemFill.html
func EnqueueSvmMemFill(commandQueue CommandQueue, svmPtr, pattern HostPointer, patternSize, size int,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	rawSvmPtr := ResolvePointer(svmPtr, true, "svmPtr")
	patternPtr := ResolvePointer(pattern, true, "pattern")
	status := C.clEnqueueSVMMemFill(
		commandQueue.handle(),
		rawSvmPtr,
		patternPtr,
		C.size_t(patternSize),
		C.size_t(size),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	runtime.KeepAlive(rawSvmPtr)
	runtime.KeepAlive(patternPtr)
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueSvmMap enqueues a command that will allow the host to update a region of an SVM buffer.
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueSVMMap.html
func EnqueueSvmMap(commandQueue CommandQueue, blocking bool, flags MemFlags, svmPtr *SvmBuffer, size int,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueSVMMap(
		commandQueue.handle(),
		C.cl_bool(BoolFrom(blocking)),
		C.cl_map_flags(flags),
		svmPtr.ptr,
		C.size_t(size),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueSvmUnmap enqueues a command to indicate that the host has completed updating the region given by an SVM
// pointer and which was specified in a previous call to EnqueueSvmMap().
//
// Since: 2.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueSVMUnmap.html
func EnqueueSvmUnmap(commandQueue CommandQueue, svmPtr *SvmBuffer, waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueSVMUnmap(
		commandQueue.handle(),
		svmPtr.ptr,
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueSvmMigrateMem enqueues a command to indicate which device a set of ranges of SVM allocations should be
// associated with.
//
// Since: 2.1
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueSVMMigrateMem.html
func EnqueueSvmMigrateMem(commandQueue CommandQueue, svmPtrs []*SvmBuffer, sizes []int, flags MemMigrationFlags,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	svmPtrAddresses := make([]uintptr, len(svmPtrs))
	for i, svmPtr := range svmPtrs {
		svmPtrAddresses[i] = uintptr(svmPtr.ptr)
	}
	var sizesPtr unsafe.Pointer
	if len(sizes) > 0 {
		sizesPtr = unsafe.Pointer(&sizes[0])
	}
	status := C.cl30EnqueueSVMMigrateMem(
		commandQueue.handle(),
		C.cl_uint(len(svmPtrs)),
		unsafe.Pointer(&svmPtrAddresses[0]),
		(*C.size_t)(sizesPtr),
		C.cl_mem_migration_flags(flags),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}
