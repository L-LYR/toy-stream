package stream

// WARNING: c must be closed before calling this function
func drain(c <-chan Item) {
	for range c {
		// do nothing
	}
}
