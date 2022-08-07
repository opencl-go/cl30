#include "api.h"

extern void cl30GoContextErrorCallback(char *, uint8_t *, size_t, intptr_t);

CL_CALLBACK void cl30CContextErrorCallback(char const *errorInfo,
    void const *privateInfoPtr, size_t privateInfoLen,
    void *userData)
{
    cl30GoContextErrorCallback((char *)(errorInfo), (uint8_t *)(privateInfoPtr), (size_t)(privateInfoLen), (intptr_t)(userData));
}

cl_context cl30CreateContext(cl_context_properties *properties,
    cl_uint numDevices, cl_device_id *devices,
    void *notify, intptr_t userData,
    cl_int *errcodeReturn)
{
    return clCreateContext(properties, numDevices, devices,
        (void (CL_CALLBACK *)(char const *, void const *, size_t, void *))(notify), (void *)(userData),
        errcodeReturn);
}

cl_context cl30CreateContextFromType(cl_context_properties *properties,
    cl_device_type deviceType,
    void *notify, intptr_t userData,
    cl_int *errcodeReturn)
{
    return clCreateContextFromType(properties, deviceType,
        (void (CL_CALLBACK *)(char const *, void const *, size_t, void *))(notify), (void *)(userData),
        errcodeReturn);
}

cl_int cl30SetContextDestructorCallback(cl_context context, void *notify, void *userData)
{
    return clSetContextDestructorCallback(context,
        (void (CL_CALLBACK *)(cl_context, void *))(notify),
        userData);
}

extern void cl30GoContextDestructorCallback(cl_context, void*);

CL_CALLBACK void cl30CContextDestructorCallback(cl_context context, void *userData)
{
    cl30GoContextDestructorCallback(context, userData);
}
