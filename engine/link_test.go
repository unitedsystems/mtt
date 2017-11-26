package engine

import (
	"fmt"
	"testing"
	"time"
)

func Test__LinkPull(t *testing.T) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)
	d, _ := r.subscribe(userName, newClient(s))

	numberOfMessages := historySize * 3 / 2
	for i := 0; i < numberOfMessages; i++ {
		r.publish(fmt.Sprintf("message%d", i), nil)
	}

	data := d.pull(time.Now().UnixNano())

	if len(data) != historySize {
		t.Errorf("failed to pull data from link, expected %d (actual %d)", historySize, len(data))
	}
}
