package main

import (
	"log"

	"google.golang.org/grpc"
	""
)

const (
	address = "localhost:50053"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

}