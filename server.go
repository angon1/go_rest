package main

import (
	"go-rest/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	messagesApi := router.PathPrefix("/api/messages").Subrouter()
	models.ConnectDataBase()
	defer models.DB.Close()
	router.HandleFunc("/", home)
	messagesApi.HandleFunc("/", messagesGet).Methods(http.MethodGet)
	messagesApi.HandleFunc("", messagesPost).Methods(http.MethodPost)
	messagesApi.HandleFunc("/{messageID}", messagesDelete).Methods(http.MethodDelete)
	messagesApi.HandleFunc("", home)
	log.Fatal(http.ListenAndServe(":8080", router))

}
