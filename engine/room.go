package engine

import (
	"fmt"
	"sync"
	"time"
)

type room struct {
	sync.RWMutex
	messages []roomMessage
	lastID   uint64 // TODO: process overflow
	name     string

	linksLock sync.RWMutex
	links     map[string]*Link // NOTE: key is user's name for this particular room

	broadcast bool
}

func newRoom(name string) *room {
	r := &room{
		name:      name,
		messages:  make([]roomMessage, historySize),
		links:     make(map[string]*Link, historySize),
		linksLock: sync.RWMutex{},
	}
	for i := range r.messages {
		r.messages[i].allocateBuffer()
	}
	go r.askClientsToPoll()
	return r
}

func (r *room) askClientsToPoll() {
	for {
		if r.broadcast {
			r.Lock()
			r.broadcast = false
			r.Unlock()
		} else {
			time.Sleep(time.Millisecond)
			continue
		}

		r.linksLock.RLock()
		for _, l := range r.links {
			select {
			case l.c.masterNotification <- struct{}{}:
			default:
			}
		}
		r.linksLock.RUnlock()
	}
}

func (r *room) publish(text string, link *Link) {
	r.Lock()
	defer r.Unlock()
	idx := r.lastID % historySize
	r.messages[idx].write(text)
	r.messages[idx].Link = link
	r.lastID++
	r.broadcast = true
}

func (r *room) subscribe(name string, c *Client) (*Link, error) {
	l := &Link{c: c, r: r}
	r.linksLock.Lock()
	if _, ok := r.links[name]; ok {
		r.linksLock.Unlock()
		return nil, fmt.Errorf("%s is not unique", name)
	}
	r.links[name] = l
	r.linksLock.Unlock()
	return l, nil
}
