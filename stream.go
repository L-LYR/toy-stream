package stream

type Stream interface {
	
	Filter(Predictor) Stream

	First() Item
	Last() Item

	Done() Stream
	Result() []error
}

var _ Stream = &_StreamImpl{}

func Range(c <-chan Item) Stream {
	return &_StreamImpl{source: c}
}

func Just(values ...interface{}) Stream {
	source := make(chan Item, len(values))
	for _, value := range values {
		source <- From(value)
	}
	close(source)
	return Range(source)
}

func GenerateBy(generator Generator) Stream {
	source := make(chan Item)
	Go(func() {
		generator(source)
		close(source) // remember to close !
	})
	return Range(source)
}
