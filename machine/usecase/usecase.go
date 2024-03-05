package usecase

import (
	"context"
	"strings"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/util"
	"github.com/jackc/pgx/v5"
)

type machineUsecase struct {
	MachineRepository domain.MachineRepository
}

func NewMachineUsecase(m domain.MachineRepository) domain.MachineUsecase {
	return &machineUsecase{MachineRepository: m}
}

func (m *machineUsecase) Create(ctx context.Context, md *domain.MachineData) error {
	query := util.CreateQueryStr([]string{md.Guid, md.Name, md.OS, md.ServiceId})
	return m.MachineRepository.Create(ctx, md.Guid, md.Name, md.OS, query)
}

func (m *machineUsecase) FindByGuid(ctx context.Context, guid string) (*domain.Machine, error) {
	return m.MachineRepository.FindByGuid(ctx, guid)
}

func (m *machineUsecase) List(ctx context.Context, q domain.MachineListQuerys) ([]domain.MachineJson, error) {
	list, err := m.MachineRepository.List(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []domain.MachineJson{}, nil
	}

	var response = make([]domain.MachineJson, 0, len(list))
	for _, machine := range list {
		response = append(response, domain.MachineJson(machine))
	}

	return response, nil
}

func (m *machineUsecase) Detail(ctx context.Context, id string) (*domain.MachineJson, error) {
	dt, err := m.MachineRepository.FindById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrMachineNotFound
		}

		return nil, err
	}

	output := domain.MachineJson(*dt)
	return &output, nil
}

func (m *machineUsecase) Count(ctx context.Context) (int, error) {
	return m.MachineRepository.Count(ctx)
}

func (m *machineUsecase) UpdateName(ctx context.Context, id, name string) error {
	dd, err := m.MachineRepository.FindById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrMachineNotFound
		}

		return err
	}

	if strings.Split(dd.Name, "|")[0] == name {
		return nil
	}

	name = util.CleanSpecialCharacters(name)

	// ALWAYS KEEP THE ORIGINAL NAME OF THE MACHINE
	if !strings.Contains(dd.Name, "|") {
		name = name + "|" + dd.Name
	} else {
		parts := strings.Split(dd.Name, "|")
		name = name + "|" + parts[len(parts)-1]
	}

	return m.MachineRepository.UpdateName(ctx, id, name)
}
