package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/unitedsystems/mtt/pb"

	"google.golang.org/grpc"
)

func spawnClient(addr, rooms string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewChatClient(conn)
	subscription, err := c.Subscribe(context.Background())
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	subscribe(subscription, rooms)
	go publisher(subscription)
	receiver(subscription)
}

func receiver(c pb.Chat_SubscribeClient) {
	for {
		messagePack, err := c.Recv()
		if err != nil {
			panic(err)
		}
		// TBD: sort messages by timestamp
		sort.Sort(messagePack)
		for _, message := range messagePack.Messages {
			fmt.Printf(
				"%s @%s %s> %s\n",
				message.Room,
				message.Username,
				time.Unix(0, message.Timestamp).Format("15:04:05"),
				message.Text,
			)
		}
	}

}

func subscribe(c pb.Chat_SubscribeClient, rooms string) {
	tokens := strings.Split(rooms, ",")
	for _, token := range tokens {
		parts := strings.Split(token, ":")
		err := c.Send(&pb.OutgoingMessage{
			Room:      parts[0],
			Username:  parts[1],
			Subscribe: true,
		})
		if err != nil {
			panic(err)
		}
	}
}

func publisher(c pb.Chat_SubscribeClient) {
	var room, message, line string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			line = scanner.Text()
		}
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			fmt.Println("wrong message format, should be: room your message here")

			continue
		}
		room = parts[0]
		message = strings.Join(parts[1:len(parts)], " ")
		err := c.Send(&pb.OutgoingMessage{
			Room: room,
			Text: message,
		})
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	addrPtr := flag.String("l", "localhost:4444", "host and port to connect to")
	roomsPtr := flag.String("r", "", "rooms and according usernames in format room1:username1,room2:username2,... ")
	flag.Parse()

	if *addrPtr == "" || *roomsPtr == "" {
		panic("bad args")
	}

	spawnClient(*addrPtr, *roomsPtr)
}
