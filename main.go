package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/erikdubbelboer/go-plugintest/data"
)

var fast atomic.Value

//go:noinline
func slowHandler(d data.Data) bool {
	if d.Check(10, 11) {
		return false
	}

	if d.Check(5, 8) {
		return false
	}

	if d.Check(150, 360) {
		return false
	}

	if d.Check(10, 20) {
		return false
	}

	if d.Check(60, 150) {
		return false
	}

	if d.Check(2, 6) {
		return false
	}

	if d.Check(100, 2000) {
		return false
	}

	if d.Check(0, 4) {
		return false
	}

	return true
}

func BenchmarkSlow(dur time.Duration) float64 {
	rand.Seed(0)

	start := time.Now()
	ops := float64(0)
	for {
		since := time.Since(start)
		if since > dur {
			break
		}
		d := data.Data{A: rand.NormFloat64() * float64(since/(dur/1000))}
		slowHandler(d)
		ops++
	}
	return ops
}

func BenchmarkFast(dur time.Duration) float64 {
	rand.Seed(0)

	start := time.Now()
	ops := float64(0)
	for {
		since := time.Since(start)
		if since > dur {
			break
		}
		d := data.Data{A: rand.NormFloat64() * float64(since/(dur/1000))}
		ff := fast.Load().(func(d data.Data) bool)
		ff(d)
		ops++
	}
	return ops
}

func main() {
	dur := time.Second * 10

	ops := BenchmarkSlow(dur)
	fmt.Println(int(dur/time.Duration(ops)), "ns/op")

	fast.Store(getFastHandler())

	go func() {
		for {
			time.Sleep(time.Second)
			fast.Store(getFastHandler())
		}
	}()

	ops = BenchmarkFast(dur)
	fmt.Println(int(dur/time.Duration(ops)), "ns/op")
}
