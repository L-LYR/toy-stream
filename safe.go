package stream

import "log"

type Task func()

func Recover(cleanups ...Task) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if r := recover(); r != nil {
		log.Print(r)
	}
}

func Run(fn Task) {
	fn()
	Recover()
}

func Go(fn Task) {
	go Run(fn)
}
