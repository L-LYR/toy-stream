package stream

type Stream interface {
	FilterBy(func(Item) bool, ...Option) Stream
	MapBy(func(Item) Item, ...Option) Stream
	GroupBy(func(Item) interface{}, ...Option) Stream
	SortBy(func(Item, Item) bool) Stream
	Flatten() Stream
	Distinct() Stream

	First() Item
	Last() Item

	Collect() ItemSlice

	Done() Stream
	Result() []error
}

var _ Stream = &_StreamImpl{}

func Range(c <-chan Item) Stream {
	return &_StreamImpl{source: c}
}

func Just(values ...Item) Stream {
	source := make(chan Item, len(values))
	for _, value := range values {
		source <- value
	}
	close(source)
	return Range(source)
}

func GenerateBy(generator func(chan<- Item)) Stream {
	source := make(chan Item)
	Go(func() {
		generator(source)
		close(source) // remember to close !
	})
	return Range(source)
}
