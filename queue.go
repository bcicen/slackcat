package main

import (
	"sync"
)

//Streaming message queue
type StreamQ struct {
	lines  []string
	unAckd []string // lines being processed
	lock   sync.RWMutex
}

func newStreamQ() *StreamQ {
	return &StreamQ{
		lines:  []string{},
		unAckd: []string{},
		lock:   sync.RWMutex{},
	}
}

func (q *StreamQ) Len() int      { return len(q.lines) + len(q.unAckd) }
func (q *StreamQ) IsEmpty() bool { return q.Len() < 1 }

func (q *StreamQ) Add(line string) {
	q.lock.Lock()
	q.lines = append(q.lines, line)
	q.lock.Unlock()
}

// Flush returns all lines in queue
func (q *StreamQ) Flush() []string {
	q.lock.Lock()
	defer q.lock.Unlock()
	for _, l := range q.lines {
		q.unAckd = append(q.unAckd, l)
	}
	q.lines = []string{}
	return q.unAckd
}

// acknowledge items from last Get() have been processed
func (q *StreamQ) Ack() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.unAckd = []string{}
}
