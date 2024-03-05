package domain

import (
	"context"
	"time"
)

type User struct {
	ID       string   `db:"id"`
	Name     string   `db:"name"`
	Username string   `db:"username"`
	Password string   `db:"password"`
	Status   bool     `db:"status"`
	Scope    []string `db:"scope"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserJson struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Status   bool     `json:"status"`
	Scope    []string `json:"scope"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCreate struct {
	Name       string
	Username   string
	Password   string
	Scope      []string
	RegisterId string
}

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"alphanum,max=50"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=150"`
}

type UserListResponse struct {
	Total int        `json:"total"`
	Count int        `json:"count"`
	Page  int        `json:"page"`
	Data  []UserJson `json:"data"`
}

type UserChangePasswordRequest struct {
	NewPassword string `json:"password" binding:"required,min=6,max=150"`
}

type UserChangePasswordResponse struct {
	Message string `json:"message"`
}

type UserCremoveScopeRequest struct {
	Scope  string `json:"scope" binding:"required"`
	UserId string `json:"user_id" binding:"required,uuid"`
}

type UserRepository interface {
	Create(c context.Context, u *User) error
	FindByUsername(c context.Context, username string) (*User, error)
	UpdatePassword(c context.Context, id, password string) error
	FindById(c context.Context, id string) (*User, error)
	List(c context.Context, page int) ([]User, error)
	AddScope(c context.Context, id, s string) error
	RemoveScope(c context.Context, id, s string) error
	ChangeStatus(c context.Context, id string) error
	Count(c context.Context) (int, error)
	Delete(c context.Context, id string) error
}

type UserUsecase interface {
	Create(c context.Context, userCreate *UserCreate) error
	FindByUsername(c context.Context, username string) (*User, error)
	FindById(c context.Context, id string) (*User, error)
	UpdatePassword(c context.Context, id, password string) error
	ExistUsername(c context.Context, username string) bool
	List(c context.Context, page int) ([]UserJson, error)
	AddScopeWithUserId(c context.Context, id, s string) error
	RemoveScopeWithUserId(c context.Context, id, s string) error
	Disable(c context.Context, id string) error
	Enable(c context.Context, id string) error
	Count(c context.Context) (int, error)
	DeleteById(c context.Context, id string) error
	DeleteByUsername(c context.Context, username string) error
}
