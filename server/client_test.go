package server

import (
	"testing"
)

func Test__ClientPoll(t *testing.T) {
	s := newServer()
	r1 := s.getRoom("room1")
	r2 := s.getRoom("room2")

	c := newClient()
	r1.Subscribe("john", c)
	r2.Subscribe("jack", c)

	r1.Publish(message{
		text: "test1",
	})
	r2.Publish(message{
		text: "test2",
	})
	r2.Publish(message{
		text: "test3",
	})

	<-c.masterNotification
	actualResult := c.poll()
	expectedLen := 3
	if len(actualResult) != expectedLen {
		t.Errorf("client failed to get %d messages (got %d)", expectedLen, len(actualResult))
	}
}
