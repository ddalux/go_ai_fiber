package usecase

import (
	"errors"
	"time"

	"github.com/ddalux/go_ai_fiber/internal/auth"
	repo "github.com/ddalux/go_ai_fiber/internal/repository"
	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
)

type UserUsecase interface {
	Register(email, password, firstname, lastname, phone, birthday string) error
	Login(email, password string) (string, error)
	Me(token string) (*repo.User, error)
}

type userUsecase struct {
	r         repo.UserRepo
	jwtSecret []byte
}

func NewUserUsecase(r repo.UserRepo, jwtSecret []byte) UserUsecase {
	return &userUsecase{r: r, jwtSecret: jwtSecret}
}

func (u *userUsecase) Register(email, password, firstname, lastname, phone, birthday string) error {
	if _, err := u.r.GetByEmail(email); err == nil {
		return ErrUserExists
	}
	// simple hash using sha256 (demo). Replace with bcrypt in production.
	hash := auth.SimpleHash(password)
	user := &repo.User{
		Email:        email,
		PasswordHash: hash,
		FirstName:    firstname,
		LastName:     lastname,
		Phone:        phone,
		Birthday:     birthday,
		CreatedAt:    time.Now(),
	}
	return u.r.Create(user)
}

func (u *userUsecase) Login(email, password string) (string, error) {
	user, err := u.r.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCreds
	}
	if !auth.VerifySimpleHash(user.PasswordHash, password) {
		return "", ErrInvalidCreds
	}
	claims := jwt.RegisteredClaims{
		Subject:   user.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", err
	}
	return ts, nil
}

func (u *userUsecase) Me(tokenStr string) (*repo.User, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return u.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidCreds
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidCreds
	}
	user, err := u.r.GetByEmail(claims.Subject)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
