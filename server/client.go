package server

import (
	"time"
)

type client struct {
	masterNotification chan struct{}
	descriptors        []*clientDescriptor
	connection         struct{}
}

func (c *client) backgroundPoll() {
	for {
		<-c.masterNotification
		c.flush(c.poll())
	}
}

func (c *client) poll() []message {
	t := time.Now().UnixNano()
	result := make([]message, 0)
	for _, cd := range c.descriptors {
		select {
		case <-cd.updateChan:
			messages := cd.pull(t)
			result = append(result, messages...)
		}
	}
	return result
}

func (c *client) flush(messages []message) {
	// write to socket here
}

func newClient() *client {
	return &client{
		masterNotification: make(chan struct{}, 1),
		descriptors:        make([]*clientDescriptor, 0),
	}
}
