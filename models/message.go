package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

type Message struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content"`
	Email       string `json:"email" validate:"required,email"`
	MessageCode string
}

func ConnectDataBase() {
	database, err := gorm.Open("sqlite3", "messagebox.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&Message{})
	DB = database
}
