package engine

import "testing"

func Test__MessageCreation(t *testing.T) {
	messageText := "test"
	m := new(roomMessage)
	m.allocateBuffer()
	m.write(messageText)
	if len(m.Buffer) != len(messageText) {
		t.Errorf("message fails to copy text to buffer on creation (length is %d)", len(m.Buffer))
	}
	if cap(m.Buffer) != maxMessageSize {
		t.Errorf("message fails to allocate buffer on creation (capacity is %d)", cap(m.Buffer))
	}
}

func Test__MessageExport(t *testing.T) {
	messageText := "test"
	roomName := "foo"
	cleintName := "john"

	m := new(roomMessage)
	m.allocateBuffer()
	m.write(messageText)

	s := NewServer()
	c := newClient(s)
	c.Subscribe(roomName, cleintName)

	r := s.getRoom(roomName)
	m.Link = r.links[cleintName]

	exported := m.export()
	if exported.Username != cleintName {
		t.Errorf("export fails to set username")
	}
	if exported.Room != roomName {
		t.Errorf("export fails to set room name")
	}
}
