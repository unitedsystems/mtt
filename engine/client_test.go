package engine

import (
	"testing"
)

func Test__ClientPoll(t *testing.T) {
	s := NewServer()
	r1 := s.getRoom("room1")
	r2 := s.getRoom("room2")

	c := newClient(s)
	c.Subscribe("room1", "john")
	c.Subscribe("room2", "jack")

	r1.publish("test1", nil)
	r2.publish("test2", nil)
	r2.publish("test3", nil)

	<-c.masterNotification
	actualResult := c.Poll()
	expectedLen := 3
	if len(actualResult) != expectedLen {
		t.Errorf("client failed to get %d messages (got %d)", expectedLen, len(actualResult))
	}
}
