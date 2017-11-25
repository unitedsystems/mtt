package engine

import (
	"fmt"
	"testing"
	"time"
)

func Test__RoomMessagePublishing(t *testing.T) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)

	numberOfMessages := historySize * 3 / 2
	for i := 0; i < numberOfMessages; i++ {
		m := fmt.Sprintf("message%d", i)
		r.publish(Message{
			Name: userName,
			Text: m,
		})
	}

	if r.lastID != numberOfMessages {
		t.Errorf("room lastID does not gets incremented %d (should be %d)", r.lastID, numberOfMessages)
	}

	actualMesssage := r.messages[0].Text
	expectedMessage := fmt.Sprintf("message%d", historySize)
	if actualMesssage != expectedMessage {
		t.Errorf("first room message does not get overwritten %s (should be %s)", expectedMessage, actualMesssage)
	}

	actualMesssage = r.messages[historySize-1].Text
	expectedMessage = fmt.Sprintf("message%d", historySize-1)
	if actualMesssage != expectedMessage {
		t.Errorf("last room message does not get overwritten %s (should be %s)", expectedMessage, actualMesssage)
	}

	actualMesssage = r.messages[numberOfMessages%historySize].Text
	expectedMessage = fmt.Sprintf("message%d", numberOfMessages%historySize)
	if actualMesssage != expectedMessage {
		t.Errorf("next after head room message should not be overwritten %s (should be %s)", expectedMessage, actualMesssage)
	}
}

func Test__RoomSubscription(t *testing.T) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)

	_, err := r.subscribe(userName, newClient(s))
	if err != nil {
		t.Errorf("got error on room subscription: %v", err)
	}
	fmt.Println()

	_, err = r.subscribe(userName, newClient(s))
	if err == nil {
		t.Errorf("didn't get error on second room subscription")
	}
}

func Test__RoomBroadcast(t *testing.T) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	messageText := "Hallo!"
	r := s.getRoom(roomName)
	c := newClient(s)
	r.subscribe(userName, c)

	r.publish(Message{
		Name: userName,
		Text: messageText,
	})

	select {
	case <-c.masterNotification:
	case <-time.NewTimer(time.Millisecond).C:
		t.Errorf("room does not update descriptors")
	}
}
