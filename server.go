package main

import (
	"log"

	"google.golang.org/grpc"
	pb ""
)

const (
	address = "localhost:50052"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

}