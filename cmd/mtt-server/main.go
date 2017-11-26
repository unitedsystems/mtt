package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/unitedsystems/mtt/pb"
	"google.golang.org/grpc"
)

func spawnServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// keepaliveOption := grpc.KeepaliveParams(keepalive.ServerParameters{Timeout: time.Second})
	grpcServer := grpc.NewServer()
	pb.RegisterChatServer(grpcServer, newGRPCServer())
	log.Printf("server is listening at %s", addr)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	addrPtr := flag.String("l", "localhost:4444", "host and port to listen to")
	flag.Parse()
	if *addrPtr == "" {
		panic("bad args")
	}
	spawnServer(*addrPtr)
}
