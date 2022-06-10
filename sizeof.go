package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Foo interface {
	Foo() Foo
	Foo2() Foo
	FooNil() Foo
}

type FooImpl struct {
	//foo FooImpl	// invalid recursive type <- memory size must be fixed
	foo          *FooImpl
	fooInterface Foo
	value        interface{}
	intVal       int
	int32Val     int32
	int64Val     int64
	stringVal    string
}

func (f FooImpl) Foo() Foo {
	return f.foo
}

func (f FooImpl) Foo2() Foo {
	return f.fooInterface
}

func (f FooImpl) FooNil() Foo {
	if f.foo == nil {
		return nil
	}
	return f.foo
}

func (f FooImpl) FooImpl() *FooImpl {
	return f.foo
}

func sizeof() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("sizeof", "=", 80))
	fmt.Println(center("", "=", 80))

	fmt.Println(structDefString(FooImpl{}))
	fmt.Println()

	foo := FooImpl{}
	fmt.Printf("foo:\t%+v\n", foo)
	foor := FooImpl{&foo, &foo, 42, 12, 13, 11, "bar bar bar bar bar"}
	fmt.Printf("foor:\t%+v\n", foor)
	fmt.Println()

	fmt.Println("foor := FooImpl{foo: &foo, value: 42}")
	fmt.Printf("Sizeof(foo):\t%v bytes\n", unsafe.Sizeof(foo))
	fmt.Printf("Sizeof(foor):\t%v bytes\n", unsafe.Sizeof(foor))
	fmt.Println()

	sh := (*reflect.StringHeader)(unsafe.Pointer(&(foor.stringVal)))
	fmt.Printf("foo.stringVal: %+v: %d\n", sh, unsafe.Sizeof(sh))
	fmt.Println()

	types := []reflect.Type{
		reflect.TypeOf(*new(int8)),
		reflect.TypeOf(*new(int16)),
		reflect.TypeOf(*new(int32)),
		reflect.TypeOf(*new(int64)),
		reflect.TypeOf(*new(uint8)),
		reflect.TypeOf(*new(uint16)),
		reflect.TypeOf(*new(uint32)),
		reflect.TypeOf(*new(uint64)),
		reflect.TypeOf(*new(float32)),
		reflect.TypeOf(*new(float64)),
		reflect.TypeOf(*new(complex64)),
		reflect.TypeOf(*new(complex128)),
		reflect.TypeOf(*new(int)),
		reflect.TypeOf(*new(uint)),
		reflect.TypeOf(*new(uintptr)),
		reflect.TypeOf(*new(string)),
		reflect.TypeOf(*new([]int64)),
		reflect.TypeOf([]int64{100}),
		reflect.TypeOf([10]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
	}
	fmt.Printf("%s%s%s\n", pad("Name", 16), pad("Kind", 16), "Size (shallow)")
	for _, t := range types {
		name := t.Name()
		switch t.Kind() {
		case reflect.Array:
			name = fmt.Sprintf("[%d]%s", t.Len(), t.Elem())
		case reflect.Slice:
			name = fmt.Sprintf("[]%s", t.Elem().Name())
		}
		fmt.Printf("%s%s%d\n", pad(name, 16), pad(t.Kind(), 16), t.Size())
	}
	fmt.Println()

	fmt.Printf("byte == uint8:  %v\n", reflect.TypeOf(*new(byte)) == reflect.TypeOf(*new(uint8)))
	fmt.Printf("rune == uint32: %v\n", reflect.TypeOf(*new(rune)) == reflect.TypeOf(*new(int32)))

	fmt.Println(SizeSummary(foor))
}
