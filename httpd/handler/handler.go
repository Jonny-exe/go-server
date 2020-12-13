package handler

// skipping primitive.E
import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Jonny-exe/go-server/httpd/defaultimage"
	"github.com/Jonny-exe/go-server/httpd/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/png"
	"strings"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"bytes"
	"github.com/anthonynsimon/bild/transform"
	_ "image/jpeg" // This is to load the jpeg encoder
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

// to install: "go get github.com/anthonynsimon/bild/..."
// installed it into: /home/a/go/src/github.com/anthonynsimon

// AddMessage adds users to mongodb database
func AddMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("AddMessage")
	var req dbmodels.Message
	json.NewDecoder(r.Body).Decode(&req)
	log.Println(req)
	model := dbmodels.MessageTime{Sender: req.Sender, Receiver: req.Receiver, Content: req.Content, Date: time.Now()}
	insertResult, err := collectionMessages.InsertOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	json.NewEncoder(w).Encode(model)
}

// Test ...
func Test(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Pass string `json:"pass"`
	}

	var req Request
	json.NewDecoder(r.Body).Decode(&req)
	pass := encryptPassword(req.Pass)

	json.NewEncoder(w).Encode(pass)
}

// AddUser adds a user to the db
func AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("AddUser")
	var req dbmodels.User
	json.NewDecoder(r.Body).Decode(&req)
	log.Println("AddUser: req: ", req)
	encryptedPass := encryptPassword(req.Pass)
	req.Pass = string(encryptedPass)
	req.ProfileImage = defaultimage.Image
	insertResult, err := collectionUsers.InsertOne(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Endpoint hit: all user endpoint")

	json.NewEncoder(w).Encode(req)
}

// GetFriendRequests ...
func GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFriendRequest")
	var req dbmodels.Name
	json.NewDecoder(r.Body).Decode(&req)
	var result dbmodels.NameAndDateStruct
	filter := bson.M{"name": req.Name}
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Println("GetFriendRequests: result: ", result)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(result.FriendRequests)
}

// UploadProfileImage ..
func UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.NameAndImageAndAreaToCrop
	log.Println("UploadProfileImage", req)
	json.NewDecoder(r.Body).Decode(&req)
	log.Println("UploadProfileImage", req.Name, req.AreaToCrop)
	image := editImage(req.Image, req.AreaToCrop)
	updateFilter := bson.M{"name": req.Name}

	update, err := collectionUsers.UpdateOne(context.TODO(), updateFilter,
		// bson.D{
		//	{"$set", bson.D{{"profileImage", friendsSlice}}},
		// }
		bson.D{
			primitive.E{Key: "$set",
				Value: bson.D{primitive.E{Key: "profileimage", Value: image}}}},
	)

	if err != nil {
		log.Println(err)
	}

	log.Println(update)
	json.NewEncoder(w).Encode(http.StatusOK)
}

func editImage(base64Image string, cropSizes dbmodels.CropSizes) string {
	// this file read is just for testing, the server does not need this
	// see https://golang.org/pkg/encoding/base64/
	var dataenc string
	var datadec []byte

	// Data to retrun
	var datacropenc string
	dataenc = base64Image
	dataenc = dataenc[strings.IndexByte(dataenc, ',')+1:]
	// argument data is []byte, exactly what we need
	// dataenc is what the server receives in the REST API call
	// a base64 encoded image of any type (JPG, PNG) in a string
	fmt.Printf("Contents of file (encoded): %v ...\n", string(dataenc)[0:127])
	datadec, err := base64.StdEncoding.DecodeString(dataenc)
	if err != nil {
		log.Println("Decode failed with error: ", err)
		log.Println(err)
	} else {
		log.Printf("Contents of file (decoded): %v ...\n", string(datadec)[0:127])
		// create an io.Reader for []byte
		r := bytes.NewReader(datadec)
		log.Println("Created io.reader for []byte")
		// Calling the generic image.Decode() will convert the bytes into an image
		// and give us type of image as a string. E.g. "png"
		imageData, imageType, err := image.Decode(r)
		if err != nil {
			fmt.Println("Decoding image failed with error: ", err)
			fmt.Println(err)
		} else {
			// fmt.Println(imageData) // imageData (type image.Image)
			fmt.Printf("Image type is '%v'.\n", imageType)
			// use the coordinates received by server via REST API
			cropRect := image.Rect(cropSizes.X, cropSizes.Y, cropSizes.X+cropSizes.Width, cropSizes.Y+cropSizes.Height) // image.Rect(x0, y0, x1, y1)
			imageCropped := transform.Crop(imageData, cropRect)

			// Resize
			imageCropped = transform.Resize(imageCropped, 256, 256, transform.Linear)

			var b bytes.Buffer
			w := io.Writer(&b) // create a byte[] io.Writer using buffer
			// Encode takes a writer interface and an image interface
			// Since we want a PNG as output we convert to PNG
			png.Encode(w, imageCropped)
			imageCroppedString := b.String()
			imageCroppedBytes := b.Bytes()
			fmt.Printf("Contents of cropped image (unencoded): %v ...\n", imageCroppedString[0:127])
			if err != nil {
				fmt.Println("Write cropped image to file for testing failed with error: ", err)
				fmt.Println(err)
			} else {
				// argument imageCroppedBytes is []byte, exactly what we need
				datacropenc = base64.StdEncoding.EncodeToString(imageCroppedBytes) // []byte(arg)
				datacropenc = "data:image/jpeg;base64," + datacropenc
				// datacropenc is what the server will store in MongoDB,
				// a cropped PNG image that is base64 encoded stored in a string.
				fmt.Printf("Contents of cropped PNG image (encoded): %v ...\n", string(datacropenc)[0:127])
				fmt.Printf("Crop completed successfully.\n")
				// data:image/jpeg;base64,
			}
		}
	}
	return datacropenc
}

// Base64ToImage ...
func Base64ToImage(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Data string `json:"data"`
	}
	var req request
	json.NewDecoder(r.Body).Decode(&req)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(req.Data))
	log.Println(reader)
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)
	log.Println(reflect.TypeOf(m), reflect.TypeOf(formatString))
	json.NewEncoder(w).Encode(http.StatusOK)
}

// GetProfileImage ...
func GetProfileImage(w http.ResponseWriter, r *http.Request) {
	log.Println("GetProfileImage")
	var req dbmodels.NameAndImage
	json.NewDecoder(r.Body).Decode(&req)
	filter := bson.M{"name": req.Name}
	var result dbmodels.Image

	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	// log.Println(result)

	json.NewEncoder(w).Encode(result.ProfileImage)
}

// func UploadImage(w http.ResponseWriter, r *http.Request) {
// 	r.ParseMultipartForm(10 << 20)
// 	file, handler, err := r.FormFile("image")
//
// 	if err != nil {
// 		log.Println(err, file, handler)
// 		json.NewEncoder(w).Encode(err)
// 		return
// 	}
//
// 	defer file.Close()
// 	log.Println("UploadImage: filename", handler.Filename)
// 	log.Println("UploadImage: file size", handler.Size)
// 	log.Println("UploadImage: MIME header", handler.Header)
//
// 	tempFile, err := ioutil.TempFile("temp-images", "*")
//
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer file.Close()
//
// 	fileBytes, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	tempFile.Write(fileBytes)
// 	json.NewEncoder(w).Encode("Everything worked")
//
// }
//

// GetFriends gets all the friends from a user
func GetFriends(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFriends")
	var req dbmodels.Name
	var result dbmodels.Friends
	json.NewDecoder(r.Body).Decode(&req)
	log.Println("GetFriends: req ", req)
	filter := bson.M{"name": req.Name}
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	log.Println("GetFriends: Found result: ", result)
	json.NewEncoder(w).Encode(result.Friends)
}

func encryptPassword(password string) []byte {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Println(err)
	}
	return bytes
}

// Login ..
func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login")
	var req dbmodels.NameAndPass
	json.NewDecoder(r.Body).Decode(&req)
	log.Println("Login: req: ", req)
	if req.Pass == "" {
		json.NewEncoder(w).Encode(false)
		log.Println("Login: returned false because empty password")
		return
	}
	type Search struct {
		Pass string `json:"pass"`
	}
	var result Search
	log.Println(reflect.TypeOf(result))
	log.Println("Login", req)
	json.NewDecoder(r.Body).Decode(&req)
	filter := bson.M{"name": req.Name}
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	var dbpass string = result.Pass
	if err != nil {
		log.Println(err)
	}

	compareResult := comparePassword([]byte(dbpass), req.Pass)
	json.NewEncoder(w).Encode(compareResult)
}

func comparePassword(dbpass []byte, pass string) bool {
	log.Println(dbpass, pass)
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
	var req dbmodels.SenderAndReceiver
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
	var req dbmodels.Name
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

func addNewFriend(user string, newFriend string) {
	// Get current friends
	var userInfo dbmodels.Friends
	getFilter := bson.M{"name": user}
	err := collectionUsers.FindOne(context.TODO(), getFilter).Decode(&userInfo)

	if err != nil {
		log.Println(err)
	}

	// Update to new friends
	friendsSlice := appendToArray(newFriend, userInfo.Friends)

	fmt.Println("addNewFriend: FriendSlice: ", friendsSlice)
	updateFilter := bson.M{"name": user}
	update, err := collectionUsers.UpdateOne(context.TODO(), updateFilter,
		// bson.D{
		//	{"$set", bson.D{{"friends", friendsSlice}}},
		// }
		bson.D{
			primitive.E{Key: "$set",
				Value: bson.D{primitive.E{Key: "friends", Value: friendsSlice}}}},
	)
	fmt.Println("addNewFriend: Modified count: ", update.ModifiedCount)
	if err != nil {
		log.Println(err)
	}
}

// AddFriend ...
func AddFriend(w http.ResponseWriter, r *http.Request) {
	log.Println("AddFriends")

	var req dbmodels.NameAndNewFriend
	json.NewDecoder(r.Body).Decode(&req)

	// Add eachother to the friend list
	addNewFriend(req.Name, req.NewFriend)
	addNewFriend(req.NewFriend, req.Name)

	// Remove the friend request
	removeFriendRequests(req.Name, req.NewFriend)
	json.NewEncoder(w).Encode(http.StatusOK)
}

// AddFriendRequest ...
func AddFriendRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("AddFriendRequest")
	var req dbmodels.NameAndNewFriend
	json.NewDecoder(r.Body).Decode(&req)
	addFriendRequest(req.Name, req.NewFriend)
	json.NewEncoder(w).Encode(http.StatusOK)
}

func addFriendRequest(user string, newFriend string) {
	// Get current requests
	log.Println("AddFriendRequest", user, newFriend)
	filter := bson.M{"name": newFriend}
	var result dbmodels.NameAndDateStruct
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	}

	// Update current requests
	var newFriendRequest dbmodels.NameAndDate
	newFriendRequest.Name = user
	newFriendRequest.Date = time.Now().Format("2006-01-02")
	var requestsSlice dbmodels.NameAndDateArray = result.FriendRequests[0:]
	requestsSlice = append(requestsSlice, newFriendRequest)
	updateFilter := bson.M{"name": newFriend}
	update, err := collectionUsers.UpdateOne(context.TODO(), updateFilter,
		// bson.D{
		//	{"$set", bson.D{{"friends", friendsSlice}}},
		// }
		bson.D{
			primitive.E{Key: "$set",
				Value: bson.D{primitive.E{Key: "friendrequests", Value: requestsSlice}}}},
	)
	if err != nil {
		log.Println(err)
	}
	log.Println(update)
}

// RemoveFriendRequest ...
func RemoveFriendRequest(w http.ResponseWriter, r *http.Request) {
	var req dbmodels.NameAndFriendToRemove
	json.NewDecoder(r.Body).Decode(&req)
	removeFriendRequests(req.Name, req.FriendToRemove)
	json.NewEncoder(w).Encode(http.StatusOK)
}

func removeFriendRequests(user string, requestToRemove string) {
	var findResult dbmodels.NameAndDateStruct
	// var result = bson.M{}
	filter := bson.M{"name": user}
	log.Println("removeFriendRequests: filter", filter)
	err := collectionUsers.FindOne(context.TODO(), filter).Decode(&findResult)

	if err != nil {
		log.Println(err)
	}

	log.Println("removeFriendRequests: findResult.FriendRequests", findResult)
	var indexToRemove int = findFriendToRemove(findResult.FriendRequests, requestToRemove)
	var newRequests dbmodels.NameAndDateArray = removeIndexFromArray(findResult.FriendRequests, indexToRemove)

	updateFilter := bson.M{"name": user}
	update, err := collectionUsers.UpdateOne(context.TODO(), updateFilter,
		// bson.D{
		//	{"$set", bson.D{{"friends", friendsSlice}}},
		// }
		bson.D{
			primitive.E{Key: "$set",
				Value: bson.D{primitive.E{Key: "friendrequests", Value: newRequests}}}},
	)

	if err != nil {
		log.Println(err, update)
	}
}

func findFriendToRemove(requests dbmodels.NameAndDateArray, friendName string) int {
	var requestIndex int
	for index, request := range requests {
		if request.Name == friendName {
			requestIndex = index
			break
		}
	}
	log.Println("findFriendToRemove: ", requestIndex)
	return requestIndex
}

func removeIndexFromArray(requests dbmodels.NameAndDateArray, index int) dbmodels.NameAndDateArray {
	log.Println("removeIndexFromArray: ", index, requests)
	return append(requests[:index], requests[index+1:]...)
}

func appendToArray(newElement string, array []string) []string {
	var newSlice []string = array[0:]
	newSlice = append(newSlice, newElement)
	return newSlice
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
