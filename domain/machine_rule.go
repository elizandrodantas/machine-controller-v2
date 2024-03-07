package domain

import (
	"context"
	"time"
)

type MachineRule struct {
	MachineId string    `db:"machine_id"`
	ServiceId string    `db:"service_id"`
	Expire    int       `db:"expire"`
	CreatedAt time.Time `db:"created_at"`
}

type MachineRuleJson struct {
	MachineId string    `json:"machine_id"`
	ServiceId string    `json:"service_id"`
	Expire    int       `json:"expire"`
	CreatedAt time.Time `json:"created_at"`
}

type RuleCreateRequest struct {
	MachineId string `json:"machine_id" validate:"required,uuid4"`
	ServiceId string `json:"service_id" validate:"required,uuid4"`
	Expire    int    `json:"expire" validate:"required,gte=0"`
}

type RuleCreateResponse struct {
	MachineId string `json:"machine_id"`
	ServiceId string `json:"service_id"`
	Expire    int    `json:"expire"`
}

type RuleHistoryResponse struct {
	MachineId string            `json:"machine_id"`
	Total     int               `json:"total"`
	Count     int               `json:"count"`
	Page      int               `json:"page"`
	Data      []MachineRuleJson `json:"data"`
}

type RuleJoinService struct {
	ServiceName string    `db:"service_name"`
	MachineId   string    `db:"machine_id"`
	ServiceId   string    `db:"service_id"`
	Expire      int       `db:"expire"`
	CreatedAt   time.Time `db:"created_at"`
}

type RuleJoinServiceJson struct {
	ServiceName string    `json:"service_name"`
	MachineId   string    `json:"machine_id"`
	ServiceId   string    `json:"service_id"`
	Expire      int       `json:"expire"`
	CreatedAt   time.Time `json:"created_at"`
}

type RuleActivesResponse struct {
	MachineId string                `json:"machine_id"`
	Total     int                   `json:"total"`
	Data      []RuleJoinServiceJson `json:"data"`
}

type RuleListResponse struct {
	Total int                   `json:"total"`
	Count int                   `json:"count"`
	Page  int                   `json:"page"`
	Data  []RuleJoinServiceJson `json:"data"`
}

type RuleRemoveRequest struct {
	MachineId string `json:"machine_id" validate:"required,uuid4"`
	ServiceId string `json:"service_id" validate:"required,uuid4"`
}

type RuleUsecase interface {
	Create(ctx context.Context, machineId, serviceId string, expire int) error
	History(ctx context.Context, machineId string, serviceId string, page int) ([]MachineRuleJson, error)
	Count(ctx context.Context) (int, error)
	CountByMachineId(ctx context.Context, machineId string, serviceId string) (int, error)
	List(ctx context.Context, page int) ([]RuleJoinServiceJson, error)
	Remove(ctx context.Context, machineId, serviceId string) error
	Actives(ctx context.Context, machineId string) ([]RuleJoinServiceJson, error)
	ActiveByMachineIdAndServiceId(ctx context.Context, machineId, serviceId string) (*MachineRule, error)
}

type RuleRepository interface {
	Create(ctx context.Context, machineId, serviceId string, expire int) error
	History(ctx context.Context, machineId string, serviceId string, page int) ([]MachineRule, error)
	Count(ctx context.Context) (int, error)
	CountByMachineId(ctx context.Context, machineId string, serviceId string) (int, error)
	List(ctx context.Context, page int) ([]RuleJoinService, error)
	Invalidate(ctx context.Context, machineId, serviceId string) error
	Actives(ctx context.Context, machineId string) ([]RuleJoinService, error)
	ActiveByMachineIdAndServiceId(ctx context.Context, machineId, serviceId string) (*MachineRule, error)
}
