package cl30_test

import (
	"testing"
	"unsafe"

	cl "github.com/opencl-go/cl30"
)

func TestNameVersion(t *testing.T) {
	t.Parallel()
	if (cl.NameVersionByteSize != unsafe.Sizeof(cl.NameVersion{})) {
		t.Errorf("byte size mismatch")
	}
}
