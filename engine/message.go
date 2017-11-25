package engine

// Message represents chat message
type Message struct {
	Timestamp int64
	Room      string
	Name      string
	Text      string
}
