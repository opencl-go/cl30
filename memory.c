#include "api.h"

extern void cl30GoMemObjectDestructorCallback(cl_mem, uintptr_t *);

static CL_CALLBACK void cl30CMemObjectDestructorCallback(cl_mem mem, void *userData)
{
    cl30GoMemObjectDestructorCallback(mem, (uintptr_t *)(userData));
}

cl_int cl30SetMemObjectDestructorCallback(cl_mem mem, uintptr_t *userData)
{
    return clSetMemObjectDestructorCallback(mem, cl30CMemObjectDestructorCallback, userData);
}
