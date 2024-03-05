package repository

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceRepository struct {
	Conn *pgxpool.Pool
}

func NewServiceRepository(Conn *pgxpool.Pool) domain.ServiceRepository {
	return &serviceRepository{Conn}
}

func (r *serviceRepository) Create(ctx context.Context, name, description string) error {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(ctx, "INSERT INTO services (name, description) VALUES ($1, $2) ON CONFLICT DO NOTHING", name, description)
	return err
}

func (r *serviceRepository) FindByName(ctx context.Context, name string) (*domain.Service, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	var service domain.Service
	err = pool.QueryRow(ctx, "SELECT * FROM services WHERE name = $1 ORDER BY created_at DESC", name).Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepository) List(ctx context.Context, page int) ([]domain.Service, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	rows, err := pool.Query(ctx, "SELECT * FROM services ORDER BY created_at DESC LIMIT 10 OFFSET $1", page)
	if err != nil {
		return nil, err
	}

	var services []domain.Service
	for rows.Next() {
		var service domain.Service
		err = rows.Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}

func (r *serviceRepository) Remove(ctx context.Context, id string) error {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(ctx, "DELETE FROM services WHERE id = $1", id)
	return err
}

func (r *serviceRepository) FindById(ctx context.Context, id string) (*domain.Service, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	var service domain.Service
	err = pool.QueryRow(ctx, "SELECT * FROM services WHERE id = $1", id).Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepository) Count(ctx context.Context) (int, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer pool.Release()

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM services").Scan(&count)
	return count, err
}
