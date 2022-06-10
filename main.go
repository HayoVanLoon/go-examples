package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func main() {
	fmt.Println()
	slices()

	fmt.Println()
	sizeof()

	fmt.Println()
	pointerReturns()

	fmt.Println()
	defers()

	fmt.Println()
	interfaces()

	fmt.Println()
	sharedMemoryNoMutex()

	fmt.Println()
	pipeline()

	fmt.Println()
	deferIntercept()

	fmt.Println()
	rogueRountine()

	fmt.Println("number of CPUs:", runtime.GOMAXPROCS(0))
	fmt.Println("running routines: ", runtime.NumGoroutine())

	//panic("show me")
}

func SizeSummary(v interface{}) string {
	return sizeSummary2(v, 0)
}

func sizeSummary2(v interface{}, depth int) string {
	v, typ := getInterestingType(v)
	val := reflect.ValueOf(v)

	prefix := ""
	for i := 0; i < depth; i += 1 {
		prefix += "\t"
	}
	b := &strings.Builder{}
	if val.Kind() == reflect.Struct {
		b.WriteString("{\n")
		for i := 0; i < val.NumField(); i += 1 {
			b.WriteString(sprintField(val.Field(i), typ.Field(i), depth))
		}
		_, _ = fmt.Fprintf(b, "%s}", prefix)
	} else {
		sprintField(val, reflect.StructField{}, depth)
	}
	return b.String()
}

func sprintField(v reflect.Value, sf reflect.StructField, depth int) string {
	v2 := safeInterface(v)
	v2, t := getInterestingType(v2)
	tName := t.Name()
	if tName == "" {
		tName = t.Elem().Name()
		if tName == "" && t.Kind() == reflect.Interface {
			tName = "interface{}"
		}
	}
	formatted := fmt.Sprintf("%v", v2)
	if t.Kind() == reflect.Struct {
		formatted = sizeSummary2(v2, depth+1)
	}
	prefix := ""
	for i := 0; i < depth; i += 1 {
		prefix += "\t"
	}
	return fmt.Sprintf("%s\t%s\t%s\t(%d B)\t: %+v\n", prefix, sf.Name, tName, t.Size(), formatted)
}

func getInterestingType(v interface{}) (interface{}, reflect.Type) {
	var typ reflect.Type
	for typ == nil { // hvl: might need recursive traversal for certain types
		typ = reflect.TypeOf(v)
		switch typ.Kind() {
		case reflect.String:
			v = *(*reflect.StringHeader)(unsafe.Pointer(&v))
			typ = reflect.TypeOf(v)
		case reflect.Slice:
			v = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
			typ = reflect.TypeOf(v)
		}
	}
	return v, typ
}

func safeInterface(v reflect.Value) interface{} {
	if v.CanInterface() {
		return v.Interface()
	}
	switch v.Type().Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Interface:
		return reflect.New(v.Type()).Interface()
	case reflect.String:
		return v.String()
	case reflect.Pointer:
		return v.Pointer()
	}
	panic(fmt.Sprintf("unsupported kind %v", v.Type().Kind()))
}

func structDefString(v interface{}) string {
	x := reflect.ValueOf(v).Type()

	sb := strings.Builder{}
	sb.WriteString(x.Name())

	sb.WriteString("{\n\t")
	for i := 0; i < x.NumField(); i++ {
		if i > 0 {
			sb.WriteString("\n\t")
		}
		f := x.Field(i)
		sb.WriteString(fmt.Sprintf("%s\t%s", f.Name, f.Type))
	}
	sb.WriteRune('\n')
	for i := 0; i < x.NumMethod(); i++ {
		sb.WriteString("\n\t")
		f := x.Method(i)
		sb.WriteString(fmt.Sprintf("%s\t%s", f.Name, f.Type))
	}
	sb.WriteString("\n}")

	return sb.String()
}

func pad(x interface{}, min int) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v", x))
	for i := sb.Len(); i < min; i += 1 {
		sb.WriteRune(' ')
	}
	return sb.String()
}

func center(s, padding string, t int) string {
	if s != "" {
		s = " " + s + " "
	}
	fill := t - len(s)
	if fill < 0 {
		return s[:t]
	}
	rep := fill / len(padding)
	b := &strings.Builder{}
	for i := 0; i < rep/2; i += 1 {
		_, _ = fmt.Fprint(b, padding)
	}
	b.WriteString(s)
	for i := 0; i < rep/2+rep%2; i += 1 {
		_, _ = fmt.Fprint(b, padding)
	}
	out := b.String()
	if len(out) > t {
		return out[:t]
	}
	return out
}
