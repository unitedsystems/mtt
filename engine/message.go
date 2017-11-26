package engine

import "time"

type roomMessage struct {
	Timestamp int64
	Link      *Link
	Buffer    []byte
}

func newMessage(text string) roomMessage {
	m := new(roomMessage)
	m.allocateBuffer()
	m.write(text)
	return *m
}

func (m *roomMessage) allocateBuffer() {
	m.Buffer = make([]byte, 0, maxMessageSize)
}

func (m *roomMessage) write(text string) {
	textBytes := []byte(text)
	m.Buffer = m.Buffer[0:len(textBytes)]
	copy(m.Buffer, textBytes)
	m.Timestamp = time.Now().UnixNano()
}

func (m *roomMessage) export() Message {
	result := Message{
		Timestamp: m.Timestamp,
		Text:      m.Buffer,
	}
	if m.Link != nil {
		result.Username = m.Link.name
		result.Room = m.Link.r.name
	}
	return result
}

// Message is interface for outgoing chat message
type Message struct {
	Timestamp int64
	Username  string
	Text      []byte
	Room      string
}
