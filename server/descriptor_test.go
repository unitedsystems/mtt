package server

import (
	"fmt"
	"testing"
	"time"
)

func Test__DescriptorPull(t *testing.T) {
	s := newServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)
	d, _ := r.Subscribe(userName, newClient())

	numberOfMessages := historySize * 3 / 2
	for i := 0; i < numberOfMessages; i++ {
		m := fmt.Sprintf("message%d", i)
		r.Publish(message{
			name: userName,
			text: m,
		})
	}

	data := d.pull(time.Now().UnixNano())

	if len(data) != historySize {
		t.Errorf("descriptor fails to pull data expected %d (actual %d)", historySize, len(data))
	}
}
