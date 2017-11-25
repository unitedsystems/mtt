package main

import (
	"log"
	"mtt/engine"
	"mtt/pb"
	"time"
)

type grpcServer struct {
	engine *engine.Server
}

func newGRPCServer() *grpcServer {
	return &grpcServer{
		engine: engine.NewServer(),
	}
}

func (s *grpcServer) sendError(ss pb.Chat_SubscribeServer, errorToSend error) {
	messagePack := new(pb.MessagePack)
	messagePack.Messages = make([]*pb.IncomingMessage, 1)
	messagePack.Messages[0] = &pb.IncomingMessage{
		Room:      "system",
		Username:  "@chatbot",
		Timestamp: time.Now().Unix(),
		Text:      errorToSend.Error(),
	}
	err := ss.Send(messagePack)
	if err != nil {
		log.Printf("can't send error:" + errorToSend.Error())
		return
	}
}

func (s *grpcServer) processSubscription(client *engine.Client, m *pb.OutgoingMessage) error {
	log.Printf("%s subscribed to %s", m.Username, m.Room)
	err := client.Subscribe(m.Room, m.Username)
	if err != nil {
		log.Printf("error during Subscribe: %#v", err)
	}
	return err
}

func (s *grpcServer) processPublication(client *engine.Client, m *pb.OutgoingMessage) error {
	log.Printf("%s sent '%s' to %s", m.Username, m.Text, m.Room)
	err := client.Publish(m.Room, m.Text)
	if err != nil {
		log.Printf("error during Publish: %#v", err)
	}
	return err
}

func (s *grpcServer) listen(ss pb.Chat_SubscribeServer, client *engine.Client) {
	for {
		m, err := ss.Recv()
		if err != nil {
			log.Printf("error in RECV on server: %#v", err)
			client.Disconnect()
			return
		}
		if m.Subscribe {
			err = s.processSubscription(client, m)
		} else {
			err = s.processPublication(client, m)
		}
		if err != nil {
			s.sendError(ss, err)
		}
	}
}

func (s *grpcServer) serve(ss pb.Chat_SubscribeServer, client *engine.Client) {
	for {
		messages := client.Poll()
		if len(messages) == 0 {
			log.Println("WARNING, polled 0 messages")
			continue
		}
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
		log.Println("Polled and sending", len(messages), "messages")
		err := ss.Send(messagePack)
		if err != nil {
			log.Printf("error in SEND on server " + err.Error())
			client.Disconnect()
			return
		}
	}
}

func (s *grpcServer) Subscribe(ss pb.Chat_SubscribeServer) error {
	client := s.engine.Connect()
	go s.serve(ss, client)
	s.listen(ss, client)
	return nil
}
