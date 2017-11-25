package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"mtt/pb"

	"google.golang.org/grpc"
)

func spawnServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServer(grpcServer, newGRPCServer())
	log.Printf("server is listening at %s", addr)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func spawnClient(addr, rooms string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewChatClient(conn)

	// stub for requests
	numberOfRooms := 3
	username := fmt.Sprintf("Kenny-%d", rand.Int())

	subscription, err := c.Subscribe(context.Background())
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}
	for i := 1; i <= numberOfRooms; i++ {
		err = subscription.Send(&pb.OutgoingMessage{
			Room:      fmt.Sprintf("room%d", i),
			Subscribe: true,
			Username:  username,
		})
		if err != nil {
			panic(err)
		}
	}

	for {
		time.Sleep(time.Second)
		err = subscription.Send(&pb.OutgoingMessage{
			Room:     fmt.Sprintf("room%d", rand.Int()%3),
			Text:     time.Now().String(),
			Username: username,
		})
		if err != nil {
			panic(err)
		}
		mp, err := subscription.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println(mp)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	modePtr := flag.String("m", "server", "mode to run chat in 'server' or 'client'")
	addrPtr := flag.String("l", "localhost:4444", "host and port to run or connect to")
	roomsPtr := flag.String("r", "", "rooms and according usernames in format room1:username1,room2:username2,... ")
	flag.Parse()

	if *modePtr == "server" {
		spawnServer(*addrPtr)
		fmt.Println("here")
	} else {
		spawnClient(*addrPtr, *roomsPtr)
	}
}
