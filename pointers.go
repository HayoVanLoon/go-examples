package main

import (
	"fmt"
)

type MyError struct {
	code int
}

func (m MyError) Error() string {
	return fmt.Sprintf("error %d", m.code)
}

func typedNilIfPositive(numbers ...int) error {
	var err *MyError
	for _, i := range numbers {
		if i < 0 {
			err = &MyError{42}
		}
		fmt.Printf("%b...", i)
	}
	return err
}

func nilIfPositive(numbers ...int) error {
	var err *MyError
	for _, i := range numbers {
		if i < 0 {
			err = &MyError{42}
		}
		fmt.Printf("%b...", i)
	}
	if err == nil {
		return nil
	}
	return err
}

func structPointerIfNegative(numbers ...int) *MyError {
	var err *MyError
	for _, i := range numbers {
		if i < 0 {
			err = &MyError{42}
		}
		fmt.Printf("%b...", i)
	}
	return err
}

func pointerReturns() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("pointerReturns", "=", 80))
	fmt.Println(center("", "=", 80))

	{
		err := nilIfPositive(-1, 2, 3)
		fmt.Println("\nnilIfPositive(-1, 2, 3): err == nil (expect false):", err == nil)
		err = typedNilIfPositive(-1, 2, 3)
		fmt.Println("\ntypedNilIfPositive(-1, 2, 3): err == nil (expect false):", err == nil)
		err = structPointerIfNegative(-1, 2, 3)
		fmt.Println("\nstructPointerIfNegative(-1, 2, 3): err == nil (expect false):", err == nil)
	}
	fmt.Println()
	{
		err := nilIfPositive(1, 2, 3)
		fmt.Println("\nnilIfPositive(1, 2, 3): err == nil (expect true):", err == nil)
		err = typedNilIfPositive(1, 2, 3)
		fmt.Println("\ntypedNilIfPositive(1, 2, 3): err == nil (expect true):", err == nil)
		fmt.Println("typedNilIfPositive(1, 2, 3): err == (*MyError)(nil) (expect true):", err == (*MyError)(nil))
		err = structPointerIfNegative(1, 2, 3)
		fmt.Println("\nstructPointerIfNegative(1, 2, 3): err == nil (expect true):", err == nil)
		fmt.Println("structPointerIfNegative(1, 2, 3): err == (*MyError)(nil) (expect true):", err == (*MyError)(nil))
	}
}
