package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"mtt/pb"

	"google.golang.org/grpc"
)

func main() {
	modePtr := flag.String("m", "server", "mode to run chat in 'server' or 'client'")
	listenPtr := flag.String("l", "localhost:4444", "host and port to run or connect to")
	roomsPtr := flag.String("r", "", "rooms and according usernames in format room1:username1,room2:username2,... ")
	flag.Parse()

	fmt.Println(*modePtr, *listenPtr, *roomsPtr)

	lis, err := net.Listen("tcp", *listenPtr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServer(grpcServer, newServer())

	go grpcServer.Serve(lis)
	time.Sleep(time.Second)
	client(*listenPtr)
}

func client(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatClient(conn)

	fmt.Println("SUBSCRIBING")
	subscription, err := c.Subscribe(context.Background())
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}
	err = subscription.Send(&pb.OutgoingMessage{
		Room:      "room1",
		Subscribe: true,
		Username:  "Kenny",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("TEXTING")
	err = subscription.Send(&pb.OutgoingMessage{
		Room: "room1",
		Text: "testme",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("RECV")
	mp, err := subscription.Recv()
	if err != nil {
		panic(err)
	}
	fmt.Println(mp)

	mp, err = subscription.Recv()
	if err != nil {
		panic(err)
	}

	fmt.Println(mp)
}
