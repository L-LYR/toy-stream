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

func RunWithRecover(fn Task) {
	fn()
	Recover()
}

func GoWithRecover(fn Task) {
	go RunWithRecover(fn)
}
