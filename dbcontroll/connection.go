package dbcontroll

import (
	"fmt"
	"go-rest/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func ConnectDataBase() {
	database, err := gorm.Open("sqlite3", "messagebox.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&models.Message{})
	DB = database
}
