#include "api.h"

extern void cl30GoProgramReleaseCallback(cl_program, uintptr_t *);

static CL_CALLBACK void cl30CProgramReleaseCallback(cl_program program, void *userData)
{
    cl30GoProgramReleaseCallback(program, (uintptr_t *)(userData));
}

cl_int cl30SetProgramReleaseCallback(cl_program program, uintptr_t *userData)
{
    return clSetProgramReleaseCallback(program, cl30CProgramReleaseCallback, userData);
}
