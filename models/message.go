package models

type Message struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content"`
	Email       string `json:"email" validate:"required,email"`
	MessageCode string
}
