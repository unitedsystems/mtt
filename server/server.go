package main

import (
	"fmt"
	"log"
	"mtt/engine"
	"mtt/pb"
)

type grpcServer struct {
	engine *engine.Server
}

func (s *grpcServer) listen(ss pb.Chat_SubscribeServer, client *engine.Client) {
	for {
		m, err := ss.Recv()
		if err != nil {
			log.Printf("error in RECV on server: %#v", err)
			return
		}
		if m.Subscribe {
			log.Printf("%s subscribed to %s", m.Username, m.Room)
			err = client.Subscribe(m.Room, m.Username)
			if err != nil {
				log.Printf("error during Subscribe: %#v", err)
			}
		} else {
			log.Printf("%s sent '%s' to %s", m.Username, m.Text, m.Room)
			err = client.Publish(m.Room, m.Text)
			if err != nil {
				log.Printf("error during Publish: %#v", err)
			}
		}
	}
}

func (s *grpcServer) serve(ss pb.Chat_SubscribeServer, client *engine.Client) {
	for {
		messages := client.Poll()
		fmt.Println("Polled", len(messages), "messages")
		messagePack := new(pb.MessagePack)
		messagePack.Messages = make([]*pb.IncomingMessage, len(messages))
		for i, message := range messages {
			messagePack.Messages[i] = &pb.IncomingMessage{
				Room:      message.Room,
				Username:  message.Name,
				Timestamp: message.Timestamp,
				Text:      message.Text,
			}
		}
		err := ss.Send(messagePack)
		if err != nil {
			log.Printf("error in SEND on server " + err.Error())
			return
		}
	}
}

func (s *grpcServer) Subscribe(ss pb.Chat_SubscribeServer) error {
	client := s.engine.Connect()
	go s.listen(ss, client)
	s.serve(ss, client)
	return nil
}

func newGRPCServer() *grpcServer {
	return &grpcServer{
		engine: engine.NewServer(),
	}
}
