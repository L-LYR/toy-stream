package stream

type Generator func(chan<- Item)

type Predictor func(Item) bool

type Transformer func(Item) Item

type TaskWrapper func(Item, chan<- Item)
