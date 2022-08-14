package cl30_test

import (
	"testing"
	"unsafe"

	cl "github.com/opencl-go/cl30"
)

func TestBufferRegionSize(t *testing.T) {
	t.Parallel()
	if (cl.BufferRegionByteSize != unsafe.Sizeof(cl.BufferRegion{})) {
		t.Errorf("byte size mismatch")
	}
}
