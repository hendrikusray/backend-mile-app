package domain

import (
	"context"
	"time"
)

type User struct {
	UserId      string    `json:"user_id"`
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Address     string    `json:"address" binding:"required"`
	Pin         string    `json:"pin" binding:"required"`
	Username    string    `json:"username" binding:"required"`
	Password    string    `json:"password" binding:"required"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserUseCase interface {
	Login(ctx context.Context, login *Login) (Token, error)
}

type UserRepository interface {
	GetUser(ctx context.Context, phoneNumber string) (*User, error)
}
