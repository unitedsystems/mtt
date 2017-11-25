package engine

import (
	"sync"
)

// Descriptor is a link between chat room and client
type Descriptor struct {
	sync.Mutex
	r      *room
	c      *Client
	lastID int
	name   string
}

func (d *Descriptor) pull(t int64) []Message {
	d.r.RLock()
	defer d.r.RUnlock()
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
