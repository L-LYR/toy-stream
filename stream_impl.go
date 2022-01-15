package stream

import (
	"reflect"
	"sort"
)

type _StreamImpl struct {
	// stream procedure base channel
	source <-chan Item
	// occurred errors in a stream procedure
	errs []error
}

func (s *_StreamImpl) apply(fn func(Item, chan<- Item), opts ...Option) Stream {
	return _NewTaskRunner(opts...).Apply(fn, s.source)
}

func (s *_StreamImpl) FilterBy(p func(Item) bool, opts ...Option) Stream {
	return s.apply(func(i Item, c chan<- Item) {
		if p(i) {
			c <- i
		}
	}, opts...)
}

func (s *_StreamImpl) MapBy(t func(Item) Item, opts ...Option) Stream {
	return s.apply(func(i Item, c chan<- Item) {
		c <- t(i)
	}, opts...)
}

func (s *_StreamImpl) GroupBy(k func(Item) Item, opts ...Option) Stream {
	groups := make(map[Item][]Item)
	for item := range s.source {
		key := k(item)
		groups[key] = append(groups[key], item)
	}
	return GenerateBy(func(c chan<- Item) {
		for _, group := range groups {
			c <- ItemSlice(group)
		}
	})
}

func (s *_StreamImpl) Flatten() Stream {
	return GenerateBy(func(c chan<- Item) {
		for item := range s.source {
			if reflect.TypeOf(item) == reflect.TypeOf(ItemSlice{}) {
				for _, innerItem := range item.(ItemSlice) {
					c <- innerItem
				}
			} else {
				c <- item
			}
		}
	})
}

func (s *_StreamImpl) Distinct() Stream {
	return GenerateBy(func(c chan<- Item) {
		filter := make(map[Item]struct{})
		for item := range s.source {
			if _, ok := filter[item]; !ok {
				c <- item
				filter[item] = struct{}{}
			}
		}
	})
}

func (s *_StreamImpl) Head(n int) Stream {
	if n < 1 {
		panic("n must be greater than 0")
	}
	return GenerateBy(func(c chan<- Item) {
		for item := range s.source {
			c <- item
			n--
			if n == 0 {
				break
			}
		}
	})
}

func (s *_StreamImpl) Tail(n int) Stream {
	ring := make([]Item, n)
	idx := 0
	return GenerateBy(func(c chan<- Item) {
		for item := range s.source {
			if idx >= n {
				ring[idx%n] = item
			} else {
				ring[idx] = item
			}
			idx++
		}
		if idx < n {
			for i := 0; i < idx; i++ {
				c <- ring[i]
			}
		} else {
			base := idx % n
			for i := 0; i < n; i++ {
				c <- ring[(base+i)%n]
			}
		}
	})
}

func (s *_StreamImpl) Reverse() Stream {
	buffer := ItemSlice{}
	for item := range s.source {
		buffer = append(buffer, item)
	}
	return GenerateBy(func(c chan<- Item) {
		for i := len(buffer) - 1; i >= 0; i-- {
			c <- buffer[i]
		}
	})
}

func (s *_StreamImpl) All(fn func(Item) bool) bool {
	for item := range s.source {
		if !fn(item) {
			go drain(s.source)
			return false
		}
	}
	return true
}

func (s *_StreamImpl) Any(fn func(Item) bool) bool {
	for item := range s.source {
		if fn(item) {
			go drain(s.source)
			return true
		}
	}
	return false
}

func (s *_StreamImpl) None(fn func(Item) bool) bool {
	return !s.Any(fn)
}

func (s *_StreamImpl) Count() int {
	cnt := 0
	for range s.source {
		cnt++
	}
	return cnt
}

func (s *_StreamImpl) First() Item {
	if len(s.source) == 0 {
		return nil
	}
	return s.fetch()
}

func (s *_StreamImpl) Last() Item {
	if len(s.source) == 0 {
		return nil
	}
	for len(s.source) > 1 {
		_ = s.fetch()
	}
	return s.fetch()
}

func (s *_StreamImpl) SortBy(cmp func(Item, Item) bool) Stream {
	result := s.Collect()
	sort.Slice(result, func(i, j int) bool {
		return cmp(result[i], result[j])
	})
	return Just(result...)
}

func (s *_StreamImpl) Collect() ItemSlice {
	result := make(ItemSlice, 0, len(s.source))
	for item := range s.source {
		result = append(result, item)
	}
	return result
}

func (s *_StreamImpl) Done() Stream {
	drain(s.source)
	s.errs = nil
	return s
}

func (s *_StreamImpl) Result() []error {
	return s.errs
}

func (s *_StreamImpl) Consume() (Item, bool) {
	item, ok := <-s.source
	return item, !ok
}

// `fetch` will try to fetch an `item` from source channel.
// If `ok` is false, it will put an `UnexpectedChannelEnd` into result `errs`.
func (s *_StreamImpl) fetch() Item {
	item, ok := <-s.source
	if !ok {
		s.errs = append(s.errs, ErrUnexpectedChannelEnd)
		return nil
	}
	return item
}
