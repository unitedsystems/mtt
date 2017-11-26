package engine

import "testing"

func Test__MessageCreation(t *testing.T) {
	messageText := "test"
	m := newMessage(messageText)
	if len(m.Buffer) != len(messageText) {
		t.Errorf("message fails to copy text to buffer on creation (length is %d)", len(m.Buffer))
	}
	if cap(m.Buffer) != maxMessageSize {
		t.Errorf("message fails to allocate buffer on creation (capacity is %d)", cap(m.Buffer))
	}
}
