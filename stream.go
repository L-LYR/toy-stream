package stream

type Stream interface {
	FilterBy(func(Item) bool, ...Option) Stream
	MapBy(func(Item) Item, ...Option) Stream
	GroupBy(func(Item) Item, ...Option) Stream
	SortBy(func(Item, Item) bool) Stream
	Flatten() Stream
	Distinct() Stream
	Head(int) Stream
	Tail(int) Stream
	Reverse() Stream

	First() Item
	Last() Item

	All(func(Item) bool) bool
	Any(func(Item) bool) bool
	None(func(Item) bool) bool
	Count() int

	Collect() ItemSlice
	Consume() (Item, bool)

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
	GoWithRecover(func() {
		generator(source)
		close(source) // remember to close !
	})
	return Range(source)
}

func Concat(streams ...Stream) Stream {
	fns := []func(chan<- Item){}
	for _, stream := range streams {
		stream := stream // awful!
		fns = append(fns, func(c chan<- Item) {
			for {
				item, isEmpty := stream.Consume()
				if isEmpty {
					break
				}
				c <- item
			}
		})
	}
	return _NewTaskRunner( /*default option*/ ).Invoke(fns...)
}
