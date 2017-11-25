package main

import (
	"log"
	"mtt/pb"
	"time"
)

type grpcServer struct {
}

func (s *grpcServer) listen(ss pb.Chat_SubscribeServer) {
	for {
		m, err := ss.Recv()
		if err != nil {
			log.Printf("error in RECV on server: %#v", err)
		}
		if m.Subscribe {
			log.Printf("%s subscribed to %s", m.Username, m.Room)
		} else {
			log.Printf("%s sent '%s' to %s", m.Username, m.Text, m.Room)
		}
	}
}

func (s *grpcServer) serve(ss pb.Chat_SubscribeServer) {
	for {
		time.Sleep(time.Second * 10)
		mp := new(pb.MessagePack)
		mp.Messages = make([]*pb.IncomingMessage, 1)
		mp.Messages[0] = &pb.IncomingMessage{Text: "some message"}
		err := ss.Send(mp)
		if err != nil {
			log.Fatalf("error in SEND on server " + err.Error())
		}
	}
}

func (s *grpcServer) Subscribe(ss pb.Chat_SubscribeServer) error {
	go s.listen(ss)
	s.serve(ss)
	return nil
}

func newGRPCServer() *grpcServer {
	return new(grpcServer)
}
