package entity

import "time"

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Phone     string    `json:"phone"`
	Birthday  string    `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
}
