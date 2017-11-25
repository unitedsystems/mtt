package engine

import "sync"

// Server provides core messaging functionality
type Server struct {
	sync.RWMutex
	rooms map[string]*room
}

func (s *Server) getRoom(name string) *room {
	s.Lock()
	if _, ok := s.rooms[name]; !ok {
		s.rooms[name] = newRoom(name)
	}
	s.Unlock()
	return s.rooms[name]
}

// NewServer creates instance of chat server engine
func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*room),
	}
}

// Connect creates client entity in the engine
func (s *Server) Connect() *Client {
	return newClient(s)
}
