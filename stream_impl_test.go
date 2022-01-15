package stream

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	as := assert.New(t)
	as.EqualValues(1, Just(1, 2, 3, 4, 5).First())
	as.EqualValues(nil, Just().First())
}

func TestLast(t *testing.T) {
	as := assert.New(t)
	as.EqualValues(5, Just(1, 2, 3, 4, 5).Last())
	as.EqualValues(nil, Just().Last())
}

func TestDone(t *testing.T) {
	as := assert.New(t)

	s := Just(1, 2, 3, 4, 5)

	as.EqualValues(1, s.First())
	as.EqualValues(4, len(s.(*_StreamImpl).source))

	as.EqualValues(5, s.Last())
	as.EqualValues(0, len(s.(*_StreamImpl).source))

	// try to fetch an empty channel
	as.EqualValues(nil, s.(*_StreamImpl).fetch())
	as.EqualValues([]error{ErrUnexpectedChannelEnd}, s.Result())

	s = s.Done()
	as.EqualValues(0, len(s.(*_StreamImpl).source))
	as.Nil(s.Result())
}

func TestFilterBy(t *testing.T) {
	as := assert.New(t)

	s := Just(1, 2, 3, 4, 5, 6, 7, 8)
	result := s.FilterBy(func(i Item) bool {
		return i.(int)%2 == 0
	}).SortBy(func(i1, i2 Item) bool {
		return i1.(int) < i2.(int)
	}).Collect()
	for i, v := range []int{2, 4, 6, 8} {
		as.EqualValues(v, result[i])
	}
	as.Nil(s.Result())
	s.Done()
}

func TestMapBy(t *testing.T) {
	as := assert.New(t)

	s := Just(1, 2, 3, 4)
	result := s.MapBy(func(i Item) Item {
		return i.(int) * 2
	}).SortBy(func(i1, i2 Item) bool {
		return i1.(int) < i2.(int)
	}).Collect()
	for i, v := range []int{2, 4, 6, 8} {
		as.EqualValues(v, result[i])
	}
	as.Nil(s.Result())
	s.Done()
}

func TestSortBy(t *testing.T) {
	as := assert.New(t)

	type Complex struct {
		s string
		i int
		f float32
	}
	cs := []Complex{
		{s: "a", i: 498, f: 1.02134},
		{s: "asbdj", i: 1, f: 1.3240},
		{s: "dssfjaj", i: 43645, f: 1435.0},
		{s: "a", i: 9999, f: 1.0234},
		{s: "kskdll", i: 5555, f: 1657.0},
	}
	result := GenerateBy(func(c chan<- Item) {
		for _, s := range cs {
			c <- s
		}
	}).SortBy(func(i1, i2 Item) bool {
		return i1.(Complex).s < i2.(Complex).s
	}).Collect()

	sort.Slice(cs, func(i, j int) bool {
		return cs[i].s < cs[j].s
	})
	for i, c := range result {
		as.EqualValues(cs[i], c.(Complex))
	}
}

func TestGroupBy(t *testing.T) {
	as := assert.New(t)

	result := Just(1, 2, 3, 4, 5, 6, 7, 8).GroupBy(func(i Item) Item {
		return i.(int)%2 == 0
	}).GroupBy(func(i Item) Item {
		return len(i.(ItemSlice))
	}).Flatten().Flatten().SortBy(func(i1, i2 Item) bool {
		return i1.(int) < i2.(int)
	}).Collect()

	as.EqualValues(ItemSlice{1, 2, 3, 4, 5, 6, 7, 8}, result)
}

func TestDistinct(t *testing.T) {
	as := assert.New(t)
	result := Just(1, 5, 1, 2, 3, 4, 2, 3, 4, 5).Distinct().SortBy(
		func(i1, i2 Item) bool { return i1.(int) < i2.(int) },
	).Collect()
	as.EqualValues(ItemSlice{1, 2, 3, 4, 5}, result)
}

func TestHead(t *testing.T) {
	as := assert.New(t)
	result := Just(1, 2, 3, 4).Head(10).Collect()
	as.EqualValues(ItemSlice{1, 2, 3, 4}, result)
	result = Just(1, 2, 3, 4, 5, 6, 7).Head(4).Collect()
	as.EqualValues(ItemSlice{1, 2, 3, 4}, result)
}

func TestTail(t *testing.T) {
	as := assert.New(t)
	result := Just(1, 2, 3, 4).Tail(10).Collect()
	as.EqualValues(ItemSlice{1, 2, 3, 4}, result)
	result = Just(1, 2, 3, 4, 5, 6, 7, 8, 9).Tail(4).Collect()
	as.EqualValues(ItemSlice{6, 7, 8, 9}, result)
}

func TestReverse(t *testing.T) {
	as := assert.New(t)
	result := Just(1, 2, 3, 4).Reverse().Collect()
	as.EqualValues(ItemSlice{4, 3, 2, 1}, result)
}

func TestAll(t *testing.T) {
	as := assert.New(t)
	isOdd := func(i Item) bool { return i.(int)%2 == 1 }
	as.True(Just(1, 3, 5, 7).All(isOdd))
}

func TestAny(t *testing.T) {
	as := assert.New(t)
	isEven := func(i Item) bool { return i.(int)%2 == 0 }
	as.True(Just(1, 2, 3, 5, 7).Any(isEven))
}
