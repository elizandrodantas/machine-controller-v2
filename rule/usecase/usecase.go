package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
)

type ruleUsecase struct {
	RuleRepository domain.RuleRepository
}

func NewRuleUsecase(r domain.RuleRepository) domain.RuleUsecase {
	return &ruleUsecase{
		RuleRepository: r,
	}
}

func (r *ruleUsecase) ActiveByMachineIdAndServiceId(ctx context.Context, machineId, serviceId string) (*domain.MachineRule, error) {
	return r.RuleRepository.ActiveByMachineIdAndServiceId(ctx, machineId, serviceId)
}

func (r *ruleUsecase) Create(ctx context.Context, machineId, serviceId string, exp int) error {
	if exp <= 0 {
		return fmt.Errorf("expire time must be greater than 0")
	}

	timeNow := time.Now()
	expire := 0

	// IF CLIENT HAS A SERVICE ACTIVE, ADD THE EXPIRE TIME TO THE NEW EXPIRE TIME
	if serviceRule, err := r.ActiveByMachineIdAndServiceId(ctx, machineId, serviceId); err == nil {
		diff := serviceRule.Expire - int(timeNow.Unix())
		if diff > 0 {
			expire += diff
		}
	}

	if err := r.Remove(ctx, machineId, serviceId); err != nil {
		return fmt.Errorf("error to invalidate: %w", err)
	}

	expire += int(timeNow.Unix()) + exp

	return r.RuleRepository.Create(ctx, machineId, serviceId, expire)
}

func (r *ruleUsecase) History(ctx context.Context, machineId string, serviceId string, page int) ([]domain.MachineRuleJson, error) {
	list, err := r.RuleRepository.History(ctx, machineId, serviceId, page)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []domain.MachineRuleJson{}, nil
	}

	var response = make([]domain.MachineRuleJson, 0, len(list))
	for _, item := range list {
		response = append(response, domain.MachineRuleJson(item))
	}

	return response, nil
}

func (r *ruleUsecase) List(ctx context.Context, page int, oact bool) ([]domain.RuleJoinServiceJson, error) {
	list, err := r.RuleRepository.List(ctx, page, oact)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []domain.RuleJoinServiceJson{}, nil
	}

	var response = make([]domain.RuleJoinServiceJson, 0, len(list))
	for _, item := range list {
		response = append(response, domain.RuleJoinServiceJson(item))
	}

	return response, nil
}

func (r *ruleUsecase) Count(ctx context.Context) (int, error) {
	return r.RuleRepository.Count(ctx)
}

func (r *ruleUsecase) CountByMachineId(ctx context.Context, machineId string, serviceId string) (int, error) {
	return r.RuleRepository.CountByMachineId(ctx, machineId, serviceId)
}

func (r *ruleUsecase) Actives(ctx context.Context, machineId string) ([]domain.RuleJoinServiceJson, error) {
	actives, err := r.RuleRepository.Actives(ctx, machineId)
	if err != nil {
		return nil, err
	}

	if len(actives) == 0 {
		return []domain.RuleJoinServiceJson{}, nil
	}

	var response = make([]domain.RuleJoinServiceJson, 0, len(actives))
	for _, item := range actives {
		response = append(response, domain.RuleJoinServiceJson(item))
	}

	return response, nil
}

func (r *ruleUsecase) Remove(ctx context.Context, machineId, serviceId string) error {
	return r.RuleRepository.Invalidate(ctx, machineId, serviceId)
}

func (r *ruleUsecase) CountOnlyActives(ctx context.Context) (int, error) {
	return r.RuleRepository.CountOnlyActives(ctx)
}
