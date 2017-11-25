package engine

import (
	"sync"
)

// Link is a link between chat room and client
type Link struct {
	sync.Mutex
	r      *room
	c      *Client
	lastID int
	name   string
}

func (l *Link) pull(t int64) []Message {
	l.r.RLock()
	defer l.r.RUnlock()
	span := l.r.lastID - l.lastID
	if span > historySize {
		span = historySize
	}
	result := make([]Message, span)
	for i := 0; i < span; i++ {
		result[i] = l.r.messages[(l.r.lastID-span+i)%historySize]
	}
	l.lastID = l.r.lastID
	return result
}
