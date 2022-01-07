package stream

type Generator func(chan<- Item)

type Predictor func(Item) bool

type TaskWrapper func(Item, chan<- Item)
