package stream

import "testing"

func Test_drain(t *testing.T) {
	type args struct {
		c <-chan Item
	}
	type TestCase struct {
		name string
		args args
	}

	tests := func() []TestCase {
		c1 := make(chan Item, 10)
		close(c1)

		c2 := make(chan Item)
		close(c2)

		c3 := make(chan Item, 20000)
		for i := 0; i < 1000; i++ {
			c3 <- Conv(i)
		}
		close(c3)

		return []TestCase{
			{name: "empty closed channel with buffer", args: args{c: c1}},
			{name: "empty closed channel without buffer", args: args{c: c2}},
			{name: "closed channel with buffer", args: args{c: c3}},
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drain(tt.args.c)
		})
	}
}
