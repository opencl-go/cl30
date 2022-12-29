package cl30_test

import (
	"reflect"
	"testing"

	cl "github.com/opencl-go/cl30"
)

func TestHostMemoryBytesForNil(t *testing.T) {
	t.Parallel()
	result := cl.HostMemoryBytes(nil)
	if len(result) != 0 {
		t.Errorf("byte size mismatch")
	}
}

func TestFixedHostMemory(t *testing.T) {
	t.Parallel()
	mem := cl.AllocFixedHostMemory(10)
	defer mem.Free()
	if mem.Size() != 10 {
		t.Errorf("byte size mismatch")
	}
	if mem.Pointer() == nil {
		t.Errorf("no pointer available")
	}
	bytes := cl.HostMemoryBytes(mem)
	if len(bytes) != 10 {
		t.Errorf("size mismatch of bytes")
	}
	mem.Free()
	if mem.Size() != 0 {
		t.Errorf("size should be zero after free")
	}
	if mem.Pointer() != nil {
		t.Errorf("pointer should be nil after free")
	}
}

func TestFixedHostMemoryDefaults(t *testing.T) {
	t.Parallel()
	var mem *cl.FixedHostMemory
	if mem.Size() != 0 {
		t.Errorf("default size should be zero")
	}
	if mem.Pointer() != nil {
		t.Errorf("default pointer should be nil")
	}
	mem.Free() // no explicit test, yet should cover extra code
}

func TestHostValueOf(t *testing.T) {
	mem := cl.HostValueOf(uint32(0x11111111))
	if got, want := mem.Size(), 4; got != want {
		t.Errorf("size not matching. got=%d want=%d", got, want)
	}
	if !reflect.DeepEqual(cl.HostMemoryBytes(mem), []byte{0x11, 0x11, 0x11, 0x11}) {
		t.Errorf("memory access invalid. got=%v", cl.HostMemoryBytes(mem))
	}
}

func TestHostReferenceOf(t *testing.T) {
	var value uint32
	mem := cl.HostReferenceOf(&value)
	if got, want := mem.Size(), 4; got != want {
		t.Errorf("size not matching. got=%d want=%d", got, want)
	}
	copy(cl.HostMemoryBytes(mem), []byte{0xCD, 0xCD, 0xCD, 0xCD})
	if value != 0xCDCDCDCD {
		t.Errorf("memory access invalid. got=0x%X", value)
	}
}

func TestHostReferenceOfNil(t *testing.T) {
	mem := cl.HostReferenceOf((*uint16)(nil))
	if got, want := mem.Size(), 0; got != want {
		t.Errorf("size not matching. got=%d want=%d", got, want)
	}
	if len(cl.HostMemoryBytes(mem)) != 0 {
		t.Errorf("memory access invalid. got=0x%X", cl.HostMemoryBytes(mem))
	}
}

func TestHostVectorOf(t *testing.T) {
	value := make([]uint32, 3)
	mem := cl.HostVectorOf(value)
	if got, want := mem.Size(), 4*3; got != want {
		t.Errorf("size not matching. got=%d want=%d", got, want)
	}
	copy(cl.HostMemoryBytes(mem), []byte{0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x33, 0x33, 0x33, 0x33})
	if !reflect.DeepEqual(value, []uint32{0x11111111, 0x22222222, 0x33333333}) {
		t.Errorf("memory access invalid. got=%v", value)
	}
}

func TestHostVectorOfNil(t *testing.T) {
	var value []uint16
	mem := cl.HostVectorOf(value)
	if got, want := mem.Size(), 0; got != want {
		t.Errorf("size not matching. got=%d want=%d", got, want)
	}
	if len(cl.HostMemoryBytes(mem)) != 0 {
		t.Errorf("memory access invalid. got=0x%X", cl.HostMemoryBytes(mem))
	}
}
