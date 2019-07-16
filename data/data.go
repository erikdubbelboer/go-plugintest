package data

import (
	"time"
)

type Data struct {
	A float64
}

var Sink float64

//go:noinline
func (d Data) Check(a, b float64) bool {
	time.Sleep(time.Microsecond)
	return d.A > a && d.A < b
}
