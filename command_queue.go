package cl30

// #include "api.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// CommandQueue describes a sequence of events for OpenCL operations.
// Create a new command-queue with CreateCommandQueueWithProperties().
type CommandQueue uintptr

func (cq CommandQueue) handle() C.cl_command_queue {
	return *(*C.cl_command_queue)(unsafe.Pointer(&cq))
}

// String provides a readable presentation of the command-queue identifier.
// It is based on the numerical value of the underlying pointer.
func (cq CommandQueue) String() string {
	return fmt.Sprintf("0x%X", uintptr(cq))
}

const (
	// QueuePropertiesProperty is a bitfield.
	//
	// Use WithQueuePropertyFlags() for convenience.
	//
	// Property value type: CommandQueuePropertiesFlags
	QueuePropertiesProperty uint64 = C.CL_QUEUE_PROPERTIES
	// QueueSizeProperty specifies the size of the device queue in bytes.
	// This can only be specified if QueueOnDevice is set in QueuePropertiesProperty.
	// This must be a value less than or equal DeviceQueueOnDeviceMaxSizeInfo.
	//
	// For best performance, this should be less than or equal DeviceQueueOnDevicePreferredSizeInfo.
	//
	// If QueueSizeProperty is not specified, the device queue is created with DeviceQueueOnDevicePreferredSizeInfo as
	// the size of the queue.
	//
	// Use WithQueueSize() for convenience.
	//
	// Property value type: Uint
	// Since: 2.0
	QueueSizeProperty uint64 = C.CL_QUEUE_SIZE
)

// CommandQueuePropertiesFlags is used to determine DeviceQueueOnDevicePropertiesInfo and DeviceQueueOnHostPropertiesInfo
// with DeviceInfo(), as well as QueuePropertiesProperty for CreateCommandQueueWithProperties() and CreateCommandQueue().
type CommandQueuePropertiesFlags C.cl_command_queue_properties

const (
	// QueueOutOfOrderExecModeEnable determines whether the commands queued in the command-queue are executed
	// in-order or out-of-order. If set, the commands in the command-queue are executed out-of-order.
	// Otherwise, commands are executed in-order.
	QueueOutOfOrderExecModeEnable CommandQueuePropertiesFlags = C.CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE
	// QueueProfilingEnable enables or disables profiling of commands in the command-queue. If set,
	// the profiling of commands is enabled. Otherwise, profiling of commands is disabled.
	QueueProfilingEnable CommandQueuePropertiesFlags = C.CL_QUEUE_PROFILING_ENABLE
	// QueueOnDevice indicates that this is a device queue. If QueueOnDevice is set, QueueOutOfOrderExecModeEnable
	// must also be set. Only out-of-order device queues are supported.
	//
	// Since: 2.0
	QueueOnDevice CommandQueuePropertiesFlags = C.CL_QUEUE_ON_DEVICE
	// QueueOnDeviceDefault indicates that this is the default device queue.
	// This can only be used with QueueOnDevice.
	//
	// Since: 2.0
	QueueOnDeviceDefault CommandQueuePropertiesFlags = C.CL_QUEUE_ON_DEVICE_DEFAULT
)

// CommandQueueProperty is one entry of properties which are taken into account when creating command-queue objects.
type CommandQueueProperty []uint64

// WithQueueSize is a convenience function to create a valid QueueSizeProperty.
// Use it in combination with CreateCommandQueueWithProperties().
func WithQueueSize(bytes Uint) CommandQueueProperty {
	return CommandQueueProperty{QueueSizeProperty, uint64(bytes)}
}

// WithQueuePropertyFlags is a convenience function to create a valid QueuePropertiesProperty.
// Use it in combination with CreateCommandQueueWithProperties().
func WithQueuePropertyFlags(flags CommandQueuePropertiesFlags) CommandQueueProperty {
	return CommandQueueProperty{QueuePropertiesProperty, uint64(flags)}
}

// CreateCommandQueueWithProperties create a host or device command-queue on a specific device.
//
// If QueuePropertiesProperty is not specified in the properties, an in-order host command-queue is created for
// the specified device.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateCommandQueueWithProperties.html
func CreateCommandQueueWithProperties(context Context, deviceID DeviceID, properties ...CommandQueueProperty) (CommandQueue, error) {
	var rawPropertyList []uint64
	for _, property := range properties {
		rawPropertyList = append(rawPropertyList, property...)
	}
	var rawProperties unsafe.Pointer
	if len(properties) > 0 {
		rawPropertyList = append(rawPropertyList, 0)
		rawProperties = unsafe.Pointer(&rawPropertyList[0])
	}
	var status C.cl_int
	commandQueue := C.clCreateCommandQueueWithProperties(
		context.handle(),
		deviceID.handle(),
		(*C.cl_command_queue_properties)(rawProperties),
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return CommandQueue(*((*uintptr)(unsafe.Pointer(&commandQueue)))), nil
}

// RetainCommandQueue increments the commandQueue reference count.
//
// CreateCommandQueueWithProperties() and CreateCommandQueue() perform an implicit retain.
// This is very helpful for 3rd party libraries, which typically get a command-queue passed to them by the application.
// However, it is possible that the application may delete the command-queue without informing the library.
// Allowing functions to attach to (i.e. retain) and release a command-queue solves the problem of a command-queue
// being used by a library no longer being valid.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clRetainCommandQueue.html
func RetainCommandQueue(commandQueue CommandQueue) error {
	status := C.clRetainCommandQueue(commandQueue.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// ReleaseCommandQueue decrements the commandQueue reference count.
//
// After the commandQueue reference count becomes zero and all commands queued to commandQueue have finished
// (eg. kernel-instances, memory object updates etc.), the command-queue is deleted.
//
// ReleaseCommandQueue() performs an implicit flush to issue any previously queued OpenCL commands in commandQueue.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clReleaseCommandQueue.html
func ReleaseCommandQueue(commandQueue CommandQueue) error {
	status := C.clReleaseCommandQueue(commandQueue.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// CommandQueueInfoName identifies properties of a command-queue, which can be queried with CommandQueueInfo().
type CommandQueueInfoName C.cl_command_queue_info

const (
	// QueueContextInfo returns the context specified when the command-queue is created.
	//
	// Returned type: Context
	QueueContextInfo CommandQueueInfoName = C.CL_QUEUE_CONTEXT
	// QueueDeviceInfo returns the device specified when the command-queue is created
	//
	// Returned type: DeviceID
	QueueDeviceInfo CommandQueueInfoName = C.CL_QUEUE_DEVICE
	// QueueReferenceCountInfo returns the command-queue reference count.
	//
	// Note: The reference count returned should be considered immediately stale. It is unsuitable for
	// general use in applications. This feature is provided for identifying memory leaks.
	//
	// Returned type: Uint
	QueueReferenceCountInfo CommandQueueInfoName = C.CL_QUEUE_REFERENCE_COUNT
	// QueuePropertiesInfo returns the currently specified properties for the command-queue.
	// These properties are specified by the value associated with the QueuePropertiesProperty passed in properties
	// argument in CreateCommandQueueWithProperties(), or the value of the properties argument in CreateCommandQueue().
	//
	// Returned type: uint64
	QueuePropertiesInfo CommandQueueInfoName = C.CL_QUEUE_PROPERTIES
	// QueuePropertiesArrayInfo returns the properties argument specified in CreateCommandQueueWithProperties().
	//
	// Returned type: []uint64
	// Since: 3.0
	QueuePropertiesArrayInfo CommandQueueInfoName = C.CL_QUEUE_PROPERTIES_ARRAY
	// QueueSizeInfo returns the size of the device command-queue. To be considered valid for this query,
	// the command-queue must be a device command-queue.
	//
	// Returned type: Uint
	// Since: 2.0
	QueueSizeInfo CommandQueueInfoName = C.CL_QUEUE_SIZE
	// QueueDeviceDefaultInfo returns the current default command-queue for the underlying device.
	//
	// Returned type: CommandQueue
	// Since: 2.1
	QueueDeviceDefaultInfo CommandQueueInfoName = C.CL_QUEUE_DEVICE_DEFAULT
)

// CommandQueueInfo queries information about a command-queue.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character. For convenience, use CommandQueueInfoString().
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetCommandQueueInfo.html
func CommandQueueInfo(commandQueue CommandQueue, paramName CommandQueueInfoName, paramSize uint, paramValue unsafe.Pointer) (uint, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetCommandQueueInfo(
		commandQueue.handle(),
		C.cl_command_queue_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uint(sizeReturn), nil
}

// Flush issues all previously queued OpenCL commands in a command-queue to the device associated with the
// command-queue.
//
// All previously queued OpenCL commands in commandQueue are issued to the device associated with commandQueue.
// Flush() only guarantees that all queued commands to commandQueue will eventually be submitted to the appropriate
// device. There is no guarantee that they will be complete after Flush() returns.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clFlush.html
func Flush(commandQueue CommandQueue) error {
	status := C.clFlush(commandQueue.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// Finish blocks until all previously queued OpenCL commands in a command-queue are issued to the associated device
// and have completed.
//
// All previously queued OpenCL commands in commandQueue are issued to the associated device, and the function blocks
// until all previously queued commands have completed. Finish() does not return until all previously queued commands
// in commandQueue have been processed and completed. Finish() is also a synchronization point.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clFinish.html
func Finish(commandQueue CommandQueue) error {
	status := C.clFinish(commandQueue.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// SetDefaultDeviceCommandQueue replaces the default command-queue on the device.
//
// This function may be used to replace a default device command-queue created with CreateCommandQueueWithProperties()
// and the QueueOnDeviceDefault flag.
//
// Since: 2.1
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clSetDefaultDeviceCommandQueue.html
func SetDefaultDeviceCommandQueue(context Context, deviceID DeviceID, commandQueue CommandQueue) error {
	status := C.clSetDefaultDeviceCommandQueue(context.handle(), deviceID.handle(), commandQueue.handle())
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}
