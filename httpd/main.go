package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/Jonny-exe/go-server/httpd/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"
)

func handleRequest() {
	handler.Connect()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/addmessage", handler.AddMessage).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/getfriends", handler.GetFriends).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/addfriend", handler.AddFriend).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/addfriendrequest", handler.AddFriendRequest).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/login", handler.Login).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/test", handler.Test).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/getwithfilter", handler.GetWithFilter).Methods("POST")
	myRouter.HandleFunc("/adduser", handler.AddUser).Methods("POST")
	myRouter.HandleFunc("/doesuserexist", handler.DoesUserExists).Methods("POST")
	myRouter.HandleFunc("/getfriendrequests", handler.GetFriendRequests).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5000"},
		AllowCredentials: true,
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
