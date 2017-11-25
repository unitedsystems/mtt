package engine

import (
	"time"
)

// Client is representation of client in chat engine
type Client struct {
	server             *Server
	masterNotification chan struct{}
	descriptors        []*Descriptor
}

func newClient(s *Server) *Client {
	return &Client{
		server:             s,
		masterNotification: make(chan struct{}, 1),
		descriptors:        make([]*Descriptor, 0),
	}
}

// Poll queries all underlying descriptors
// and gets returns unsorted pack of new messages
func (c *Client) Poll() []Message {
	t := time.Now().UnixNano()
	result := make([]Message, 0)
	for _, cd := range c.descriptors {
		select {
		case <-cd.updateChan:
			messages := cd.pull(t)
			result = append(result, messages...)
		}
	}
	return result
}

// Disconnect cleans up connections
func (c *Client) Disconnect() {
}

// Subscribe binds connection with a room
func (c *Client) Subscribe(room, username string) error {
	r := c.server.getRoom(room)
	_, err := r.subscribe(username, c)
	return err
}

// Publish sends message to specified room
func (c *Client) Publish(room, text string) error {
	r := c.server.getRoom(room)
	r.publish(Message{
		Text: text,
	})
	return nil
}
