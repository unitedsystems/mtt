package engine

import (
	"fmt"
	"sync"
	"time"
)

type room struct {
	sync.RWMutex
	messages []Message
	lastID   int
	name     string

	clientLock  sync.RWMutex
	descriptors map[string]*Descriptor // NOTE: key is user's name for this particular room

	broadcastChan chan struct{}
}

func newRoom(name string) *room {
	r := &room{
		name:          name,
		messages:      make([]Message, historySize),
		descriptors:   make(map[string]*Descriptor, historySize),
		broadcastChan: make(chan struct{}, 1),
		clientLock:    sync.RWMutex{},
	}
	go r.askClientsToPoll()
	return r
}

func (r *room) askClientsToPoll() {
	for {
		<-r.broadcastChan
		r.clientLock.RLock()
		for _, d := range r.descriptors {
			select {
			case d.c.masterNotification <- struct{}{}:
			default:
			}
		}
		r.clientLock.RUnlock()
	}
}

func (r *room) publish(m Message) {
	r.Lock()
	m.Timestamp = time.Now().UnixNano()
	r.messages[r.lastID%historySize] = m
	r.lastID++
	r.Unlock()

	select {
	case r.broadcastChan <- struct{}{}: // shedules broadcast
	default:
	}
}

func (r *room) subscribe(name string, c *Client) (*Descriptor, error) {
	d := &Descriptor{c: c, r: r}
	r.clientLock.Lock()
	if _, ok := r.descriptors[name]; ok {
		r.clientLock.Unlock()
		return nil, fmt.Errorf("%s is not unique", name)
	}
	r.descriptors[name] = d
	r.clientLock.Unlock()
	return d, nil
}
