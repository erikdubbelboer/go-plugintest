package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/erikdubbelboer/plugintest/data"
)

var fast func(d data.Data) bool

//go:noinline
func BenchmarkSlow(n int) {
	rand.Seed(0)

	for i := 0; i < n; i++ {
		d := data.Data{A: rand.Intn(300+150) % 300}
		slowHandler(d)
	}
}

//go:noinline
func BenchmarkFast(n int) {
	rand.Seed(0)

	for i := 0; i < n; i++ {
		d := data.Data{A: rand.Intn(300+150) % 300}
		fast(d)
	}
}

func main() {
	BenchmarkSlow(1000)
	start := time.Now()
	BenchmarkSlow(100000)
	fmt.Println(time.Now().Sub(start))

	fast = getFastHandler()
	BenchmarkFast(1000)
	fast = getFastHandler()
	start = time.Now()
	BenchmarkFast(100000)
	fmt.Println(time.Now().Sub(start))
}
