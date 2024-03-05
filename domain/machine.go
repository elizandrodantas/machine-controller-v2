package domain

import (
	"context"
	"time"
)

const (
	MACHINERESPONSE_SUCCESS     = "success"
	MACHINERESPONSE_CREATE      = "register"
	MACHINERESPONSE_UNAVAILABLE = "unavailable"
)

type Machine struct {
	ID    string `db:"id"`
	Guid  string `db:"guid"`
	Name  string `db:"name"`
	OS    string `db:"os"`
	Query string `db:"query"`

	CreatedAt time.Time `db:"created_at"`
}

type MachineData struct {
	Guid      string `json:"guid" binding:"required,uuid"`
	Name      string `json:"name" binding:"required,max=100"`
	OS        string `json:"os" binding:"required,max=100"`
	Expire    int    `json:"expire" binding:"required"`
	ServiceId string `json:"service_id" binding:"required,uuid4"`
}

type MachineRequest struct {
	Data string `json:"data" binding:"required,min=32,base64"`
}

type MachineUpdateNameRequest struct {
	Name string `json:"name" binding:"required,max=100"`
	Id   string `json:"id" binding:"required,uuid4"`
}

type MachineResponse struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	Identify        string `json:"identify"`
	ServiceIdentify string `json:"service_identify"`
	Name            string `json:"name"`
}

type MachineJson struct {
	ID    string `json:"id"`
	Guid  string `json:"guid"`
	Name  string `json:"name"`
	OS    string `json:"os"`
	Query string `json:"query"`

	CreatedAt time.Time `json:"created_at"`
}

type MachineListQuerys struct {
	MachineId string `form:"machine_id" json:"machine_id" binding:"omitempty,uuid4"`
	Query     string `form:"q" json:"q" binding:"omitempty,max=100"`
	Page      int    `form:"page" json:"page" binding:"omitempty"`
	OS        string `form:"os" json:"os" binding:"omitempty,max=100"`
}

type MachineListResponse struct {
	Total int           `json:"total"`
	Count int           `json:"count"`
	Page  int           `json:"page"`
	Data  []MachineJson `json:"data"`
}

type MachineRepository interface {
	Create(ctx context.Context, guid, name, os, query string) error
	FindByGuid(ctx context.Context, guid string) (*Machine, error)
	List(ctx context.Context, q MachineListQuerys) ([]Machine, error)
	FindById(ctx context.Context, id string) (*Machine, error)
	Count(ctx context.Context) (int, error)
	UpdateName(ctx context.Context, id, name string) error
}

type MachineUsecase interface {
	Create(ctx context.Context, md *MachineData) error
	FindByGuid(ctx context.Context, guid string) (*Machine, error)
	List(ctx context.Context, q MachineListQuerys) ([]MachineJson, error)
	Detail(ctx context.Context, id string) (*MachineJson, error)
	Count(ctx context.Context) (int, error)
	UpdateName(ctx context.Context, id, name string) error
}
