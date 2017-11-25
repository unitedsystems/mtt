package server

import (
	"fmt"
	"testing"
	"time"
)

func Test__RoomMessagePublishing(t *testing.T) {
	s := newServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)

	numberOfMessages := historySize * 3 / 2
	for i := 0; i < numberOfMessages; i++ {
		m := fmt.Sprintf("message%d", i)
		r.Publish(message{
			name: userName,
			text: m,
		})
	}

	if r.lastID != numberOfMessages {
		t.Errorf("room lastID does not gets incremented %d (should be %d)", r.lastID, numberOfMessages)
	}

	actualMesssage := r.messages[0].text
	expectedMessage := fmt.Sprintf("message%d", historySize)
	if actualMesssage != expectedMessage {
		t.Errorf("first room message does not get overwritten %s (should be %s)", expectedMessage, actualMesssage)
	}

	actualMesssage = r.messages[historySize-1].text
	expectedMessage = fmt.Sprintf("message%d", historySize-1)
	if actualMesssage != expectedMessage {
		t.Errorf("last room message does not get overwritten %s (should be %s)", expectedMessage, actualMesssage)
	}

	actualMesssage = r.messages[numberOfMessages%historySize].text
	expectedMessage = fmt.Sprintf("message%d", numberOfMessages%historySize)
	if actualMesssage != expectedMessage {
		t.Errorf("next after head room message should not be overwritten %s (should be %s)", expectedMessage, actualMesssage)
	}
}

func Test__RoomSubscription(t *testing.T) {
	s := newServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)

	_, err := r.Subscribe(userName, newClient())
	if err != nil {
		t.Errorf("got error on room subscription: %v", err)
	}

	_, err = r.Subscribe(userName, newClient())
	if err == nil {
		t.Errorf("didn't get error on second room subscription")
	}
}

func Test__RoomBroadcast(t *testing.T) {
	s := newServer()
	roomName := "foo"
	userName := "John"
	messageText := "Hallo!"
	r := s.getRoom(roomName)
	d, _ := r.Subscribe(userName, newClient())

	r.Publish(message{
		name: userName,
		text: messageText,
	})

	select {
	case <-d.updateChan:
	case <-time.NewTimer(time.Millisecond).C:
		t.Errorf("room does not update descriptors")
	}
}
