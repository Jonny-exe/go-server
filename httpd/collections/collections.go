package collections

import (
	"context"
	"fmt"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Connecting to MongoDB")
	connectionKey := os.Getenv("DB_CONNECTION")
	log.Println(connectionKey)
	Client, err := mongo.NewClient(options.Client().ApplyURI(connectionKey))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection succesfull")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println(cancel)
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer Client.Disconnect(ctx)

	Users := Client.Database("my_database").Collection("posts")
	fmt.Println(Users, Client)

}
