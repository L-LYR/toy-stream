package stream

type Item interface {
	Order
}

type Order interface {
	EqualWith(Order) bool
	LessThan(Order) bool
}

var _ Order = i64(0)
var _ Order = u64(0)
var _ Order = f64(0)
var _ Order = str("")
