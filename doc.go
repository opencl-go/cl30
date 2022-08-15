// Package cl30 provides a wrapper API to OpenCL 3.0.
//
// If you require a different API level, refer to the opencl-go project (https://opencl-go.github.com) to see which
// versions are available.
//
// To build and work with this library, you need an OpenCL SDK installed on your system.
// Refer to the documentation on opencl-go (https://opencl-go.github.com) on how to do this.
//
// The API requires knowledge of the OpenCL API. While the wrapper hides some low-level C-API details,
// there is still heavy use of `unsafe.Pointer` and the potential for memory access-violations if used wrong.
//
// This library wraps/represents constants, types, and functions as closely to the original API as possible -
// while applying Go idioms.
// It also provides convenience functions where practical. For example: PlatformInfoString().
//
// The library uses a Go-idiomatic naming pattern that avoids name clashes:
// For example, "cl_device_local_mem_type" and "CL_DEVICE_LOCAL_MEM_TYPE" would both resolve to "DeviceLocalMemType".
// For the type of the constant, choose a suffix that matches the parameter where it is used.
// For example, "cl_platform_info" is expressed as "PlatformInfoName" as it identifies the name of the parameter
// in PlatformInfo().
//
// References:
//
// See also: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/
//
// The API documentation is, in part, based on the official asciidoctor source files from https://github.com/KhronosGroup/OpenCL-Docs,
// licensed under the Creative Commons Attribution 4.0 International License; see https://creativecommons.org/licenses/by/4.0/ .
package cl30
