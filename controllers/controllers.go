package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go-rest/dbcontroll"
	"go-rest/models"

	"net/http"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func MessagesGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	email := request.URL.Query().Get("email")
	var msgs []models.Message
	if email != "" {
		dbcontroll.DB.Where("email = ?", email).Find(&msgs)
	} else {
		dbcontroll.DB.Find(&msgs)
	}
	fmt.Fprintf(writer, "message: %+v email: %s", msgs, email)
}

func generateRandomCode(msg models.Message) (md5Hash string) {
	tempStr := msg.Title + msg.Content + msg.Email + time.Now().String()
	messageString := []byte(tempStr)
	md5Hash = fmt.Sprintf("%x", md5.Sum(messageString))
	return md5Hash
}

func createMessage(msg models.Message) (hash string) {
	hash = generateRandomCode(msg)
	msg.MessageCode = hash
	dbcontroll.DB.Create(&models.Message{
		ID:          0,
		Title:       msg.Title,
		Content:     msg.Content,
		Email:       msg.Email,
		MessageCode: msg.MessageCode,
	})
	return hash
}

func MessagesPost(writer http.ResponseWriter, request *http.Request) {
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

func MessagesDelete(writer http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	messageID := pathParams["messageID"]
	var msg models.Message
	result := dbcontroll.DB.Where("message_code = ?", messageID).First(&msg)
	if result.RowsAffected != 0 {
		dbcontroll.DB.Delete(&msg)
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func Home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(`{"message": "not found"}`))
}

func ConnectDb() {
	dbcontroll.ConnectDataBase()
}

func CloseDb() {
	dbcontroll.DB.Close()
}
