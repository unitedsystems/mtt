package engine

import (
	"sync"
)

// Descriptor is a link between chat room and client
type Descriptor struct {
	sync.Mutex
	r          *room
	c          *Client
	lastID     int
	name       string
	updateChan chan struct{}
}

func (d *Descriptor) pull(t int64) []Message {
	d.r.RLock()
	defer d.r.RUnlock()
	if d.lastID == d.r.lastID {
		return make([]Message, 0)
	}
	span := d.r.lastID - d.lastID
	if span > historySize {
		span = historySize
	}
	result := make([]Message, span)
	for i := 0; i < span; i++ {
		result[i] = d.r.messages[(d.r.lastID-span+i)%historySize]
	}
	d.lastID = d.r.lastID
	return result
}
