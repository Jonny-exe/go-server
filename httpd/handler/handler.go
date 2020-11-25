package handler

// skipping primitive.E
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jonny-exe/go-server/httpd/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // FRED
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

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
	json.NewDecoder(r.Body).Decode(&req)
	insertResult, err := collectionUsers.InsertOne(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Endpoint hit: all user endpoint")

	json.NewEncoder(w).Encode(req)
}

// GetFriends adds a friend to a user in the db
func GetFriends(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFriends")
	// setupResponse(&w, r)
	// if (*r).Method == "OPTIONS" {
	// 	fmt.Println("GetFriends options")
	// 	return
	// }
	var req dbmodels.GetFriendsRequest
	json.NewDecoder(r.Body).Decode(&req)
	var result bson.M
	filter := bson.M{"name": req.Name}
	log.Println(filter)
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	log.Println("Found result: ", result)
	// insertResult, err := collectionMessages.InsertOne(context.TODO(), req)

	json.NewEncoder(w).Encode(result)
}

// GetWithFilter ...
func GetWithFilter(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.FriendRequest
	fmt.Println("GetWithFilter")
	json.NewDecoder(r.Body).Decode(&req)
	var result bson.M
	filter := bson.M{"name": req.User}
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Found result: ", result["friends"])
	// insertResult, err := collectionMessages.InsertOne(context.TODO(), req)

	json.NewEncoder(w).Encode(result["friends"])
}

// AddFriend ...
func AddFriend(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddFriends")

	// Get
	var getRequest dbmodels.FriendRequest
	// var getResult bson.M
	var getResult dbmodels.FriendResult
	json.NewDecoder(r.Body).Decode(&getRequest)
	getFilter := bson.M{"name": getRequest.User}
	err := collectionUsers.FindOne(context.TODO(), getFilter).Decode(&getResult)

	if err != nil {
		log.Println(err)
	}

	// Update
	var req dbmodels.FriendRequest
	var updateResult bson.M
	json.NewDecoder(r.Body).Decode(&req)
	var friendsSlice []string = getResult.Friends[0:]
	friendsSlice = append(friendsSlice, getRequest.NewFriend)

	fmt.Println(reflect.TypeOf(friendsSlice), friendsSlice)
	fmt.Println("Req.User: ", getRequest.User)
	// TODO: make this but with inserting it, not creating a new array
	updateFilter := bson.M{"name": getRequest.User}
	update, err := collectionUsers.UpdateOne(context.TODO(), updateFilter,
		// bson.D{
		//	{"$set", bson.D{{"friends", friendsSlice}}},
		// },
		// FRED
		bson.D{
			primitive.E{Key: "$set",
				Value: bson.D{primitive.E{Key: "friends", Value: friendsSlice}}}},
	)

	// example use of bson.D to avoid Linter warning of "primitive.E composite literal uses unkeyed fields" :
	// bson.D{primitive.E{Key: "$set", Value: bson.D{ primitive.E{Key: "name", Value: supplier.Name}, primitive.E{Key: "documentType", Value: supplier.DocumentType} }, }}

	fmt.Println(update)

	fmt.Println("Result: ", update.ModifiedCount)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Found result 1: ", updateResult["friends"])
	// insertResult, err := collectionMessages.InsertOne(context.TODO(), req)

	json.NewEncoder(w).Encode(updateResult["friends"])
}

// HomePage ...
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint HIt")
}

// func GetFriends(w http.ResponseWriter, r *http.Request) {
// 	var req dbmodels.UserModel
// 	fmt.Println("AddMessage")
// 	json.NewDecoder(r.Body).Decode(&req)
// 	insertResult, err := collectionUsers.FindOne(context.TODO(), req)
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
// 	fmt.Println("Endpoint hit: all user endpoint")
//
// 	json.NewEncoder(w).Encode(req)
// }

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
	collectionUsers = database.Collection("postusers")
}
