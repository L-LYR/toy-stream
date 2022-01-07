package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	as := assert.New(t)
	as.EqualValues(1, Just(1, 2, 3, 4, 5).First())
	as.EqualValues(nil, Just().First())
}

func TestLast(t *testing.T) {
	as := assert.New(t)
	as.EqualValues(5, Just(1, 2, 3, 4, 5).Last())
	as.EqualValues(nil, Just().Last())
}

func TestDone(t *testing.T) {
	as := assert.New(t)

	s := Just(1, 2, 3, 4, 5)

	as.EqualValues(1, s.First())
	as.EqualValues(4, len(s.(*_StreamImpl).source))

	as.EqualValues(5, s.Last())
	as.EqualValues(0, len(s.(*_StreamImpl).source))

	// try to fetch an empty channel
	as.EqualValues(nil, s.(*_StreamImpl).fetch())
	as.EqualValues([]error{ErrUnexpectedChannelEnd}, s.Result())

	s = s.Done()
	as.EqualValues(0, len(s.(*_StreamImpl).source))
	as.EqualValues([]error{}, s.Result())
}