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

	clientLock sync.RWMutex
	links      map[string]*Link // NOTE: key is user's name for this particular room

	broadcastChan chan struct{}
}

func newRoom(name string) *room {
	r := &room{
		name:          name,
		messages:      make([]Message, historySize),
		links:         make(map[string]*Link, historySize),
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
		for _, d := range r.links {
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

func (r *room) subscribe(name string, c *Client) (*Link, error) {
	d := &Link{c: c, r: r}
	r.clientLock.Lock()
	if _, ok := r.links[name]; ok {
		r.clientLock.Unlock()
		return nil, fmt.Errorf("%s is not unique", name)
	}
	r.links[name] = d
	r.clientLock.Unlock()
	return d, nil
}
