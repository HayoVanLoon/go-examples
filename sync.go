package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Ticker interface {
	Ticks() <-chan time.Time
	Stop()
}

type ticker time.Ticker

func (t ticker) Ticks() <-chan time.Time {
	return t.C
}

func (t *ticker) Stop() {
	(*time.Ticker)(t).Stop()
}

func newTicker(d time.Duration) Ticker {
	return (*ticker)(time.NewTicker(d))
}

type IntBuffer struct {
	data []int
	r    int
	w    int
	mux  sync.RWMutex
}

func NewIntBuffer(size int) *IntBuffer {
	return &IntBuffer{data: make([]int, size)}
}

func (b *IntBuffer) Scan() int {
	if b.Empty() {
		return math.MaxInt
	}
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.data[b.r%len(b.data)]
}

func (b *IntBuffer) Read() int {
	if b.Empty() {
		return math.MaxInt
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	i := b.data[b.r%len(b.data)]
	b.r += 1
	return i
}

func (b *IntBuffer) Write(i int) bool {
	if b.Full() {
		return false
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	b.data[b.w%len(b.data)] = i
	if b.w == math.MaxInt {
		// prevent theoretical overflow
		d := b.w - b.r
		b.r = b.r % len(b.data)
		b.w = b.r + d
	}
	b.w += 1
	return true
}

func (b *IntBuffer) Full() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.w-b.r == len(b.data)
}

func (b *IntBuffer) Empty() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.r == b.w
}

func (b *IntBuffer) Sat() int {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.w - b.r
}

func (b *IntBuffer) String() string {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return fmt.Sprintf("IntBuffer(%d){w: %d, r: %d, w-r: %d}", len(b.data), b.w, b.r, b.w-b.r)
}

func produce(ctx context.Context, rate int, random bool) (<-chan int, <-chan int) {
	out := make(chan int)
	total := make(chan int, 1)
	delay := 1000 / rate
	go func() {
		defer close(out)
		ticker := newTicker(time.Millisecond)
		defer ticker.Stop()

		acc := 0
		defer func() {
			total <- acc
			close(total)
		}()

		var out2 chan<- int
		const signal = 1
		cd := delay
		for {
			select {
			case <-ticker.Ticks():
				cd -= 1
				if cd == 0 {
					out2 = out
					cd = delay
					if random {
						cd *= (rand.Intn(9) + 1) / 10
					}
				}
			case out2 <- signal:
				acc += signal
				out2 = nil
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, total
}

func duplicate(in <-chan int, num int, bsize int) ([]chan int, []chan int) {
	chs := make([]chan int, num)
	bufs := make([]*IntBuffer, num)
	bps := make([]chan int, num)
	for i := range chs {
		chs[i] = make(chan int)
		bufs[i] = NewIntBuffer(bsize)
		bps[i] = make(chan int, 1)
	}
	go func() {
		defer func() {
			for i := range chs {
				close(chs[i])
				lost := 0
				for !bufs[i].Empty() {
					lost += bufs[i].Read()
				}
				bps[i] <- lost
				close(bps[i])
			}
		}()
		for x := range in {
			for i, ch := range chs {
				bufs[i].Write(x)
				select {
				case ch <- bufs[i].Scan():
					_ = bufs[i].Read()
				default:
				}
			}
		}
	}()
	return chs, bps
}

func limit(in <-chan int, rate, size int) (<-chan int, <-chan int) {
	out := make(chan int)
	ds := make(chan int, 1)
	go func() {
		defer close(out)
		ticker := newTicker(1 * time.Second)
		defer ticker.Stop()

		buf := make([]int, size)
		w, r := 0, 0

		var out2 chan<- int
		rec := 0
		s := 0
		d := 0
		lim := rate

		defer func() {
			ds <- d
			close(ds)
		}()

		for {
			select {
			case <-ticker.Ticks():
				lim += rate
				fmt.Printf("received/sent/discarded/buffered: %d / %d / %d / %d\n", rec, s, d, w-r)
				if in == nil && s < lim {
					out2 = out
				}
			case i, ok := <-in:
				if !ok {
					in = nil
					fmt.Println("draining", w-r)
				}
				rec += 1
				if w-r < size {
					buf[w%size] = i
					w += 1
				} else {
					d += 1
				}
				if s < lim {
					out2 = out
				}
			case out2 <- buf[r%size]:
				s += 1
				r += 1
				if r == w {
					out2 = nil
					if in == nil {
						return
					}
				}
				if s == lim {
					out2 = nil
				}
			}
		}
	}()
	return out, ds
}

func sum(ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		acc := 0
		defer func() {
			out <- acc
			close(out)
		}()
		for {
			select {
			case i, ok := <-ch:
				if !ok {
					return
				}
				acc += i
			}
		}
	}()
	return out
}

func pipeline() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("pipeline", "=", 80))
	fmt.Println(center("", "=", 80))

	ctx, cancel := context.WithCancel(context.Background())

	const (
		prodRate  = 100
		limitRate = 25
		duration  = 4
	)

	ch, total := produce(ctx, prodRate, false)
	chs, bps := duplicate(ch, 2, 10)

	normal := chs[0]
	limited, discarded := limit(chs[1], limitRate, (prodRate-limitRate)*duration)

	normalSum := sum(normal)
	limitedSum := sum(limited)

	time.Sleep(3 * time.Second)
	fmt.Println("canceling")
	cancel()

	r := struct {
		total     int
		normal    int
		limited   int
		discarded int
		bpNormal  int
		bpLimited int
	}{}
	r.total = <-total
	r.normal = <-normalSum
	r.limited = <-limitedSum
	r.discarded = <-discarded
	r.bpNormal = <-bps[0]
	r.bpLimited = <-bps[1]

	fmt.Printf("total produced:  %d\n", r.total)
	fmt.Printf("total normal:    %d\n", r.normal)
	fmt.Printf("total limited:   %d\n", r.limited)
	fmt.Printf("total discarded: %d (%2.0f %%)\n", r.discarded, 100*float32(r.discarded)/float32(r.total))
	fmt.Printf("backpressure 0:  %d\n", r.bpNormal)
	fmt.Printf("backpressure 1:  %d\n", r.bpLimited)

	fmt.Println("done")
	time.Sleep(duration * time.Second)
}

func sharedMemoryNoMutex() {
	fmt.Println(center("", "=", 80))
	fmt.Println(center("sharedMemoryNoMutex", "=", 80))
	fmt.Println(center("", "=", 80))

	var done bool
	var a int

	inc := func() {
		for !done {
			if a == 0 {
				a += 1
			}
		}
	}

	dec := func() {
		for !done {
			if a > 0 {
				a -= 1
			}
		}
	}

	go inc()
	go dec()
	go dec()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		t := time.Now().Add(1 * time.Second)
		acc := 0
		for time.Now().Before(t) {
			if a < 0 {
				acc += 1
			}
		}
		fmt.Println("a:", a, "\t#(a<0)", acc)
		done = true
		wg.Done()
	}()
	wg.Wait()
}
