package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/erikdubbelboer/go-plugintest/data"
)

var (
	Perf *[8]int64
)

func getFastHandler() func(d data.Data) bool {
	if Perf == nil {
		p := [8]int64{}
		Perf = &p
	}

	conditions := []string{
		`
			if d.Check(10, 11) {
				return false
			}
		`,
		`
			if d.Check(5, 8) {
				return false
			}
		`,
		`
			if d.Check(150, 360) {
				return false
			}
		`,
		`
			if d.Check(10, 20) {
				return false
			}
		`,
		`
			if d.Check(60, 150) {
				return false
			}
		`,
		`
			if d.Check(2, 6) {
				return false
			}
		`,
		`
			if d.Check(100, 2000) {
				return false
			}
		`,
		`
			if d.Check(0, 4) {
				return false
			}
		`,
	}

	indices := make([]int, len(conditions))
	for i, _ := range indices {
		indices[i] = i
	}

	sort.Slice(indices, func(i, j int) bool {
		return atomic.LoadInt64(&Perf[indices[i]]) > atomic.LoadInt64(&Perf[indices[j]])
	})

	fmt.Println(*Perf)
	fmt.Println(indices)

	fast := `
		// +build ignore

		package main

		import (
			"sync/atomic"

			"github.com/erikdubbelboer/go-plugintest/data"
		)

		// If we don't introduce a new type "plugin" will return "plugin already loaded" once in a while.
		type X` + strconv.FormatInt(time.Now().UnixNano(), 16) + `x int

		var Perf [8]int64

		func Handle(d data.Data) bool {
	`
	x := make([]int64, 0)
	for _, i := range indices {
		code := strings.SplitN(strings.TrimSpace(conditions[i]), "\n", 2)
		fast += code[0]
		fast += "\natomic.AddInt64(&Perf[" + strconv.Itoa(i) + "], 1)\n"
		fast += code[1] + "\n"

		x = append(x, Perf[i])
	}

	fast += `
			return true
		}
	`

	formatted, err := format.Source([]byte(fast))
	if err != nil {
		panic(err)
	}

	filename := "fast-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".so"

	if err := ioutil.WriteFile("fast.go", formatted, 0644); err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", filename, "fast.go")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	p, err := plugin.Open(filename)
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
	Perf = pf.(*[8]int64)

	if err := os.Remove(filename); err != nil {
		panic(err)
	}

	return h.(func(d data.Data) bool)
}
