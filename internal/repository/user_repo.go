package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	MemberCode   string    `gorm:"uniqueIndex;default:null" json:"member_code"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Phone        string    `json:"phone"`
	Birthday     string    `json:"birthday"`
	CreatedAt    time.Time `json:"created_at"`
	Points       int       `json:"points" gorm:"default:0"`
}

type UserRepo interface {
	Create(u *User) error
	GetByEmail(email string) (*User, error)
	GetByMemberCode(code string) (*User, error)
	Update(u *User) error
	CreateTransaction(t *Transaction) error
	GetRecentRecipientsBySender(senderEmail string, limit int) ([]User, error)
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

func (r *gormUserRepo) GetByMemberCode(code string) (*User, error) {
	var u User
	if err := r.db.Where("member_code = ?", code).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormUserRepo) Update(u *User) error {
	return r.db.Save(u).Error
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FromEmail string    `json:"from_email"`
	ToEmail   string    `json:"to_email"`
	ToMember  string    `json:"to_member_code"`
	Amount    int       `json:"amount"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *gormUserRepo) CreateTransaction(t *Transaction) error {
	return r.db.Create(t).Error
}

// Get recent recipients for a sender (distinct by to_email) ordered by last transfer desc
func (r *gormUserRepo) GetRecentRecipientsBySender(senderEmail string, limit int) ([]User, error) {
	var users []User
	// join transactions to users to get recipient info
	rows, err := r.db.Raw(`SELECT u.* FROM user_repos u JOIN (SELECT to_email, MAX(created_at) as maxt FROM transactions WHERE from_email = ? GROUP BY to_email ORDER BY maxt DESC LIMIT ?) t ON u.email = t.to_email`, senderEmail, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		if err := r.db.ScanRows(rows, &u); err == nil {
			users = append(users, u)
		}
	}
	return users, nil
}
