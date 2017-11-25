package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"mtt/pb"

	"google.golang.org/grpc"
)

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
	addrPtr := flag.String("l", "localhost:4444", "host and port to connect to")
	roomsPtr := flag.String("r", "", "rooms and according usernames in format room1:username1,room2:username2,... ")
	flag.Parse()

	spawnClient(*addrPtr, *roomsPtr)
}
