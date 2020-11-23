package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/Jonny-exe/go-server/httpd/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Article ...
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint HIt")
}

func handleRequest() {
	handler.Connect()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", handler.HomePage).Methods("POST")
	myRouter.HandleFunc("/addmessage", handler.AddMessage).Methods("POST")
	myRouter.HandleFunc("/getfriends", handler.GetFriends).Methods("POST")
	myRouter.HandleFunc("/adduser", handler.AddUser).Methods("POST")
	// myRouter.HandleFunc("/getfriends", handler.GetFriends).Methods("POST")

	fmt.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main() {
	handleRequest()
}
