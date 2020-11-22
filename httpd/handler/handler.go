package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jonny-exe/go-server/httpd/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

// Article ....
type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
}

// Articles ...
type Articles []Article

// AddMessage adds users to mongodb database
func AddMessage(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.MessageRequest
	fmt.Println("AddMessage")
	json.NewDecoder(r.Body).Decode(&req)
	model := dbmodels.MessageModel{Sender: req.Sender, Receiver: req.Receiver, Content: req.Content, Date: time.Now()}
	insertResult, err := collectionMessages.InsertOne(context.TODO(), model)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Endpoint hit: all articles endpoint")

	json.NewEncoder(w).Encode(model)
}

// AddUser adds a user to the db
func AddUser(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.UserModel
	fmt.Println("AddMessage")
	json.NewDecoder(r.Body).Decode(&req)
	insertResult, err := collectionMessages.InsertOne(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Endpoint hit: all user endpoint")

	json.NewEncoder(w).Encode(req)
}

// AddFriend adds a friend to a user in the db
func AddFriend(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.FriendRequest
	fmt.Println("AddFriend")
	json.NewDecoder(r.Body).Decode(&req)
	err := collectionMessages.FindOne(context.TODO(), bson.D{}).Decode(&req)
	fmt.Println(err)
	// insertResult, err := collectionMessages.InsertOne(context.TODO(), req)
	fmt.Println("Found : ")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Endpoint hit: all friend endpoint")

	json.NewEncoder(w).Encode(req)
}

// HomePage ...
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint HIt")
}

// Mongodb types
// *mongo.Database
// *mongo.Collection
// *mongo.Client
var client *mongo.Client
var database *mongo.Database
var collectionUsers *mongo.Collection
var collectionMessages *mongo.Collection

// Connect connects to the mongodb db
func Connect() {
	enverr := godotenv.Load()
	fmt.Println(enverr)
	fmt.Println("Connecting to MongoDB")
	connectionKey := os.Getenv("DB_CONNECTION")
	fmt.Println(connectionKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(connectionKey))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection succesfull")
	fmt.Println(cancel)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Dosent need to be closed, its better not to close i
	// defer client.Disconnect(ctx)

	database = client.Database("test")
	collectionMessages = database.Collection("postmessages")
}
