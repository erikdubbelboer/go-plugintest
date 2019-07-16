// +build ignore

package main

import (
	"sync/atomic"

	"github.com/erikdubbelboer/go-plugintest/data"
)

// If we don't introduce a new type "plugin" will return "plugin already loaded" once in a while.
type X15b1d100382619b8x int

var Perf [8]int64

func Handle(d data.Data) bool {
	if d.Check(100, 2000) {
		atomic.AddInt64(&Perf[6], 1)
		return false
	}
	if d.Check(60, 150) {
		atomic.AddInt64(&Perf[4], 1)
		return false
	}
	if d.Check(10, 20) {
		atomic.AddInt64(&Perf[3], 1)
		return false
	}
	if d.Check(2, 6) {
		atomic.AddInt64(&Perf[5], 1)
		return false
	}
	if d.Check(0, 4) {
		atomic.AddInt64(&Perf[7], 1)
		return false
	}
	if d.Check(5, 8) {
		atomic.AddInt64(&Perf[1], 1)
		return false
	}
	if d.Check(150, 360) {
		atomic.AddInt64(&Perf[2], 1)
		return false
	}
	if d.Check(10, 11) {
		atomic.AddInt64(&Perf[0], 1)
		return false
	}

	return true
}
