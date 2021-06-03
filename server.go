package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/maya-fisher/birthday-service/proto"
	"github.com/maya-fisher/birthday-service/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedBirthdaysServer
}

var Birthday_collection *mongo.Collection

func (s *server) CreateBirthdayPersonBy(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetBirthdayResponse, error) {

	log.Printf("Received: %v", in.GetPerson())

	time := time.Unix(in.GetPerson().Birthday, 0)

	_, err := Birthday_collection.InsertOne(ctx, bson.D{
		{Key: "Name", Value: in.GetPerson().Name},
		{Key: "Birthday", Value: time},
		{Key: "UserID", Value: in.GetPerson().UserId},
	})

	if err != nil {
		return &pb.GetBirthdayResponse{Person: nil}, err

	} else {
		person, err := getBrthdayByID(in.GetPerson().UserId, ctx)
		return &pb.GetBirthdayResponse{Person: person}, err
	}
}

func (s *server) UpdateBirthdayByIdAndName(ctx context.Context, in *pb.GetBirthdayRequest) (*pb.GetBirthdayResponse, error) {

	fmt.Println("PERSON:", in.GetPerson())
	if in.GetPerson().Name != "" {
		update := bson.M{"$set": bson.M{"Name": in.GetPerson().Name}}
		filter := bson.M{"UserID": bson.M{"$eq": in.GetPerson().UserId}}
		result, _ := Birthday_collection.UpdateOne(
			context.Background(),
			filter,
			update,
		)

		fmt.Println("UpdateOne() result:", result)

		person, err := getBrthdayByID(in.GetPerson().UserId, ctx)

		if err != nil {

			return &pb.GetBirthdayResponse{Person: nil}, err
		} else {

			return &pb.GetBirthdayResponse{Person: person}, nil

		}
	}

	time := time.Unix(in.GetPerson().Birthday, 0)
	update := bson.M{"$set": bson.M{"Birthday": time}}
	filter := bson.M{"UserID": bson.M{"$eq": in.GetPerson().UserId}}

	result, _ := Birthday_collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	fmt.Println("UpdateOne() result:", result)

	person, err := getBrthdayByID(in.GetPerson().UserId, ctx)

	if err != nil {

		return &pb.GetBirthdayResponse{Person: nil}, err
	} else {
		return &pb.GetBirthdayResponse{Person: person}, nil

	}
}

func (s *server) DeleteBirthdayByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetBirthdayResponse, error) {

	person, err := getBrthdayByID(in.GetUserId(), ctx)

	log.Printf("Received: %v", in.GetUserId())

	result, _ := Birthday_collection.DeleteOne(ctx, bson.M{"UserID": in.GetUserId()})

	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	if err != nil {
		return &pb.GetBirthdayResponse{Person: nil}, err
	} else {

		return &pb.GetBirthdayResponse{Person: person}, nil

	}
}

func (s *server) GetBirthdayPersonByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetBirthdayResponse, error) {

	person, err := getBrthdayByID(in.GetUserId(), ctx)

	log.Printf("Received: %v", in.GetUserId())

	if err != nil {

		return &pb.GetBirthdayResponse{Person: nil}, err
	} else {

		return &pb.GetBirthdayResponse{Person: person}, nil

	}
}

func getBrthdayByID(id string, ctx context.Context) (*pb.Person, error) {

	var res bson.M

	err := Birthday_collection.FindOne(ctx, bson.M{"UserID": id}).Decode(&res)
	if err == mongo.ErrNoDocuments {

		person := &pb.Person{}

		return person, err

	}

	name := fmt.Sprintf("%v", res["Name"])
	userId := fmt.Sprintf("%v", res["UserID"])
	unconverted_birthday := fmt.Sprintf("%v", res["Birthday"])
	birthday, err := strconv.ParseInt(unconverted_birthday, 10, 64)
	if err != nil {
	}

	person := &pb.Person{
		Name:     name,
		Birthday: birthday,
		UserId:   userId,
	}

	return person, nil
}

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGO_URL))
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

	lis, err := net.Listen("tcp", config.PORT)
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
