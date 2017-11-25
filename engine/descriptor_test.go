package engine

import (
	"fmt"
	"testing"
	"time"
)

func Test__DescriptorPull(t *testing.T) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)
	d, _ := r.subscribe(userName, newClient(s))

	numberOfMessages := historySize * 3 / 2
	for i := 0; i < numberOfMessages; i++ {
		m := fmt.Sprintf("message%d", i)
		r.publish(Message{
			Name: userName,
			Text: m,
		})
	}

	data := d.pull(time.Now().UnixNano())

	if len(data) != historySize {
		t.Errorf("descriptor fails to pull data expected %d (actual %d)", historySize, len(data))
	}
}
