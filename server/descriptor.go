package server

import (
	"sync"
)

type clientDescriptor struct {
	sync.Mutex
	r          *room
	c          *client
	lastID     int
	name       string
	updateChan chan struct{}
}

func (d *clientDescriptor) pull(t int64) []message {
	d.r.RLock()
	defer d.r.RUnlock()
	if d.lastID == d.r.lastID {
		return make([]message, 0)
	}
	span := d.r.lastID - d.lastID
	if span > historySize {
		span = historySize
	}
	result := make([]message, span)
	for i := 0; i < span; i++ {
		result[i] = d.r.messages[(d.r.lastID-span+i)%historySize]
	}
	d.lastID = d.r.lastID
	return result
}
