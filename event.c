#include "api.h"

extern void cl30GoEventCallback(cl_event, cl_int, void*);

static CL_CALLBACK void cl30CEventCallback(cl_event event, cl_int commandStatus, void *userData)
{
    cl30GoEventCallback(event, commandStatus, (intptr_t *)(userData));
}

cl_int cl30SetEventCallback(cl_event event, cl_int callbackType, intptr_t *userData)
{
    return clSetEventCallback(event, callbackType, cl30CEventCallback, userData);
}
