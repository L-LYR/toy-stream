package stream

type Item interface {
	Order
}

type Order interface {
	EqualWith(Order) bool
	LessThan(Order) bool
}

var _ Item = i64(0)
var _ Item = u64(0)
var _ Item = f64(0)
var _ Item = str("")
