package domain

import (
	"context"
	"time"
)

type Notes struct {
	ID        string    `db:"id"`
	Text      string    `db:"text"`
	MachineId string    `db:"machine_id"`
	UserId    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type NotesList struct {
	ID        string    `db:"id"`
	Text      string    `db:"text"`
	MachineId string    `db:"machine_id"`
	UserId    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserName  string    `db:"user_name"`
}

type NotesListJson struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	MachineId string    `json:"machine_id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
}

type NotesCreateRequest struct {
	Text      string `json:"text" binding:"required,min=1,max=5000"`
	MachineId string `json:"machine_id" binding:"required,uuid"`
}

type NotesListResponse struct {
	Total int             `json:"total"`
	Count int             `json:"count"`
	Page  int             `json:"page"`
	Data  []NotesListJson `json:"data"`
}

type NotesUpdateRequest struct {
	Text string `json:"text" binding:"required"`
	ID   string `json:"id" binding:"required,uuid"`
}

type NotesDeleteRequest struct {
	ID string `json:"id" binding:"required,uuid"`
}

type NotesRepository interface {
	Create(c context.Context, text, machineId, userId string) error
	FindById(c context.Context, id string) (Notes, error)
	ListByMachineId(c context.Context, machineId string, page int) ([]NotesList, error)
	ChangeNote(c context.Context, id, text string) error
	Delete(c context.Context, id string) error
	CountByMachineId(c context.Context, machineId string) (int, error)
}

type NotesUsecase interface {
	Create(c context.Context, text, machineId, userId string) error
	ListByMachineId(c context.Context, machineId string, page int) ([]NotesListJson, error)
	ChangeNote(c context.Context, id, userId, text string) error
	Delete(c context.Context, userId, noteId string) error
	CountByMachineId(c context.Context, machineId string) (int, error)
}
