// Package cl30 provides a wrapper API to OpenCL 3.0.
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
