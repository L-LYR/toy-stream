package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJust(t *testing.T) {
	as := assert.New(t)
	values := []Item{1, 2, 3, 4, 5}
	s := Just(values...)
	si, ok := s.(*_StreamImpl)
	as.True(ok)
	as.Equal(5, len(si.source))
	for _, value := range values {
		as.Equal(si.fetch(), value)
	}
}

func TestGenerateBy(t *testing.T) {
	as := assert.New(t)
	g := func(c chan<- Item) {
		for i := 0; i < 5; i++ {
			c <- i
		}
	}
	s := GenerateBy(g)
	for _, item := range []Item{0, 1, 2, 3, 4} {
		as.Equal(s.(*_StreamImpl).fetch(), item)
	}
}
