package main

import (
	// "context"
	"log"
	"net"

	pb "github.com/maya-fisher/birthday-service/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50053"
	port = ":50053"
)

type server struct {
	pb.UnimplementedBirthdaysServer
}

type Person struct {
	Name string
	Birthday int64
}

func (s *server) CreateBirthdayPersonBy(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetIdResponse, error) {

	log.Printf("Received: %v", in.GetPerson()) 
	return &pb.GetIdResponse{Id: in.GetId(), }
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBirthdaysServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }

	// defer conn.Close()

	// c := pb.NewBirthdaysClient(conn)
	// fmt.Println("client", c)
}