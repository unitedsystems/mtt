package engine

// Link is a link between chat room and client
type Link struct {
	r      *room
	c      *Client
	lastID uint64
	name   string
	buffer []Message
}

// pull gets last messages from room
// it keeps track of last message id in the room
// and gets updates up to specified date
// new updates will be pulled on the next call
func (l *Link) pull(maxTimestamp int64) []Message {
	room := l.r
	if l.buffer == nil {
		l.buffer = make([]Message, 0, historySize)
	}
	room.RLock()
	defer room.RUnlock()
	span := room.lastID - l.lastID
	if span == 0 {
		return nil
	}
	if span > historySize {
		span = historySize
	}
	result := l.buffer[0:0]
	for i := uint64(0); i < span; i++ {
		idx := (room.lastID - span + i) % historySize
		if room.messages[idx].Timestamp > maxTimestamp {
			return result
		}
		result = append(result, room.messages[idx].export())
		l.lastID = room.lastID - span + i + 1
	}
	return result
}

func (l *Link) updateAvailable() {
	select {
	case l.c.masterNotification <- struct{}{}:
	default:
	}
}
