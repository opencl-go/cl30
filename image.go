package cl30

// #include "api.h"
import "C"
import "unsafe"

// ChannelOrder describes the sequence and nature of the color channels of an image.
type ChannelOrder C.cl_channel_order

// Constants for possible ChannelOrder values.
const (
	ChannelOrderR         ChannelOrder = C.CL_R
	ChannelOrderA         ChannelOrder = C.CL_A
	ChannelOrderRg        ChannelOrder = C.CL_RG
	ChannelOrderRa        ChannelOrder = C.CL_RA
	ChannelOrderRgb       ChannelOrder = C.CL_RGB
	ChannelOrderRgba      ChannelOrder = C.CL_RGBA
	ChannelOrderBgra      ChannelOrder = C.CL_BGRA
	ChannelOrderArgb      ChannelOrder = C.CL_ARGB
	ChannelOrderIntensity ChannelOrder = C.CL_INTENSITY
	ChannelOrderLuminance ChannelOrder = C.CL_LUMINANCE

	ChannelOrderRx   ChannelOrder = C.CL_Rx
	ChannelOrderRgx  ChannelOrder = C.CL_RGx
	ChannelOrderRgbx ChannelOrder = C.CL_RGBx

	ChannelOrderDepth   ChannelOrder = C.CL_DEPTH
	ChannelOrderStencil ChannelOrder = C.CL_DEPTH_STENCIL

	ChannelOrderSrgb  ChannelOrder = C.CL_sRGB
	ChannelOrderSrgbx ChannelOrder = C.CL_sRGBx
	ChannelOrderSrgba ChannelOrder = C.CL_sRGBA
	ChannelOrderSbgra ChannelOrder = C.CL_sBGRA
	ChannelOrderAbgr  ChannelOrder = C.CL_ABGR
)

// ChannelType describes the resolution and value range of color channels.
type ChannelType C.cl_channel_type

// Constants for possible ChannelType values.
const (
	ChannelTypeSnormInt8      ChannelType = C.CL_SNORM_INT8
	ChannelTypeSnormInt16     ChannelType = C.CL_SNORM_INT16
	ChannelTypeUnormInt8      ChannelType = C.CL_UNORM_INT8
	ChannelTypeUnormInt16     ChannelType = C.CL_UNORM_INT16
	ChannelTypeUnormShort565  ChannelType = C.CL_UNORM_SHORT_565
	ChannelTypeUnormShort555  ChannelType = C.CL_UNORM_SHORT_555
	ChannelTypeUnormInt101010 ChannelType = C.CL_UNORM_INT_101010
	ChannelTypeSignedInt8     ChannelType = C.CL_SIGNED_INT8
	ChannelTypeSignedInt16    ChannelType = C.CL_SIGNED_INT16
	ChannelTypeSignedInt32    ChannelType = C.CL_SIGNED_INT32
	ChannelTypeUnsignedInt8   ChannelType = C.CL_UNSIGNED_INT8
	ChannelTypeUnsignedInt16  ChannelType = C.CL_UNSIGNED_INT16
	ChannelTypeUnsignedInt32  ChannelType = C.CL_UNSIGNED_INT32
	ChannelTypeHalfFloat      ChannelType = C.CL_HALF_FLOAT
	ChannelTypeFloat          ChannelType = C.CL_FLOAT

	ChannelTypeUnormInt24 ChannelType = C.CL_UNORM_INT24

	ChannelTypeUnormInt1010102 ChannelType = C.CL_UNORM_INT_101010_2
)

// ImageFormatByteSize is the size, in bytes, of the ImageFormat structure.
const ImageFormatByteSize = unsafe.Sizeof(C.cl_image_format{})

// ImageFormat describes how the bytes of the byte buffer shall be interpreted as pixel values.
type ImageFormat struct {
	ChannelOrder ChannelOrder
	ChannelType  ChannelType
}

// ImageDescByteSize is the size, in bytes, of the ImageDesc structure.
const ImageDescByteSize = unsafe.Sizeof(C.cl_image_desc{})

// ImageDesc describes the dimensions of the image.
type ImageDesc struct {
	ImageType    MemObjectType
	Width        uintptr
	Height       uintptr
	Depth        uintptr
	ArraySize    uintptr
	RowPitch     uintptr
	SlicePitch   uintptr
	NumMipLevels uint32
	NumSamples   uint32
	MemObject    MemObject
}

// CreateImage creates a 1D image, 1D image buffer, 1D image array, 2D image, 2D image array, or 3D image object.
//
// Since: 1.2
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateImage.html
func CreateImage(context Context, flags MemFlags, format ImageFormat, desc ImageDesc, hostPtr unsafe.Pointer) (MemObject, error) {
	var status C.cl_int
	mem := C.clCreateImage(
		context.handle(),
		C.cl_mem_flags(flags),
		(*C.cl_image_format)(unsafe.Pointer(&format)),
		(*C.cl_image_desc)(unsafe.Pointer(&desc)),
		hostPtr,
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return MemObject(*((*uintptr)(unsafe.Pointer(&mem)))), nil
}

// CreateImageWithProperties creates a 1D image, 1D image buffer, 1D image array, 2D image, 2D image array,
// or 3D image object.
//
// Since: 3.0
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clCreateImageWithProperties.html
func CreateImageWithProperties(context Context, flags MemFlags, format ImageFormat, desc ImageDesc, hostPtr unsafe.Pointer,
	properties ...MemProperty) (MemObject, error) {
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
	mem := C.clCreateImageWithProperties(
		context.handle(),
		(*C.cl_mem_properties)(rawProperties),
		C.cl_mem_flags(flags),
		(*C.cl_image_format)(unsafe.Pointer(&format)),
		(*C.cl_image_desc)(unsafe.Pointer(&desc)),
		hostPtr,
		&status)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return MemObject(*((*uintptr)(unsafe.Pointer(&mem)))), nil
}

// SupportedImageFormats returns the list of image formats supported by an OpenCL implementation.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetSupportedImageFormats.html
func SupportedImageFormats(context Context, flags MemFlags, imageType MemObjectType) ([]ImageFormat, error) {
	requiredCount := C.cl_uint(0)
	status := C.clGetSupportedImageFormats(
		context.handle(),
		C.cl_mem_flags(flags),
		C.cl_mem_object_type(imageType),
		0,
		nil,
		&requiredCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	if requiredCount == 0 {
		return nil, nil
	}
	formats := make([]ImageFormat, int(requiredCount))
	returnedCount := C.cl_uint(0)
	status = C.clGetSupportedImageFormats(
		context.handle(),
		C.cl_mem_flags(flags),
		C.cl_mem_object_type(imageType),
		requiredCount,
		(*C.cl_image_format)(unsafe.Pointer(&formats[0])),
		&returnedCount)
	if status != C.CL_SUCCESS {
		return nil, StatusError(status)
	}
	return formats[:returnedCount], nil
}

// MappedImage describes an image as it was mapped into host memory.
type MappedImage struct {
	Ptr        unsafe.Pointer
	RowPitch   uintptr
	SlicePitch uintptr
}

// EnqueueMapImage enqueues a command to map a region of an image object into the host address space and
// returns a description of this mapped region.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueMapImage.html
func EnqueueMapImage(commandQueue CommandQueue,
	image MemObject, blocking bool, flags MapFlags, origin, region [3]uintptr,
	waitList []Event, event *Event) (MappedImage, error) {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	var mapped MappedImage
	var status C.cl_int
	mapped.Ptr = C.clEnqueueMapImage(
		commandQueue.handle(),
		image.handle(),
		C.cl_bool(BoolFrom(blocking)),
		C.cl_map_flags(flags),
		(*C.size_t)(unsafe.Pointer(&origin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		(*C.size_t)(unsafe.Pointer(&mapped.RowPitch)),
		(*C.size_t)(unsafe.Pointer(&mapped.SlicePitch)),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)),
		&status)
	if status != C.CL_SUCCESS {
		return MappedImage{}, StatusError(status)
	}
	return mapped, nil
}

// ImageInfoName identifies properties of an image object, which can be queried with ImageInfo().
type ImageInfoName C.cl_image_info

const (
	// ImageFormatInfo return the image format descriptor specified when the image was created.
	//
	// Returned type: ImageFormat
	ImageFormatInfo ImageInfoName = C.CL_IMAGE_FORMAT
	// ImageElementSizeInfo returns the size of each element of the image memory object given by image in bytes.
	//
	// Returned type: uintptr
	ImageElementSizeInfo ImageInfoName = C.CL_IMAGE_ELEMENT_SIZE
	// ImageRowPitchInfo returns the row pitch in bytes of a row of elements of the image object given by image.
	//
	// Returned type: uintptr
	ImageRowPitchInfo ImageInfoName = C.CL_IMAGE_ROW_PITCH
	// ImageSlicePitchInfo returns the slice pitch in bytes of a 2D slice for the 3D image object or size of each
	// image in a 1D or 2D image array given by image.
	//
	// Returned type: uintptr
	ImageSlicePitchInfo ImageInfoName = C.CL_IMAGE_SLICE_PITCH
	// ImageWidthInfo returns the width of the image in pixels.
	//
	// Returned type: uintptr
	ImageWidthInfo ImageInfoName = C.CL_IMAGE_WIDTH
	// ImageHeightInfo returns the height of the image in pixels.
	// For a 1D image, 1D image buffer and 1D image array object, height is 0.
	//
	// Returned type: uintptr
	ImageHeightInfo ImageInfoName = C.CL_IMAGE_HEIGHT
	// ImageDepthInfo returns the depth of the image in pixels.
	// For a 1D image, 1D image buffer, 2D image or 1D and 2D image array object, depth is 0.
	//
	// Returned type: uintptr
	ImageDepthInfo ImageInfoName = C.CL_IMAGE_DEPTH
	// ImageArraySizeInfo returns the number of images in the image array.
	// If image is not an image array, 0 is returned.
	//
	// Returend type: uintptr
	// Since: 1.2
	ImageArraySizeInfo ImageInfoName = C.CL_IMAGE_ARRAY_SIZE
	// ImageBufferInfo returns the buffer object associated with image.
	//
	// Returned type: MemObject
	// Since: 1.2
	// Deprecated: 2.0
	ImageBufferInfo ImageInfoName = C.CL_IMAGE_BUFFER
	// ImageNumMipLevelsInfo returns the MIP level count associated with the image.
	//
	// Returned type: uint32
	// Since: 1.2
	ImageNumMipLevelsInfo ImageInfoName = C.CL_IMAGE_NUM_MIP_LEVELS
	// ImageNumSamplesInfo returns the sample count associated with the image.
	//
	// Returned type: uint32
	// Since: 1.2
	ImageNumSamplesInfo ImageInfoName = C.CL_IMAGE_NUM_SAMPLES
)

// ImageInfo returns information specific to an image object.
//
// The provided size need to specify the size of the available space pointed to the provided value in bytes.
//
// The returned number is the required size, in bytes, for the queried information.
// Call the function with a zero size and nil value to request the required size. This helps in determining
// the necessary space for dynamic information, such as arrays.
//
// Raw strings are with a terminating NUL character.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clGetImageInfo.html
func ImageInfo(image MemObject, paramName ImageInfoName, paramSize uintptr, paramValue unsafe.Pointer) (uintptr, error) {
	sizeReturn := C.size_t(0)
	status := C.clGetImageInfo(
		image.handle(),
		C.cl_image_info(paramName),
		C.size_t(paramSize),
		paramValue,
		&sizeReturn)
	if status != C.CL_SUCCESS {
		return 0, StatusError(status)
	}
	return uintptr(sizeReturn), nil
}

// EnqueueReadImage enqueues a command to read from an image or image array object to host memory.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueReadImage.html
func EnqueueReadImage(commandQueue CommandQueue, image MemObject, blocking bool, origin, region [3]uintptr,
	rowPitch, slicePitch uintptr, ptr HostMemory,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueReadImage(
		commandQueue.handle(),
		image.handle(),
		C.cl_bool(BoolFrom(blocking)),
		(*C.size_t)(unsafe.Pointer(&origin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		C.size_t(rowPitch),
		C.size_t(slicePitch),
		ResolvePointer(ptr, !blocking, "ptr"),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueWriteImage enqueues a command to write to an image or image array object from host memory.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueWriteImage.html
func EnqueueWriteImage(commandQueue CommandQueue, image MemObject, blocking bool, origin, region [3]uintptr,
	rowPitch, slicePitch uintptr, ptr HostMemory,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueWriteImage(
		commandQueue.handle(),
		image.handle(),
		C.cl_bool(BoolFrom(blocking)),
		(*C.size_t)(unsafe.Pointer(&origin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		C.size_t(rowPitch),
		C.size_t(slicePitch),
		ResolvePointer(ptr, !blocking, "ptr"),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueFillImage enqueues a command to fill an image object with a specified color.
//
// The fill color is a single floating point value if the channel order is ChannelOrderDepth.
// Otherwise, the fill color is a four-component RGBA floating-point color value if the image channel data type
// is not an unnormalized signed or unsigned integer type, is a four-component signed integer value if the image
// channel data type is an unnormalized signed integer type and is a four-component unsigned integer value if the image
// channel data type is an unnormalized unsigned integer type.
// The fill color will be converted to the appropriate image channel format and order associated with image.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueFillImage.html
func EnqueueFillImage(commandQueue CommandQueue, image MemObject, fillColor HostMemory, origin, region [3]uintptr,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueFillImage(
		commandQueue.handle(),
		image.handle(),
		ResolvePointer(fillColor, false, "fillColor"),
		(*C.size_t)(unsafe.Pointer(&origin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueCopyImage enqueues a command to copy image objects.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueCopyImage.html
func EnqueueCopyImage(commandQueue CommandQueue, srcImage, dstImage MemObject, srcOrigin, dstOrigin, region [3]uintptr,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueCopyImage(
		commandQueue.handle(),
		srcImage.handle(),
		dstImage.handle(),
		(*C.size_t)(unsafe.Pointer(&srcOrigin[0])),
		(*C.size_t)(unsafe.Pointer(&dstOrigin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueCopyImageToBuffer enqueues a command to copy an image object to a buffer object.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueCopyImageToBuffer.html
func EnqueueCopyImageToBuffer(commandQueue CommandQueue, srcImage, dstBuffer MemObject, srcOrigin, region [3]uintptr, dstOffset uintptr,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueCopyImageToBuffer(
		commandQueue.handle(),
		srcImage.handle(),
		dstBuffer.handle(),
		(*C.size_t)(unsafe.Pointer(&srcOrigin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		C.size_t(dstOffset),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}

// EnqueueCopyBufferToImage enqueues a command to copy a buffer object to an image object.
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/clEnqueueCopyBufferToImage.html
func EnqueueCopyBufferToImage(commandQueue CommandQueue, srcBuffer, dstImage MemObject, srcOffset uintptr, srcOrigin, region [3]uintptr,
	waitList []Event, event *Event) error {
	var rawWaitList unsafe.Pointer
	if len(waitList) > 0 {
		rawWaitList = unsafe.Pointer(&waitList[0])
	}
	status := C.clEnqueueCopyBufferToImage(
		commandQueue.handle(),
		srcBuffer.handle(),
		dstImage.handle(),
		C.size_t(srcOffset),
		(*C.size_t)(unsafe.Pointer(&srcOrigin[0])),
		(*C.size_t)(unsafe.Pointer(&region[0])),
		C.cl_uint(len(waitList)),
		(*C.cl_event)(rawWaitList),
		(*C.cl_event)(unsafe.Pointer(event)))
	if status != C.CL_SUCCESS {
		return StatusError(status)
	}
	return nil
}
