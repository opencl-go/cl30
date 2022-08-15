package cl30

// #include "api.h"
// extern cl_int cl30SetEventCallback(cl_event event, cl_int callbackType, uintptr_t *userData);
import "C"
import (
	"fmt"
	"unsafe"
)

// Event objects are used as synchronization points between different commands within a context.
// Enqueue* functions offer the option to return a new event object, and a manually controlled event object can
// be created with CreateUserEvent().
type Event uintptr

func (event Event) handle() C.cl_event {
	return *(*C.cl_event)(unsafe.Pointer(&event))
}

// String provides a readable presentation of the event identifier.
// It is based on the numerical value of the underlying pointer.
func (event Event) String() string {
	return fmt.Sprintf("0x%X", uintptr(event))
}

// CreateUserEvent creates a user event object.
// User events allow applications to enqueue commands that wait on a user event to finish before the command is
// executed by the device.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateUserEvent.html
func CreateUserEvent(context Context) (Event, error) {
	var status C.cl_int
	event := C.clCreateUserEvent(context.handle(), &status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return Event(*((*uintptr)(unsafe.Pointer(&event)))), nil
}

// SetUserEventStatus sets the execution status of a user event object.
//
// If there are enqueued commands with user events in the eventWaitList argument of Enqueue* commands, the user
// must ensure that the status of these user events being waited on are set using SetUserEventStatus() before any
// OpenCL APIs that release OpenCL objects except for event objects are called; otherwise the behavior is undefined.
//
// The provided executionStatus can either be EventCommandCompleteStatus or a negative value to indicate an error.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetUserEventStatus.html
func SetUserEventStatus(event Event, executionStatus int) error {
	status := C.clSetUserEventStatus(event.handle(), C.cl_int(executionStatus))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// WaitForEvents waits on the host thread for commands identified by event objects to complete.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clWaitForEvents.html
func WaitForEvents(events []Event) error {
	var rawEvents unsafe.Pointer
	if len(events) > 0 {
		rawEvents = unsafe.Pointer(&events[0])
	}
	status := C.clWaitForEvents(C.cl_uint(len(events)), (*C.cl_event)(rawEvents))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EventInfoName identifies properties of an event, which can be queried with EventInfo().
type EventInfoName C.cl_event_info

const (
	// EventCommandQueueInfo returns the command-queue associated with event. For user event objects,
	// a zero value is returned.
	//
	// Returned type: CommandQueue
	EventCommandQueueInfo EventInfoName = C.CL_EVENT_COMMAND_QUEUE
	// EventContextInfo returns the context associated with event.
	//
	// Returned type: Context
	EventContextInfo EventInfoName = C.CL_EVENT_CONTEXT
	// EventCommandTypeInfo return the command type associated with the event.
	//
	// Returned type: EventCommandType
	EventCommandTypeInfo EventInfoName = C.CL_EVENT_COMMAND_TYPE
	// EventReferenceCountInfo returns the event reference count.
	//
	// Note: The reference count returned should be considered immediately stale. It is unsuitable for
	// general use in applications. This feature is provided for identifying memory leaks.
	//
	// Returned type: uint32
	EventReferenceCountInfo EventInfoName = C.CL_EVENT_REFERENCE_COUNT
	// EventCommandExecutionStatusInfo returns the execution status of the command identified by event.
	//
	// Error codes are identified by a negative integer value: The command was abnormally terminated - this may be
	// caused by a bad memory access for example. These error codes come from the same set of error codes that are
	// returned from the platform or runtime API calls as status codes.
	//
	// Returned type: EventCommandExecutionStatus
	EventCommandExecutionStatusInfo EventInfoName = C.CL_EVENT_COMMAND_EXECUTION_STATUS
)

// EventCommandType identifies the associated command with an event.
type EventCommandType C.cl_uint

const (
	// CommandNdRangeKernel events are created by EnqueueNDRangeKernel().
	CommandNdRangeKernel EventCommandType = C.CL_COMMAND_NDRANGE_KERNEL
	// CommandTask events are created by EnqueueTask().
	CommandTask EventCommandType = C.CL_COMMAND_TASK
	// CommandNativeKernel events are created by EnqueueNativeKernel().
	CommandNativeKernel EventCommandType = C.CL_COMMAND_NATIVE_KERNEL
	// CommandReadBuffer events are created by EnqueueReadBuffer().
	CommandReadBuffer EventCommandType = C.CL_COMMAND_READ_BUFFER
	// CommandWriteBuffer events are created by EnqueueWriteBuffer().
	CommandWriteBuffer EventCommandType = C.CL_COMMAND_WRITE_BUFFER
	// CommandCopyBuffer events are created by EnqueueCopyBuffer().
	CommandCopyBuffer EventCommandType = C.CL_COMMAND_COPY_BUFFER
	// CommandReadImage events are created by EnqueueReadImage().
	CommandReadImage EventCommandType = C.CL_COMMAND_READ_IMAGE
	// CommandWriteImage events are created by EnqueueWriteImage().
	CommandWriteImage EventCommandType = C.CL_COMMAND_WRITE_IMAGE
	// CommandCopyImage events are created by EnqueueCopyImage().
	CommandCopyImage EventCommandType = C.CL_COMMAND_COPY_IMAGE
	// CommandCopyImageToBuffer events are created by EnqueueCopyImageToBuffer().
	CommandCopyImageToBuffer EventCommandType = C.CL_COMMAND_COPY_IMAGE_TO_BUFFER
	// CommandCopyBufferToImage events are created by EnqueueCopyBufferToImage().
	CommandCopyBufferToImage EventCommandType = C.CL_COMMAND_COPY_BUFFER_TO_IMAGE
	// CommandMapBuffer events are created by EnqueueMapBuffer().
	CommandMapBuffer EventCommandType = C.CL_COMMAND_MAP_BUFFER
	// CommandMapImage events are created by EnqueueMapImage().
	CommandMapImage EventCommandType = C.CL_COMMAND_MAP_IMAGE
	// CommandUnmapMemObject events are created by EnqueueUnmapMemObject().
	CommandUnmapMemObject EventCommandType = C.CL_COMMAND_UNMAP_MEM_OBJECT
	// CommandMarker events are created by EnqueueMarker() and EnqueueMarkerWithWaitList().
	CommandMarker EventCommandType = C.CL_COMMAND_MARKER
	// CommandReadBufferRect events are created by EnqueueReadBufferRect().
	//
	// Since: 1.1
	CommandReadBufferRect EventCommandType = C.CL_COMMAND_READ_BUFFER_RECT
	// CommandWriteBufferRect events are created by EnqueueWriteBufferRect().
	//
	// Since: 1.1
	CommandWriteBufferRect EventCommandType = C.CL_COMMAND_WRITE_BUFFER_RECT
	// CommandCopyBufferRect events are created by EnqueueCopyBufferRect().
	//
	// Since: 1.1
	CommandCopyBufferRect EventCommandType = C.CL_COMMAND_COPY_BUFFER_RECT
	// CommandUser events are created by CreateUserEvent().
	//
	// Since: 1.1
	CommandUser EventCommandType = C.CL_COMMAND_USER
	// CommandBarrier events are created by EnqueueBarrier() and EnqueueBarrierWithWaitList().
	//
	// Since: 1.2
	CommandBarrier EventCommandType = C.CL_COMMAND_BARRIER
	// CommandMigrateMemObjects events are created by EnqueueMigrateMemObjects().
	//
	// Since: 1.2
	CommandMigrateMemObjects EventCommandType = C.CL_COMMAND_MIGRATE_MEM_OBJECTS
	// CommandFillBuffer events are created by EnqueueFillBuffer().
	//
	// Since: 1.2
	CommandFillBuffer EventCommandType = C.CL_COMMAND_FILL_BUFFER
	// CommandFillImage events are created by EnqueueFillImage().
	//
	// Since: 1.2
	CommandFillImage EventCommandType = C.CL_COMMAND_FILL_IMAGE

	// CommandSvmFree events are created by EnqueueSvmFree().
	//
	// Since: 2.0
	CommandSvmFree EventCommandType = C.CL_COMMAND_SVM_FREE
	// CommandSvmMemcpy events are created by EnqueueSvmMemcpy().
	//
	// Since: 2.0
	CommandSvmMemcpy EventCommandType = C.CL_COMMAND_SVM_MEMCPY
	// CommandSvmMemFill events are created by EnqueueSvmMemFill().
	//
	// Since: 2.0
	CommandSvmMemFill EventCommandType = C.CL_COMMAND_SVM_MEMFILL
	// CommandSvmMap events are created by EnqueueSvmMap().
	//
	// Since: 2.0
	CommandSvmMap EventCommandType = C.CL_COMMAND_SVM_MAP
	// CommandSvmUnmap events are created by EnqueueSvmUnmap().
	//
	// Since: 2.0
	CommandSvmUnmap EventCommandType = C.CL_COMMAND_SVM_UNMAP

	// CommandSvmMigrateMem events are created by EnqueueSvmMigrateMem().
	//
	// Since: 3.0
	CommandSvmMigrateMem EventCommandType = C.CL_COMMAND_SVM_MIGRATE_MEM
)

// EventCommandExecutionStatus describes the execution status of an event.
// Negative values are error status values.
type EventCommandExecutionStatus C.cl_int

const (
	// EventCommandQueuedStatus means the command has been enqueued in the command-queue.
	EventCommandQueuedStatus EventCommandExecutionStatus = C.CL_QUEUED
	// EventCommandSubmittedStatus means the enqueued command has been submitted by the host to the device associated
	// with the command-queue.
	EventCommandSubmittedStatus EventCommandExecutionStatus = C.CL_SUBMITTED
	// EventCommandRunningStatus means the device is currently executing this command.
	EventCommandRunningStatus EventCommandExecutionStatus = C.CL_RUNNING
	// EventCommandCompleteStatus means that the command has completed.
	EventCommandCompleteStatus EventCommandExecutionStatus = C.CL_COMPLETE
)

// EventInfo queries information about an event.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetEventInfo.html
func EventInfo(event Event, paramName EventInfoName, paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetEventInfo(
		event.handle(),
		C.cl_event_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// RetainEvent increments the event reference count.
// The OpenCL commands that return an event perform an implicit retain.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clRetainEvent.html
func RetainEvent(event Event) error {
	status := C.clRetainEvent(event.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ReleaseEvent decrements the event reference count.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clReleaseEvent.html
func ReleaseEvent(event Event) error {
	status := C.clReleaseEvent(event.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EventProfilingInfoName identifies profiling properties of an event, which can be queried with EventProfilingInfo().
type EventProfilingInfoName C.cl_profiling_info

const (
	// ProfilingCommandQueuedInfo describes the current device time counter in nanoseconds when the command identified
	// by the event is enqueued in a command-queue by the host.
	//
	// Returned type: uint64
	ProfilingCommandQueuedInfo EventProfilingInfoName = C.CL_PROFILING_COMMAND_QUEUED
	// ProfilingCommandSubmitInfo describes the current device time counter in nanoseconds when the command identified
	// by the event that has been enqueued is submitted by the host to the device associated with the command-queue.
	//
	// Returned type: uint64
	ProfilingCommandSubmitInfo EventProfilingInfoName = C.CL_PROFILING_COMMAND_SUBMIT
	// ProfilingCommandStartInfo describes the current device time counter in nanoseconds when the command identified
	// by the event starts execution on the device.
	//
	// Returned type: uint64
	ProfilingCommandStartInfo EventProfilingInfoName = C.CL_PROFILING_COMMAND_START
	// ProfilingCommandEndInfo describes the current device time counter in nanoseconds when the command identified
	// by the event has finished execution on the device.
	//
	// Returned type: uint64
	ProfilingCommandEndInfo EventProfilingInfoName = C.CL_PROFILING_COMMAND_END
	// ProfilingCommandCompleteInfo describes the current device time counter in nanoseconds when the command identified
	// by the event and any child commands enqueued by this command on the device have finished execution.
	//
	// Returned type: uint64
	// Since: 2.0
	ProfilingCommandCompleteInfo EventProfilingInfoName = C.CL_PROFILING_COMMAND_COMPLETE
)

// EventProfilingInfo returns profiling information for the command associated with event if profiling is enabled.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetEventProfilingInfo.html
func EventProfilingInfo(event Event, paramName EventProfilingInfoName, paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetEventProfilingInfo(
		event.handle(),
		C.cl_profiling_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// SetEventCallback registers a user callback function for a specific command execution status.
//
// The command execution callback values for which a callback can be registered are: EventCommandSubmittedStatus,
// EventCommandRunningStatus, or EventCommandCompleteStatus.
//
// The provided callback will receive an error in case execution failed, or nil if the requested execution status
// has been reached.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetEventCallback.html
func SetEventCallback(event Event, callbackType EventCommandExecutionStatus, callback func(error)) error {
	callbackUserData, err := userDataFor(callback)
	if err != nil {
		return err
	}
	status := C.cl30SetEventCallback(event.handle(), C.cl_int(callbackType), callbackUserData.ptr)
	if status != C.CL_SUCCESS {
		callbackUserData.Delete()
		return StatusError(status)
	}
	return nil
}

//export cl30GoEventCallback
func cl30GoEventCallback(_ Event, commandStatus C.cl_int, userData *C.uintptr_t) {
	callbackUserData := userDataFrom(userData)
	callback := callbackUserData.Value().(func(error))
	callbackUserData.Delete()
	var err error
	if commandStatus < 0 {
		err = StatusError(commandStatus)
	}
	callback(err)
}

// EnqueueMarkerWithWaitList enqueues a marker command which waits for either a list of events to complete,
// or all previously enqueued commands to complete.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueMarkerWithWaitList.html
func EnqueueMarkerWithWaitList(commandQueue CommandQueue, waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueMarkerWithWaitList(
		commandQueue.handle(),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueBarrierWithWaitList is a synchronization point that enqueues a barrier operation.
//
// The barrier command either waits for a list of events to complete, or if the list is empty it waits for all
// commands previously enqueued in commandQueue to complete before it completes. This command blocks command
// execution, that is, any following commands enqueued after it do not execute until it completes.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueBarrierWithWaitList.html
func EnqueueBarrierWithWaitList(commandQueue CommandQueue, waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueBarrierWithWaitList(
		commandQueue.handle(),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}
