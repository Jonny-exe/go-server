package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Jonny-exe/go-server/httpd/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func handleRequest() {
	handler.Connect()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/test", handler.Base64ToImage).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/addmessage", handler.AddMessage).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/getfriends", handler.GetFriends).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/addfriend", handler.AddFriend).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/addfriendrequest", handler.AddFriendRequest).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/removefriendrequest", handler.RemoveFriendRequest).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/login", handler.Login).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/test", handler.Test).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/getwithfilter", handler.GetWithFilter).Methods("POST")
	myRouter.HandleFunc("/adduser", handler.AddUser).Methods("POST")
	myRouter.HandleFunc("/doesuserexist", handler.DoesUserExists).Methods("POST")
	myRouter.HandleFunc("/getfriendrequests", handler.GetFriendRequests).Methods("POST")
	myRouter.HandleFunc("/uploadprofileimage", handler.UploadProfileImage).Methods("POST")
	myRouter.HandleFunc("/getprofileimage", handler.GetProfileImage).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://jonny.sytes.net", "http://192.168.0.19"},
		AllowCredentials: false,

		// Enable Debugging for testing, consider disabling in production
		// To debug turn this to true
		Debug: false,
	})

	var PORT int = 5000
	corsHandler := c.Handler(myRouter)
	fmt.Println("Listening on port: ", PORT)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), corsHandler))

}

func main() {
	handleRequest()
}
