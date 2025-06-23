package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func trace_dfs(ptr uintptr, visited map[uintptr]struct{}) {
	if ptr == 0 {
		return
	}

	if _, ok := visited[ptr]; ok {
		return
	}

	visited[ptr] = struct{}{}

	value := *(*uintptr)(unsafe.Pointer(ptr))

	if value != 0 {
		trace_dfs(value, visited)
	}
}

func keys_from_map(m map[uintptr]struct{}) []uintptr {
	r := make([]uintptr, len(m))
	i := 0
	for ptr := range m {
		r[i] = ptr
		i++
	}
	return r
}

func Trace(stacks [][]uintptr) []uintptr {
	visited := make(map[uintptr]struct{})

	for _, stack := range stacks {
		for _, ptr := range stack {
			if ptr != 0 {
				trace_dfs(ptr, visited)
			}
		}
	}

	return  keys_from_map(visited)
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}
	pointers := Trace(stacks)

	assert.ElementsMatch(t, pointers, expectedPointers)
}
