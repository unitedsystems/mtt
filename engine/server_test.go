package engine

import "testing"

func Test__RoomSpawning(t *testing.T) {
	s := newServer()
	roomName := "foo"

	firstCall := s.getRoom(roomName)
	if firstCall == nil {
		t.Error("server fails to create a room")
	}
	secondCall := s.getRoom(roomName)
	if firstCall != secondCall {
		t.Error("server creates new room each time")
	}
}
