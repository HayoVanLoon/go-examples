package main

import "fmt"

func func1() (ret int) {
	defer func() {
		fmt.Println("75...")
		ret = 75
	}()
	fmt.Println("74...")
	return 74
}

func defers() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("defers", "=", 80))
	fmt.Println(center("", "=", 80))

	x := func1()
	fmt.Println("...", x)
}

func tryButOhNoes() error {
	// this function could for instance be File.Close.
	return fmt.Errorf("oh noes")
}

func crashOriginalLost() (err error) {
	defer func() {
		// this will overwrite the initial error
		err = tryButOhNoes()
	}()

	err = fmt.Errorf("crashOriginalLost")
	return
}

func crash() (err error) {
	defer func() {
		// this will give priority to the initial error, ignoring the latter
		if e := tryButOhNoes(); err == nil && e != nil {
			err = e
		}
	}()

	err = fmt.Errorf("crash")
	return
}

type StackedError interface {
	error
	Previous() error
}

type stackedError struct {
	error
	prev error
}

func (e stackedError) Previous() error {
	return e.prev
}

func NewStackedError(err, previous error) StackedError {
	return &stackedError{err, previous}
}

func crashKeepBoth() (err error) {
	defer func() {
		// keep both at the cost of a bit more documentation mentioning the
		// special error type's capabilities
		if e := tryButOhNoes(); e != nil {
			if err == nil {
				err = e
			} else {
				err = NewStackedError(e, err)
			}
		}
	}()

	err = fmt.Errorf("crash")
	return
}

func deferIntercept() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("deferIntercept", "=", 80))
	fmt.Println(center("", "=", 80))

	fmt.Println(crashOriginalLost())
	fmt.Println(crash())
	err := crashKeepBoth()
	fmt.Println(err.(StackedError).Previous(), err)
}
