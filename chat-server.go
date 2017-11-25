package main

import (
	"fmt"
	"mtt/pb"
)

type chatServer struct {
}

func (s *chatServer) Subscribe(ss pb.Chat_SubscribeServer) error {
	for {
		out, err := ss.Recv()
		if err != nil {
			panic("error in RECV on server " + err.Error())
		}

		fmt.Println("subscribe SRV", out)

		mp := new(pb.MessagePack)
		mp.Messages = make([]*pb.IncomingMessage, 1)
		mp.Messages[0] = &pb.IncomingMessage{Text: "got your message"}
		err = ss.Send(mp)
		if err != nil {
			panic("error in SEND on server " + err.Error())
		}
	}
}

func newServer() *chatServer {
	return new(chatServer)
}
