# Go wrapper library for OpenCL 3.0 

This library provides a complete wrapper for the OpenCL 3.0 API.
If you require a different API level, refer to [the opencl-go project][opencl-go] to see which versions are available.

**This is work-in-progress. The wrapper is not yet in a state to provide useful functionality**

## Usage

To build and work with this library, you need an OpenCL SDK installed on your system.
Refer to [the documentation on opencl-go][opencl-go] on how to do this. 

The API requires knowledge of the [OpenCL API][opencl-api]. While the wrapper hides some low-level C-API details,
there is still heavy use of `unsafe.Pointer` and the potential for memory access-violations if used wrong.

[opencl-api]: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/
[opencl-go]: https://opencl-go.github.com

## License

This project is based on the MIT License. See `LICENSE` file.

The API documentation is, in part, based on the official asciidoctor source files from https://github.com/KhronosGroup/OpenCL-Docs,
licensed under the Creative Commons Attribution 4.0 International License; see https://creativecommons.org/licenses/by/4.0/ .
