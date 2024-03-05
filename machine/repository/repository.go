package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type machineRepository struct {
	Conn *pgxpool.Pool
}

func NewMachineRepository(conn *pgxpool.Pool) domain.MachineRepository {
	return &machineRepository{Conn: conn}
}

func (m *machineRepository) Create(ctx context.Context, guid, name, os, query string) error {
	pool, err := m.Conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(ctx, "INSERT INTO machines (guid, name, os, query) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", guid, name, os, query)
	return err
}

func (m *machineRepository) FindByGuid(ctx context.Context, guid string) (*domain.Machine, error) {
	pool, err := m.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	machine := domain.Machine{}
	err = pool.QueryRow(ctx, "SELECT * FROM machines WHERE guid = $1 ORDER BY created_at DESC", guid).
		Scan(&machine.ID, &machine.Guid, &machine.Name, &machine.OS, &machine.Query, &machine.CreatedAt)
	return &machine, err
}

func (m *machineRepository) List(ctx context.Context, q domain.MachineListQuerys) ([]domain.Machine, error) {
	pool, err := m.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	var sql strings.Builder
	var params []interface{}

	sql.WriteString("SELECT * FROM machines WHERE 1=1")

	if q.MachineId != "" {
		sql.WriteString(" AND id = $" + strconv.Itoa(len(params)+1))
		params = append(params, q.MachineId)
	}

	if q.Query != "" {
		sql.WriteString(" AND query LIKE $" + strconv.Itoa(len(params)+1))
		params = append(params, "%"+q.Query+"%")
	}

	if q.OS != "" {
		sql.WriteString(" AND os = $" + strconv.Itoa(len(params)+1))
		params = append(params, q.OS)
	}

	sql.WriteString(" ORDER BY created_at DESC LIMIT 10 OFFSET $" + strconv.Itoa(len(params)+1))
	params = append(params, q.Page*10)

	rows, _ := pool.Query(ctx, sql.String(), params...)
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[domain.Machine])
}

func (m *machineRepository) FindById(ctx context.Context, id string) (*domain.Machine, error) {
	pool, err := m.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	machine := domain.Machine{}
	err = pool.QueryRow(ctx, "SELECT * FROM machines WHERE id = $1", id).
		Scan(&machine.ID, &machine.Guid, &machine.Name, &machine.OS, &machine.Query, &machine.CreatedAt)
	return &machine, err
}

func (m *machineRepository) Count(ctx context.Context) (int, error) {
	pool, err := m.Conn.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer pool.Release()

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM machines").Scan(&count)
	return count, err
}

func (m *machineRepository) UpdateName(ctx context.Context, id, name string) error {
	pool, err := m.Conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(ctx, "UPDATE machines SET name = $1 WHERE id = $2", name, id)
	return err
}
