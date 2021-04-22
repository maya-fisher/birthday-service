package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"strconv"
	pb "github.com/maya-fisher/birthday-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port = ":50054"
)

type server struct {
	pb.UnimplementedBirthdaysServer
}

var Birthday_collection *mongo.Collection


func (s *server) CreateBirthdayPersonBy(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetIdResponse, error) {

	time := time.Unix(in.GetPerson().Birthday, 0)

	Birthday_collection.InsertOne(ctx, bson.D{
		{Key:"Name",Value: in.GetPerson().Name,},
		{Key:"Birthday",Value: time},
		{Key: "UserID", Value: in.GetPerson().UserId},
	})

	log.Printf("Received: %v", in.GetPerson()) 

	return &pb.GetIdResponse{Id: in.GetPerson().UserId}, nil
}  


func (s *server) UpdateBirthdayByIdAndName(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetIdResponse, error){

	time := time.Unix(in.GetPerson().Birthday, 0)
	update := bson.M{"$set": bson.M{"Birthday": time}}
	filter := bson.M{"UserID": bson.M{"$eq": in.GetPerson().UserId}}

	result, err := Birthday_collection.UpdateOne(
        context.Background(),
		filter,
        update,
    )

	if err != nil {
        fmt.Println("UpdateOne() result ERROR:", err)
    } else {
        fmt.Println("UpdateOne() result:", result)
	}

	filterFind, err := Birthday_collection.Find(ctx, bson.M{"UserID": in.GetPerson().UserId})

	if err != nil {
		log.Fatal(err)
	}	
	
	var filtered_users []bson.D


	if err = filterFind.All(ctx, &filtered_users); err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("%v", filtered_users[0][1].Value)
	userId := fmt.Sprintf("%v", filtered_users[0][3].Value)
	unconverted_birthday := fmt.Sprintf("%v", filtered_users[0][2].Value)
	birthday, err := strconv.ParseInt(unconverted_birthday, 10, 64)

	person := &pb.Person{
		Name: name,
		Birthday: birthday,
		UserId: userId,
	}


	fmt.Println(person)

	return &pb.GetIdResponse{Id: in.GetPerson().UserId}, nil
} 


func (s *server) GetBirthdayPersonByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetBirthdayResponse, error) {

	filter, err := Birthday_collection.Find(ctx, bson.M{"UserID": in.GetId()})

	if err != nil {
		log.Fatal(err)
	}

	var filtered_users []bson.D

	if err = filter.All(ctx, &filtered_users); err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("%v", filtered_users[0][1].Value)
	userId := fmt.Sprintf("%v", filtered_users[0][3].Value)
	unconverted_birthday := fmt.Sprintf("%v", filtered_users[0][2].Value)
	birthday, err := strconv.ParseInt(unconverted_birthday, 10, 64)

	person := &pb.Person{
		Name: name,
		Birthday: birthday,
		UserId: userId,
	}

	log.Printf("Received: %v", in.GetId()) 

	return &pb.GetBirthdayResponse{Person: person}, nil
} 


func (s *server) DeleteBirthdayByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetIdResponse, error) {

	result, err := Birthday_collection.DeleteOne(ctx, bson.M{"UserID": in.GetId()})

	if err != nil {
    	log.Fatal(err)
	}

	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	log.Printf("Received: %v", in.GetId()) 

	return &pb.GetIdResponse{Id: in.GetId()}, nil
} 


func main() {

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