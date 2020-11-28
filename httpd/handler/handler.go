package handler

// skipping primitive.E
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jonny-exe/go-server/httpd/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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
	log.Println("AddMessage")
	json.NewDecoder(r.Body).Decode(&req)
	log.Println(req)
	model := dbmodels.MessageModel{Sender: req.Sender, Receiver: req.Receiver, Content: req.Content, Date: time.Now()}
	insertResult, err := collectionMessages.InsertOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	json.NewEncoder(w).Encode(model)
}

// AddUser adds a user to the db
func AddUser(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.FriendResult
	json.NewDecoder(r.Body).Decode(&req)
	insertResult, err := collectionUsers.InsertOne(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Endpoint hit: all user endpoint")

	json.NewEncoder(w).Encode(req)
}

// GetFriends gets all the friends from a user
func GetFriends(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFriends")
	var req dbmodels.GetFriendsRequest
	var result bson.M
	json.NewDecoder(r.Body).Decode(&req)
	log.Println("GetFriends: req ", req)
	filter := bson.M{"name": req.Name}
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	log.Println("GetFriends: Found result: ", result)
	json.NewEncoder(w).Encode(result)
}

func encryptPassword(password string) []byte {
	log.Println("Test")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Println(err)
	}
	return bytes
}

func comparePassword(dbpass []byte, pass string) bool {
	fmt.Println(dbpass, pass)
	err := bcrypt.CompareHashAndPassword(dbpass, []byte(pass))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// GetWithFilter gets all the messages from a certain user
func GetWithFilter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetWithFilter")
	var req dbmodels.GetWithFilterRequest
	json.NewDecoder(r.Body).Decode(&req)
	log.Println("This is GetWithFilter req: ", req)
	var result []bson.M
	// {"receiver": req.Receiver, "sender": req.Sender},
	// {"receiver": req.Sender, "sender": req.Receiver},

	filter := bson.D{
		primitive.E{
			Key: "$or", Value: []interface{}{
				bson.M{"receiver": req.Receiver, "sender": req.Sender},
				bson.M{"receiver": req.Sender, "sender": req.Receiver},
			}},
	}
	cursor, err := collectionMessages.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		log.Println(err)
	}
	fmt.Println("Found result: ", result)
	// insertResult, err := collectionMessages.InsertOne(context.TODO(), req)

	json.NewEncoder(w).Encode(result)
}

// DoesUserExists ...
func DoesUserExists(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.GetFriendsRequest
	json.NewDecoder(r.Body).Decode(&req)
	filter := bson.M{"name": req.Name}
	options := &options.CountOptions{}
	options.SetLimit(1)
	count, err := collectionUsers.CountDocuments(context.TODO(), filter, options)
	if err != nil {
		log.Println(err)
	}
	// var result string = "Hi"

	var result bool
	if count == 0 {
		result = false
	} else {
		result = true
	}
	json.NewEncoder(w).Encode(result)
}

// AddFriend ...
func AddFriend(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddFriends")

	// Get
	var getRequest dbmodels.FriendRequest
	var getResult dbmodels.FriendResult
	json.NewDecoder(r.Body).Decode(&getRequest)
	log.Println(getRequest)
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
	fmt.Println("FriendSlice: ", friendsSlice)
	updateFilter := bson.M{"name": getRequest.User}
	update, err := collectionUsers.UpdateOne(context.TODO(), updateFilter,
		// bson.D{
		//	{"$set", bson.D{{"friends", friendsSlice}}},
		// }
		bson.D{
			primitive.E{Key: "$set",
				Value: bson.D{primitive.E{Key: "friends", Value: friendsSlice}}}},
	)
	fmt.Println("Modified count: ", update.ModifiedCount)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(updateResult["friends"])
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
		log.Println(err)
	}
	fmt.Println("Connection succesfull")
	fmt.Println(cancel)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
	}
	// Dosent need to be closed, its better not to close i
	// defer client.Disconnect(ctx)

	database = client.Database("test")
	collectionMessages = database.Collection("postmessages")
	collectionUsers = database.Collection("postusers")
}
