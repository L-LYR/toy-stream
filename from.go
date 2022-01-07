package stream

import "reflect"

func From(value interface{}) Item {
	if value == nil {
		return nil
	}
	underlying := reflect.ValueOf(value)
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return i64(underlying.Int())
	case uint, uint8, uint16, uint32, uint64:
		return u64(underlying.Uint())
	case float32, float64:
		return f64(underlying.Float())
	case string:
		return str(underlying.String())
	case Item:
		return v
	default:
		panic("Unsupport!")
	}
}
