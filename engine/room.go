package engine

import (
	"fmt"
	"sync"
	"time"
)

// room refers to chat room
type room struct {
	sync.RWMutex
	messages []Message
	lastID   int

	clientLock        sync.RWMutex
	clientDescriptors map[string]*Descriptor

	broadcastChan chan struct{}
}

func (r *room) broadcast() {
	for {
		<-r.broadcastChan
		r.clientLock.RLock()
		for _, d := range r.clientDescriptors {
			select {
			case d.updateChan <- struct{}{}:
			default:
			}
		}
		r.clientLock.RUnlock()
		r.clientLock.RLock()
		for _, d := range r.clientDescriptors {
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
	d := &Descriptor{
		c:          c,
		r:          r,
		name:       name,
		updateChan: make(chan struct{}, 1),
	}
	r.clientLock.Lock()
	defer r.clientLock.Unlock()
	if _, ok := r.clientDescriptors[name]; ok {
		return nil, fmt.Errorf("%s is not unique", name)
	}
	c.descriptors = append(c.descriptors, d)
	r.clientDescriptors[name] = d

	return d, nil
}

func newRoom(name string) *room {
	r := &room{
		messages:          make([]Message, historySize),
		clientDescriptors: make(map[string]*Descriptor, historySize),
		broadcastChan:     make(chan struct{}, 1),
		clientLock:        sync.RWMutex{},
	}
	go r.broadcast()
	return r
}
