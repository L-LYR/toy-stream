package stream

type _StreamImpl struct {
	// stream procedure base channel
	source <-chan Item
	// occurred errors in a stream procedure
	errs []error
}

func (s *_StreamImpl) Filter(p Predictor) Stream {
	return _NewTaskRunner().Do(
		func(i Item, c chan<- Item) {
			if p(i) {
				c <- i
			}
		}, s.source)
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

func (s *_StreamImpl) Done() Stream {
	drain(s.source)
	s.errs = make([]error, 0)
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
