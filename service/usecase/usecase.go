package usecase

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5"
)

type usecaseService struct {
	ServiceRepository domain.ServiceRepository
}

func NewServiceUsecase(serviceRepository domain.ServiceRepository) domain.ServiceUsecase {
	return &usecaseService{serviceRepository}
}

func (u *usecaseService) Create(ctx context.Context, name, description string) error {
	return u.ServiceRepository.Create(ctx, name, description)
}

func (u *usecaseService) FindByName(ctx context.Context, name string) (*domain.Service, error) {
	return u.ServiceRepository.FindByName(ctx, name)
}

func (u *usecaseService) List(ctx context.Context, page int) ([]domain.ServiceJson, error) {
	list, err := u.ServiceRepository.List(ctx, page)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []domain.ServiceJson{}, nil
	}

	var services = make([]domain.ServiceJson, 0, len(list))
	for _, service := range list {
		services = append(services, domain.ServiceJson(service))
	}

	return services, nil
}

func (u *usecaseService) Remove(ctx context.Context, id string) error {
	if _, err := u.FindById(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}

		return err
	}

	return u.ServiceRepository.Remove(ctx, id)
}

func (u *usecaseService) FindById(ctx context.Context, id string) (*domain.Service, error) {
	return u.ServiceRepository.FindById(ctx, id)
}

func (u *usecaseService) Count(ctx context.Context) (int, error) {
	return u.ServiceRepository.Count(ctx)
}
