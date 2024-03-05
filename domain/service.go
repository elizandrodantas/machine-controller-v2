package domain

import (
	"context"
	"time"
)

type Service struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type ServiceJson struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServiceCreateRequest struct {
	Name        string `json:"name" binding:"required,gte=3,lte=50"`
	Description string `json:"description"`
}

type ServiceListResponse struct {
	Count int           `json:"count"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Data  []ServiceJson `json:"data"`
}

type ServiceRemoveRequest struct {
	ID string `json:"id" binding:"required,uuid"`
}

type ServiceRepository interface {
	Create(ctx context.Context, name, description string) error
	List(ctx context.Context, page int) ([]Service, error)
	Remove(ctx context.Context, id string) error
	FindByName(ctx context.Context, name string) (*Service, error)
	FindById(ctx context.Context, id string) (*Service, error)
	Count(ctx context.Context) (int, error)
}

type ServiceUsecase interface {
	Create(ctx context.Context, name, description string) error
	List(ctx context.Context, page int) ([]ServiceJson, error)
	Remove(ctx context.Context, id string) error
	FindByName(ctx context.Context, name string) (*Service, error)
	FindById(ctx context.Context, id string) (*Service, error)
	Count(ctx context.Context) (int, error)
}
