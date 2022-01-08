package stream

import (
	"reflect"
	"testing"
)

func Test_Conv(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want Item
	}{
		{name: "int", arg: 1024, want: i64(1024)},
		{name: "string", arg: "xxx", want: str("xxx")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Conv(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("from() = %v, want %v", got, tt.want)
			}
		})
	}
}
