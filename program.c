#include "api.h"

extern void cl30GoProgramBuildCallback(cl_program, uintptr_t *);

static CL_CALLBACK void cl30CProgramBuildCallback(cl_program program, void *userData)
{
    cl30GoProgramBuildCallback(program, (uintptr_t *)(userData));
}

cl_int cl30BuildProgram(cl_program program,
    cl_uint numDevices, cl_device_id *devices,
    char *options, uintptr_t *userData)
{
    return clBuildProgram(program, numDevices, devices, options,
        (userData != NULL) ? cl30CProgramBuildCallback : NULL, userData);
}

extern void cl30GoProgramCompileCallback(cl_program, uintptr_t *);

static CL_CALLBACK void cl30CProgramCompileCallback(cl_program program, void *userData)
{
    cl30GoProgramCompileCallback(program, (uintptr_t *)(userData));
}

cl_int cl30CompileProgram(cl_program program,
    cl_uint numDevices, cl_device_id *devices,
    char *options,
    cl_uint numInputHeaders, cl_program *headers, char const **includeNames,
    uintptr_t *userData)
{
    return clCompileProgram(program, numDevices, devices, options,
        numInputHeaders, headers, includeNames,
        (userData != NULL) ? cl30CProgramCompileCallback : NULL, userData);
}

extern void cl30GoProgramLinkCallback(cl_program, uintptr_t *);

static CL_CALLBACK void cl30CProgramLinkCallback(cl_program program, void *userData)
{
    cl30GoProgramLinkCallback(program, (uintptr_t *)(userData));
}

cl_program cl30LinkProgram(cl_context context,
    cl_uint numDevices, cl_device_id *devices,
    char *options,
    cl_uint numInputPrograms, cl_program *programs,
    uintptr_t *userData,
    cl_int *errReturn)
{
    return clLinkProgram(context, numDevices, devices, options,
        numInputPrograms, programs,
        (userData != NULL) ? cl30CProgramLinkCallback : NULL, userData,
        errReturn);
}
