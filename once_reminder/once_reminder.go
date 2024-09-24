package oncereminder

import (
	"time"
)

type OnceReminder struct {
	duration time.Duration
	task     func()
	done     chan bool
}

func NewOnceReminder(duration time.Duration, task func()) *OnceReminder {
	return &OnceReminder{
		duration: duration,
		task:     task,
		done:     make(chan bool),
	}
}

func (or *OnceReminder) Start() {
	go func() {
		select {
		case <-time.After(or.duration):
			or.task()
		case <-or.done:
			return
		}
	}()
}

func (or *OnceReminder) Cancel() {
	close(or.done)
}
