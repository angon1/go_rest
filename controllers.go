package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go-rest/models"

	"net/http"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func messagesGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	email := request.URL.Query().Get("email")
	var msgs []models.Message
	if email != "" {
		models.DB.Where("email = ?", email).Find(&msgs)
	} else {
		models.DB.Find(&msgs)
	}
	fmt.Fprintf(writer, "message: %+v email: %s", msgs, email)
}

func GenerateRandomCode(msg models.Message) (md5Hash string) {
	tempStr := msg.Title + msg.Content + msg.Email + time.Now().String()
	messageString := []byte(tempStr)
	md5Hash = fmt.Sprintf("%x", md5.Sum(messageString))
	return md5Hash
}

func createMessage(msg models.Message) (hash string) {
	hash = GenerateRandomCode(msg)
	msg.MessageCode = hash
	models.DB.Create(&models.Message{
		ID:          0,
		Title:       msg.Title,
		Content:     msg.Content,
		Email:       msg.Email,
		MessageCode: msg.MessageCode,
	})
	return hash
}

func messagesPost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var msg models.Message
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	jsonValidator := validator.New()
	err = jsonValidator.Struct(msg)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, `{"errors":"%v"}`, validationErrors)
	} else {
		writer.WriteHeader(http.StatusCreated)
		MessageCode := createMessage(msg)
		writer.Write([]byte(fmt.Sprintf(`{"messageCode":"%s"}`, MessageCode)))
	}
}

func messagesDelete(writer http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	messageID := pathParams["messageID"]
	var msg models.Message
	result := models.DB.Where("message_code = ?", messageID).First(&msg)
	if result.RowsAffected != 0 {
		models.DB.Delete(&msg)
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(`{"message": "not found"}`))
}
