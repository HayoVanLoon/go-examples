package main

import (
	"context"
	"fmt"
	"time"
)

// Vehicle is a basic interface
type Vehicle interface {
	Move()
}

// Car is an embedded interface
type Car interface {
	Vehicle
	Run(context.Context, chan<- string)
	Honk()
}

// Polo is a type of Car
type Polo struct {
}

func (v Polo) Move() {
	fmt.Println("Vroom.")
}

func (v Polo) Honk() {
	fmt.Println("Honk.")
}

func (v Polo) Run(ctx context.Context, lane chan<- string) {
	defer close(lane)
	step := 500 * time.Millisecond
	ch := time.Tick(step)
	for i := time.Duration(0); i < 120*time.Second; i += step {
		select {
		case <-ctx.Done():
			fmt.Println("\npolo", ctx.Err())
			return
		default:
		}
		<-ch
		select {
		case lane <- "c":
		default:
		}
	}
	panic("should not reach this, should not due to timeout")
}

// Caddy is a type of Car
type Caddy struct {
	cargo interface{}
}

func (c *Caddy) Move() {
	fmt.Println("Vroom.")
}

func (c *Caddy) Honk() {
	fmt.Println("Honk.")
}

func (c *Caddy) Run(ctx context.Context, lane chan<- string) {
	defer close(lane)
	step := 500 * time.Millisecond
	ch := time.Tick(step)
	for i := time.Duration(0); i < 2*time.Second; i += step {
		<-ch
		select {
		case <-ctx.Done():
			panic("should not reach this, should stop before timeout")
		case lane <- "c":
		default:
		}
	}
	fmt.Println("\ncaddy", ctx.Err())
}

func (c *Caddy) Load(v interface{}) error {
	if c.cargo != nil {
		return fmt.Errorf("unload first")
	}
	c.cargo = v
	return nil
}

func (c *Caddy) Unload() (interface{}, error) {
	if c.cargo == nil {
		return nil, fmt.Errorf("nothing loaded")
	}
	v := c.cargo
	c.cargo = nil
	return v, nil
}

func CreateCars() []Car {
	return []Car{
		Polo{},
		&Caddy{},
	}
}

func interfaces() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("interfaces", "=", 80))
	fmt.Println(center("", "=", 80))

	cars := CreateCars()

	for _, c := range cars {
		c.Move()
		c.Honk()
		switch x := c.(type) {
		case *Caddy:
			if err := x.Load("foo"); err != nil {
				panic(err)
			}
			if _, err := x.Unload(); err != nil {
				panic(err)
			}
			if err := x.Load("bar"); err != nil {
				panic(err)
			}
		}
	}
}
