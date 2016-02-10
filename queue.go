package main

import (
	"sync"
)

//Streaming message queue
type StreamQ struct {
	queue []string
	lock  sync.RWMutex
}

func newStreamQ() *StreamQ {
	return &StreamQ{
		queue: []string{},
		lock:  sync.RWMutex{},
	}
}

func (q *StreamQ) add(line string) {
	q.lock.Lock()
	q.queue = append(q.queue, line)
	q.lock.Unlock()
}

func (q *StreamQ) isEmpty() bool {
	return (len(q.queue) < 1)
}

func (q *StreamQ) flush() []string {
	q.lock.Lock()
	items := q.queue
	q.queue = []string{}
	q.lock.Unlock()
	return items
}
