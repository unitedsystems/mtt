package engine

import (
	"fmt"
	"testing"
	"time"
)

func Benchmark__SingleChatSingleClient(b *testing.B) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	r := s.getRoom(roomName)

	c := newClient(s)
	c.Subscribe(roomName, userName)
	go func() {
		for {
			c.Poll()
		}
	}()

	time.Sleep(time.Millisecond)
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		r.publish("test", nil)
	}
}

func Benchmark__SingleChatManyClientsPoll(b *testing.B) {
	s := NewServer()
	roomName := "foo"
	userName := "John"
	numberOfClients := 10
	r := s.getRoom(roomName)

	for i := 0; i < numberOfClients; i++ {
		c := newClient(s)
		c.Subscribe(roomName, fmt.Sprintf("%s-%d", userName, i))
		go func() {
			for {
				c.Poll()
			}
		}()
	}

	time.Sleep(time.Millisecond)
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		r.publish("test", nil)
	}
}

func Benchmark__SingleClientManyRooms(b *testing.B) {
	s := NewServer()
	userName := "John"
	roomName := "room"
	numberOfRooms := 10
	c := newClient(s)
	go func() {
		for {
			c.Poll()
		}
	}()

	rooms := make([]*room, numberOfRooms)
	for i := 0; i < numberOfRooms; i++ {
		name := fmt.Sprintf("%s-%d", roomName, i)
		rooms[i] = s.getRoom(name)
		c.Subscribe(name, userName)
	}

	time.Sleep(time.Millisecond)
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		for j := 0; j < numberOfRooms; j++ {
			rooms[j].publish("test", nil)
		}
	}
}

func Benchmark__ManyToMany(b *testing.B) {
	s := NewServer()
	userName := "John"
	roomName := "room"
	numberOfRooms := 10
	numberOfClients := 10

	clients := make([]*Client, numberOfClients)
	for i := 0; i < numberOfClients; i++ {
		clients[i] = newClient(s)
		t := i
		go func() {
			for {
				clients[t].Poll()
			}
		}()
	}

	rooms := make([]*room, numberOfRooms)
	for i := 0; i < numberOfRooms; i++ {
		name := fmt.Sprintf("%s-%d", roomName, i)
		rooms[i] = s.getRoom(name)
		for j := 0; j < numberOfClients; j++ {
			clients[j].Subscribe(rooms[i].name, fmt.Sprintf("%s-%d-%d", userName, i, j))
		}
	}

	time.Sleep(time.Millisecond)
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		for j := 0; j < numberOfRooms; j++ {
			rooms[j].publish("test", nil)
		}
	}
}
