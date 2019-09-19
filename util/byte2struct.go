package util

import (
	"go_stress_test/entity"
	"unsafe"
)

type TestStructTobytes struct {
	data int64
	msg  int8
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func StructToByte(srcStruct *entity.Header) []byte {
	Len := unsafe.Sizeof(*srcStruct)
	testBytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(srcStruct)),
		cap:  int(Len),
		len:  int(Len),
	}
	data := *(*[]byte)(unsafe.Pointer(testBytes))

	return data
}
