package engine

import (
	"testing"
)

func Test__ClientPoll(t *testing.T) {
	s := NewServer()
	r1 := s.getRoom("room1")
	r2 := s.getRoom("room2")

	c := newClient(s)
	r1.subscribe("john", c)
	r2.subscribe("jack", c)

	r1.publish(Message{
		Text: "test1",
	})
	r2.publish(Message{
		Text: "test2",
	})
	r2.publish(Message{
		Text: "test3",
	})

	<-c.masterNotification
	actualResult := c.Poll()
	expectedLen := 3
	if len(actualResult) != expectedLen {
		t.Errorf("client failed to get %d messages (got %d)", expectedLen, len(actualResult))
	}
}
