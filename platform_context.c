#include "api.h"

extern void cl30GoContextErrorCallback(char *, uint8_t *, size_t, uintptr_t *);

static CL_CALLBACK void cl30CContextErrorCallback(char const *errorInfo,
    void const *privateInfoPtr, size_t privateInfoLen,
    uintptr_t *userData)
{
    cl30GoContextErrorCallback((char *)(errorInfo), (uint8_t *)(privateInfoPtr), privateInfoLen, (uintptr_t *)(userData));
}

cl_context cl30CreateContext(cl_context_properties *properties,
    cl_uint numDevices, cl_device_id *devices,
    uintptr_t *userData,
    cl_int *errcodeReturn)
{
    return clCreateContext(properties, numDevices, devices,
        (userData != NULL) ? cl30CContextErrorCallback : NULL, userData,
        errcodeReturn);
}

cl_context cl30CreateContextFromType(cl_context_properties *properties,
    cl_device_type deviceType,
    uintptr_t *userData,
    cl_int *errcodeReturn)
{
    return clCreateContextFromType(properties, deviceType,
        (userData != NULL) ? cl30CContextErrorCallback : NULL, userData,
        errcodeReturn);
}

extern void cl30GoContextDestructorCallback(cl_context, intptr_t *);

static CL_CALLBACK void cl30CContextDestructorCallback(cl_context context, void *userData)
{
    cl30GoContextDestructorCallback(context, (intptr_t *)(userData));
}

cl_int cl30SetContextDestructorCallback(cl_context context, intptr_t *userData)
{
    return clSetContextDestructorCallback(context, cl30CContextDestructorCallback, userData);
}
