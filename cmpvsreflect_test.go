package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// BenchmarkComparisons compares speed of reflect.DeepEqual versus cmp.Equal.
func BenchmarkComparisons(b *testing.B) {
	type Foo struct {
		S string
	}
	type Bar struct {
		S   string
		Foo Foo
	}
	type Bla struct {
		Bars []Bar
	}
	type Vla struct {
		Foo Foo
		Bar Bar
		Bla Bla
	}
	left := Vla{
		Foo: Foo{"apple"},
		Bar: Bar{"banana", Foo{"lemon"}},
		Bla: Bla{
			Bars: []Bar{
				{"meow", Foo{"canary"}},
				{"woof", Foo{"dog"}},
				{"ribbit", Foo{"cat"}},
			},
		},
	}
	right := Vla{
		Foo: Foo{"apple"},
		Bar: Bar{"banana", Foo{"lemon"}},
		Bla: Bla{
			Bars: []Bar{
				{"meow", Foo{"canary"}},
				{"woof", Foo{"dog"}},
				{"ribbit", Foo{"CAT!!1!"}},
			},
		},
	}
	b.Run("cmp", func(b *testing.B) {
		for i := 0; i < b.N; i += 1 {
			_ = cmp.Equal(left, right)
		}
	})

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i += 1 {
			_ = reflect.DeepEqual(left, right)
		}
	})
}
