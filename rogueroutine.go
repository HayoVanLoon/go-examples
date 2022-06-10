package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

var stdout = make(chan string)

// MotorCar is an embedded interface
type MotorCar interface {
	Car
	Run(context.Context, chan<- string)
}

// MotorPolo is a type of Car
type MotorPolo struct {
	Polo
}

func (v MotorPolo) Run(ctx context.Context, lane chan<- string) {
	defer close(lane)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.NewTimer(120 * time.Second)
	defer deadline.Stop()
	var out chan<- string
	for {
		select {
		case <-ticker.C:
			out = lane
		case out <- "p":
			out = nil
		case <-deadline.C:
			return
		case <-ctx.Done():
			stdout <- fmt.Sprintf("[ polo stopped: %s ]", ctx.Err())
			return
		}
	}
}

// MotorCaddy is a type of Car
type MotorCaddy struct {
	Caddy
	cargo interface{}
}

func (c *MotorCaddy) Run(ctx context.Context, lane chan<- string) {
	defer func() {
		stdout <- fmt.Sprintf("[ caddy stopped: %v ]", ctx.Err())
		close(lane)
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	arrival := time.NewTimer(2 * time.Second)
	defer arrival.Stop()

	var out chan<- string
	for {
		select {
		case <-ticker.C:
			out = lane
		case out <- "c":
			out = nil
		case <-arrival.C:
			return
		case <-ctx.Done():
			panic("should not reach this, should stop before timeout")
		}
	}
}

type UnstoppableCar struct {
}

func (u UnstoppableCar) Move() {
	fmt.Println("*creak*")
}

func (u UnstoppableCar) Honk() {
	stdout <- "<< wub >>"
}

func (u UnstoppableCar) Run(ctx context.Context, lane chan<- string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	h := make(chan bool)
	go func() {
		for range h {
			u.Honk()
		}
	}()

	var honk chan bool
	var out chan<- string
	var done bool
	for {
		select {
		case <-ticker.C:
			if !done {
				out = lane
			} else {
				honk = h
			}
		case <-ctx.Done():
			done = true
		case honk <- true:
			honk = nil
		case out <- "x":
			out = nil
		}
	}
}

func pipeTo[T any](ctx context.Context, in <-chan T, out chan<- T, onClose func()) {
	defer onClose()
	for {
		select {
		case <-ctx.Done():
			return
		case x, ok := <-in:
			if !ok {
				return
			}
			out <- x
		default:
		}
	}
}

const (
	delayUntilCancel = 8 * time.Second
	delayAfterCancel = 3 * time.Second
	delayBrickWall   = delayUntilCancel + delayAfterCancel + 1*time.Second
)

func rogueRountine() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("rogueRountine", "=", 80))
	fmt.Println(center("", "=", 80))

	cars := []MotorCar{
		MotorPolo{},
		&MotorCaddy{},
		UnstoppableCar{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		fmt.Println("\ncanceled output channel")
		fmt.Println("running routines: ", runtime.NumGoroutine())
	}()

	go func() {
		b := strings.Builder{}
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-stdout:
				b.WriteString(s)
				fmt.Print("\r" + b.String())
			}
		}
	}()

	var wg sync.WaitGroup
	var cancels []context.CancelFunc

	for _, c := range cars {
		carCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
		wg.Add(1)
		cancels = append(cancels, func() {
			wg.Done()
			cancel()
		})

		lane := make(chan string)
		go pipeTo(ctx, lane, stdout, func() {})
		go c.Run(carCtx, lane)
	}

	defer func() {
		wg.Wait()
		stdout <- "[ wait group done ]"
	}()

	stdout <- "[ cancel in 8 seconds ]"
	time.Sleep(delayUntilCancel)

	for _, c := range cancels {
		c()
	}
	stdout <- "[ canceled contexts ]"

	time.Sleep(delayAfterCancel)
}
