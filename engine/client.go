package engine

import (
	"fmt"
	"sync"
	"time"
)

// Client is representation of client in chat engine
type Client struct {
	sync.RWMutex
	server             *Server
	masterNotification chan struct{}
	links              map[string]*Link // NOTE: key is room name, thus we are sure that username is unique for this room
}

func newClient(s *Server) *Client {
	return &Client{
		server:             s,
		masterNotification: make(chan struct{}, 1),
		links:              make(map[string]*Link),
	}
}

// Poll queries all underlying links
// and gets returns unsorted pack of new messages
func (c *Client) Poll() []Message {
	<-c.masterNotification
	c.Lock()
	defer c.Unlock()
	t := time.Now().UnixNano()
	result := make([]Message, 0, historySize)
	for _, l := range c.links {
		messages := l.pull(t)
		result = append(result, messages...)
	}
	return result
}

// Disconnect cleans up connections
func (c *Client) Disconnect() {
	c.Lock()
	defer c.Unlock()
	oldLinks := c.links
	c.links = nil
	for _, d := range oldLinks {
		delete(d.r.links, d.name)
	}
}

// Subscribe binds connection with a room
func (c *Client) Subscribe(room, username string) error {
	r := c.server.getRoom(room)
	d, err := r.subscribe(username, c)
	if err != nil {
		panic(err)
	}
	c.Lock()
	c.links[r.name] = d
	c.Unlock()
	// let's read history
	select {
	case c.masterNotification <- struct{}{}:
	default:
	}
	return err
}

// Publish sends message to specified room
func (c *Client) Publish(room, text string) error {
	if _, ok := c.links[room]; !ok {
		return fmt.Errorf("can't send message to %s (not subscribed)", room)
	}
	r := c.server.getRoom(room)
	r.publish(text, c.links[room])
	return nil
}
