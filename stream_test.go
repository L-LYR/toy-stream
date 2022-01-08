package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJust(t *testing.T) {
	as := assert.New(t)
	values := []interface{}{1, 2, 3, 4, 5}
	s := Just(values...)
	si, ok := s.(*_StreamImpl)
	as.True(ok)
	as.Equal(5, len(si.source))
	for _, value := range values {
		as.Equal(si.fetch(), Conv(value))
	}
}

func TestGenerateBy(t *testing.T) {
	as := assert.New(t)
	g := func(c chan<- Item) {
		for i := 0; i < 5; i++ {
			c <- Conv(i)
		}
	}
	s := GenerateBy(g)
	for _, item := range []i64{0, 1, 2, 3, 4} {
		as.Equal(s.(*_StreamImpl).fetch(), item)
	}
}
