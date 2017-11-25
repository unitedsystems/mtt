package server

import (
	"fmt"
	"sync"
	"time"
)

// room refers to chat room
type room struct {
	sync.RWMutex
	messages []message
	lastID   int

	clientLock        sync.RWMutex
	clientDescriptors map[string]*clientDescriptor

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
			select {
			case d.c.masterNotification <- struct{}{}:
			default:
			}
		}
		r.clientLock.RUnlock()
	}
}

// concurrent
func (r *room) Publish(m message) {
	r.Lock()
	m.timestamp = time.Now().UnixNano()
	r.messages[r.lastID%historySize] = m
	r.lastID++
	r.Unlock()

	select {
	case r.broadcastChan <- struct{}{}: // shedules broadcast
	default:
	}
}

func (r *room) Subscribe(name string, c *client) (*clientDescriptor, error) {
	d := &clientDescriptor{
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

func (r *room) Unsubscribe(cd *clientDescriptor) {
	r.clientLock.Lock()
	defer r.clientLock.Unlock()
	delete(r.clientDescriptors, cd.name)
}

func newRoom(name string) *room {
	r := &room{
		messages:          make([]message, historySize),
		clientDescriptors: make(map[string]*clientDescriptor, historySize),
		broadcastChan:     make(chan struct{}, 1),
		clientLock:        sync.RWMutex{},
	}
	go r.broadcast()
	return r
}
