package data

type Data struct {
	A int
}

var Sink int

//go:noinline
func (d Data) Check(a, b int) bool {
	for i := 0; i < 10000; i++ {
		Sink += d.A
	}
	return d.A > a && d.A < b
}
