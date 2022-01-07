package stream

type i64 int64

func (l i64) EqualWith(r Order) bool { return l == r.(i64) }
func (l i64) LessThan(r Order) bool  { return l < r.(i64) }

type u64 uint64

func (l u64) EqualWith(r Order) bool { return l == r.(u64) }
func (l u64) LessThan(r Order) bool  { return l < r.(u64) }

type f64 float64

func (l f64) EqualWith(r Order) bool { return l == r.(f64) }
func (l f64) LessThan(r Order) bool  { return l < r.(f64) }

type str string

func (l str) EqualWith(r Order) bool { return l == r.(str) }
func (l str) LessThan(r Order) bool  { return l < r.(str) }
