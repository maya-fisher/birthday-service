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


func (s *server) CreateBirthdayPersonBy(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetBirthdayResponse, error) {

	time := time.Unix(in.GetPerson().Birthday, 0)

	result, err := Birthday_collection.InsertOne(ctx, bson.D{
		{Key:"Name",Value: in.GetPerson().Name,},
		{Key:"Birthday",Value: time},
		{Key: "UserID", Value: in.GetPerson().UserId},
	})

	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println("result:", result)
	}


	log.Printf("Received: %v", in.GetPerson()) 

	person := getBrthdayByID(in.GetPerson().UserId, ctx)

	return &pb.GetBirthdayResponse{Person: person}, err
}  


func (s *server) UpdateBirthdayByIdAndName(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetBirthdayResponse, error){

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

	person := getBrthdayByID(in.GetPerson().UserId, ctx)

	return &pb.GetBirthdayResponse{Person: person}, nil
} 


func (s *server) GetBirthdayPersonByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetBirthdayResponse, error) {

	person := getBrthdayByID(in.GetUserId(), ctx)

	log.Printf("Received: %v", in.GetUserId()) 

	return &pb.GetBirthdayResponse{Person: person}, nil
} 


func (s *server) DeleteBirthdayByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetBirthdayResponse, error) {

	person := getBrthdayByID(in.GetUserId(), ctx)

	result, err := Birthday_collection.DeleteOne(ctx, bson.M{"UserID": in.GetUserId()})

	if err != nil {
    	log.Fatal(err)
	}

	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	log.Printf("Received: %v", in.GetUserId()) 

	return &pb.GetBirthdayResponse{Person: person}, nil
} 


func getBrthdayByID(id string, ctx context.Context) (*pb.Person) {

	filter, err := Birthday_collection.Find(ctx, bson.M{"UserID": id})

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

	return person

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