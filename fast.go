
		// +build ignore

		package main

		import (
			"sync/atomic"

			"github.com/erikdubbelboer/plugintest/data"
		)

		var Perf [3]int64

		func Handle(d data.Data) bool {
	
			if d.Check(0, 100) {
				atomic.AddInt64(&Perf[2], 1)
				return false
			}
		
			if d.Check(100, 200) {
				atomic.AddInt64(&Perf[1], 1)
				return false
			}
		
			if d.Check(200, 300) {
				atomic.AddInt64(&Perf[0], 1)
				return false
			}
		
			return true
		}
	