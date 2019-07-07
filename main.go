package main

import (
	"io/ioutil"
	"os/exec"
	"plugin"
	"sort"
	"sync/atomic"

	"github.com/erikdubbelboer/plugintest/data"
)

var (
	Perf *[3]int64
)

//go:noinline
func slowHandler(d data.Data) bool {
	if d.Check(200, 300) {
		return false
	}

	if d.Check(100, 200) {
		return false
	}

	if d.Check(0, 100) {
		return false
	}

	return true
}

func getFastHandler() func(d data.Data) bool {
	if Perf == nil {
		p := [3]int64{}
		Perf = &p
	}

	conditions := []string{
		`
			if d.Check(200, 300) {
				atomic.AddInt64(&Perf[0], 1)
				return false
			}
		`,
		`
			if d.Check(100, 200) {
				atomic.AddInt64(&Perf[1], 1)
				return false
			}
		`,
		`
			if d.Check(0, 100) {
				atomic.AddInt64(&Perf[2], 1)
				return false
			}
		`,
	}

	indices := []int{0, 1, 2}

	sort.Slice(indices, func(i, j int) bool {
		return atomic.LoadInt64(&Perf[indices[i]]) > atomic.LoadInt64(&Perf[indices[j]])
	})

	fast := `
		// +build ignore

		package main

		import (
			"sync/atomic"

			"github.com/erikdubbelboer/plugintest/data"
		)

		var Perf [3]int64

		func Handle(d data.Data) bool {
	`
	for _, i := range indices {
		fast += conditions[i]
	}
	fast += `
			return true
		}
	`

	if err := ioutil.WriteFile("fast.go", []byte(fast), 0644); err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "fast.so", "fast.go")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	p, err := plugin.Open("fast.so")
	if err != nil {
		panic(err)
	}
	h, err := p.Lookup("Handle")
	if err != nil {
		panic(err)
	}
	pf, err := p.Lookup("Perf")
	if err != nil {
		panic(err)
	}
	Perf = pf.(*[3]int64)

	return h.(func(d data.Data) bool)
}
