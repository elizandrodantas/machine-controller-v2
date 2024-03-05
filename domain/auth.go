package domain

import (
	"context"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=150"`
}

type AuthResponse struct {
	Token  string `json:"access_token"`
	Type   string `json:"typen_type"`
	Expire int64  `json:"expire"`
}

type AuthUsecase interface {
	Autenticate(c context.Context, username, password string) (*User, error)
}
