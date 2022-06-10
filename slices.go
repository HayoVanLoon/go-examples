package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func iterateSlice[T any](xs []T) {
	fmt.Println("Using pointer arithmetic ...")
	if len(xs) == 0 {
		panic("empty or nil slice")
	}
	start := unsafe.Pointer(&xs[0])
	end := unsafe.Pointer(&xs[len(xs)-1])
	step := reflect.TypeOf(xs).Elem().Size()
	for p := start; ; p = unsafe.Pointer(uintptr(p) + step) {
		fmt.Printf("(%x %v),\n", p, *(*int32)(p))
		if uintptr(p) == uintptr(end) {
			// avoids calculating an illegal pointer value
			return
		}
	}
}

func slices() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("slices", "=", 80))
	fmt.Println(center("", "=", 80))

	fmt.Println("\nxs := []int32{13, 42, 11}")
	xs := []int32{13, 42, 11}
	fmt.Printf("(%x %v),\n", &xs, xs)
	for i := range xs {
		fmt.Printf("(%x %v),\n", &xs[i], xs[i])
	}
	fmt.Println("size of xs: ", SizeSummary(xs))

	fmt.Println("\nys := []*int{&a, &b, &c, &d}")
	a, b, c, d := 33, 22, 11, 44
	ys := []*int{&a, &b, &c, &d}
	fmt.Printf("(%x %v),\n", &ys, ys)
	for i := range ys {
		fmt.Printf("(%x %v),\n", &ys[i], ys[i])
	}
	fmt.Println("size of ys: ", SizeSummary(ys))

	fmt.Println("\nzs := &[]int{13, 42, 11}")
	zs := &[]int{13, 42, 11}
	fmt.Printf("(%x %v),\n", &zs, zs)
	for i := range *zs {
		fmt.Printf("(%x %v),\n", &(*zs)[i], (*zs)[i])
	}
	fmt.Println("size of zs: ", SizeSummary(zs))

	fmt.Println()
	iterateSlice(xs)
}
