package main

import (
	"context"
	"log"
	"net"
	"fmt"

	pb "github.com/maya-fisher/birthday-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

var Birthday_collection *mongo.Collection

func (s *server) CreateBirthdayPersonBy(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetIdResponse, error) {
	result, err := Birthday_collection.InsertOne(ctx, bson.D{
		{Key:"name",Value: in.GetPerson().Name,},
		{Key:"birthday",Value: in.GetPerson().Birthday},
	})
	
	if err != nil {
		log.Fatal(err)
	}

	InsertedID := result.InsertedID
	fmt.Println(InsertedID)
	log.Printf("Received: %v", in.GetPerson()) 
	return &pb.GetIdResponse{Id: "InsertedID"}, nil
}

func (s *server) UpdateBirthdayByIdAndName(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetIdResponse, error){
	log.Printf("Received: %v", in.GetPerson()) 
	return &pb.GetIdResponse{Id: "122"}, nil
}

func (s *server) GetBirthdayPersonByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetBirthdayResponse, error) {
	log.Printf("Received: %v", in.GetId()) 
	return &pb.GetBirthdayResponse{}, nil
}

func (s *server) DeleteBirthdayByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetIdResponse, error) {
	log.Printf("Received: %v", in.GetId()) 
	return &pb.GetIdResponse{}, nil
}

var birthday_collection *mongo.Collection

func main() {

	// connection to mongodb 

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/birthday_service"))
    if err != nil {
        log.Fatal(err)
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    defer client.Disconnect(ctx)

	db := client.Database("birthday_service")
	Birthday_collection = db.Collection("birthday")


	// 


	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBirthdaysServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}