package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Phone        string    `json:"phone"`
	Birthday     string    `json:"birthday"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRepo interface {
	Create(u *User) error
	GetByEmail(email string) (*User, error)
}

type gormUserRepo struct {
	db *gorm.DB
}

func NewGormUserRepo(db *gorm.DB) UserRepo {
	return &gormUserRepo{db: db}
}

func (r *gormUserRepo) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *gormUserRepo) GetByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
