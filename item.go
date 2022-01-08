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

type ItemSlice []Item

func (l ItemSlice) EqualWith(r Order) bool {
	if len(l) != len(r.(ItemSlice)) {
		return false
	}
	for idx, item := range l {
		if !item.EqualWith(r.(ItemSlice)[idx]) {
			return false
		}
	}
	return true
}

func (l ItemSlice) LessThan(r Order) bool {
	if len(l) != len(r.(ItemSlice)) {
		return false
	}
	for idx, item := range l {
		if !item.LessThan(r.(ItemSlice)[idx]) {
			return false
		}
	}
	return true
}
