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

func (s *_StreamImpl) FilterBy(p func(Item) bool, opts ...Option) Stream {
	return _NewTaskRunner(opts...).Do(
		func(i Item, c chan<- Item) {
			if p(i) {
				c <- i
			}
		}, s.source)
}

func (s *_StreamImpl) MapBy(t func(Item) Item, opts ...Option) Stream {
	return _NewTaskRunner(opts...).Do(
		func(i Item, c chan<- Item) {
			c <- t(i)
		}, s.source)
}

func (s *_StreamImpl) GroupBy(k func(Item) interface{}, opts ...Option) Stream {
	groups := make(map[interface{}][]Item)
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
