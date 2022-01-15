package stream

import "sync"

const (
	DEFAULT_WORKER_NUM = 4
)

type _TaskConfig struct {
	workerNum int
}

func (c *_TaskConfig) check() {
	c.bePositiveOrDefault()
}

func (c *_TaskConfig) bePositiveOrDefault() {
	if c.workerNum <= 0 {
		c.workerNum = DEFAULT_WORKER_NUM
	}
}

type Option func(*_TaskConfig)

func WithWorkerNum(n int) Option {
	return func(c *_TaskConfig) {
		c.workerNum = n
	}
}

type _TaskRunner struct {
	result  chan Item
	limiter chan struct{}
	_TaskConfig
}

func _NewTaskRunner(opts ...Option) *_TaskRunner {
	t := &_TaskRunner{}
	for _, opt := range opts {
		if opt != nil {
			opt(&t._TaskConfig)
		}
	}
	t.check()
	t.result = make(chan Item, t.workerNum)
	t.limiter = make(chan struct{}, t.workerNum)
	return t
}

func (t *_TaskRunner) Apply(fn func(Item, chan<- Item), source <-chan Item) Stream {
	go func() {
		wg := &sync.WaitGroup{}
		t.limiter <- struct{}{}
		for item := range source {
			item := item // awful!
			wg.Add(1)
			GoWithRecover(func() {
				fn(item, t.result)
				wg.Done()
				<-t.limiter
			})
		}
		wg.Wait()
		close(t.result)
	}()
	return Range(t.result)
}

func (t *_TaskRunner) Invoke(fns ...func(chan<- Item)) Stream {
	go func() {
		wg := &sync.WaitGroup{}
		t.limiter <- struct{}{}
		for _, fn := range fns {
			fn := fn
			wg.Add(1)
			GoWithRecover(func() {
				fn(t.result)
				wg.Done()
				<-t.limiter
			})
		}
		wg.Wait()
		close(t.result)
	}()
	return Range(t.result)
}