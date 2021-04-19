package main

import (
	"go-rest/controllers"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	messagesApi := router.PathPrefix("/api/messages").Subrouter()
	controllers.ConnectDb()
	defer controllers.CloseDb()
	router.HandleFunc("/", controllers.Home)
	messagesApi.HandleFunc("/", controllers.MessagesGet).Methods(http.MethodGet)
	messagesApi.HandleFunc("", controllers.MessagesPost).Methods(http.MethodPost)
	messagesApi.HandleFunc("/{messageID}", controllers.MessagesDelete).Methods(http.MethodDelete)
	messagesApi.HandleFunc("", controllers.Home)
	log.Fatal(http.ListenAndServe(":8080", router))

}
