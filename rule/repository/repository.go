package repository

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ruleRepository struct {
	Conn *pgxpool.Pool
}

func NewRuleRepository(Conn *pgxpool.Pool) domain.RuleRepository {
	return &ruleRepository{Conn}
}

func (r *ruleRepository) ActiveByMachineIdAndServiceId(ctx context.Context, machineId, serviceId string) (*domain.MachineRule, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	var machineRule domain.MachineRule
	err = pool.QueryRow(ctx,
		"SELECT * FROM machine_rules WHERE machine_id = $1 AND service_id = $2 AND "+
			"expire > EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::INTEGER ORDER BY created_at DESC LIMIT 1",
		machineId, serviceId).Scan(&machineRule.MachineId, &machineRule.ServiceId, &machineRule.Expire, &machineRule.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &machineRule, nil
}

func (r *ruleRepository) Create(ctx context.Context, machineId, serviceId string, expire int) error {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(ctx,
		"INSERT INTO machine_rules (machine_id, service_id, expire) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		machineId, serviceId, expire)
	if err != nil {
		return err
	}

	return nil
}

func (r *ruleRepository) List(ctx context.Context, page int) ([]domain.RuleJoinService, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	rows, err := pool.Query(ctx,
		"SELECT s.name as service_name, mr.* FROM machine_rules mr INNER JOIN services s ON s.id = mr.service_id "+
			"ORDER BY mr.created_at DESC LIMIT 10 OFFSET $1",
		page*10)

	if err != nil {
		return nil, err
	}

	machineRules, err := pgx.CollectRows(rows, pgx.RowToStructByPos[domain.RuleJoinService])
	if err != nil {
		return nil, err
	}

	return machineRules, nil
}

func (r *ruleRepository) History(ctx context.Context, machineId string, serviceId string, page int) ([]domain.MachineRule, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	var rows pgx.Rows
	if serviceId != "" {
		rows, err = pool.Query(ctx,
			"SELECT * FROM machine_rules WHERE machine_id = $1 AND service_id = $2 ORDER BY created_at DESC LIMIT 10 OFFSET $3",
			machineId, serviceId, page*10)
	} else {
		rows, err = pool.Query(ctx,
			"SELECT * FROM machine_rules WHERE machine_id = $1 ORDER BY created_at DESC LIMIT 10 OFFSET $2",
			machineId, page*10)
	}

	if err != nil {
		return nil, err
	}

	machineRules, err := pgx.CollectRows(rows, pgx.RowToStructByPos[domain.MachineRule])
	if err != nil {
		return nil, err
	}

	return machineRules, nil
}

func (r *ruleRepository) Count(ctx context.Context) (int, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer pool.Release()

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM machine_rules").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ruleRepository) CountByMachineId(ctx context.Context, machineId string, serviceId string) (int, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer pool.Release()

	var count int
	if serviceId != "" {
		err = pool.QueryRow(ctx,
			"SELECT COUNT(*) FROM machine_rules WHERE machine_id = $1 AND service_id = $2",
			machineId, serviceId).Scan(&count)
	} else {
		err = pool.QueryRow(ctx,
			"SELECT COUNT(*) FROM machine_rules WHERE machine_id = $1",
			machineId).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ruleRepository) Actives(ctx context.Context, machineId string) ([]domain.RuleJoinService, error) {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	rows, _ := pool.Query(ctx,
		"SELECT s.name as service_name, mr.* FROM machine_rules mr INNER JOIN services s ON s.id = mr.service_id "+
			"WHERE mr.machine_id = $1 AND mr.expire > EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::INTEGER ORDER BY mr.created_at DESC",
		machineId)
	rulesActives, err := pgx.CollectRows(rows, pgx.RowToStructByPos[domain.RuleJoinService])
	if err != nil {
		return nil, err
	}

	return rulesActives, nil
}

func (r ruleRepository) Invalidate(ctx context.Context, machineId, serviceId string) error {
	pool, err := r.Conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(ctx,
		"UPDATE machine_rules SET expire = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::INTEGER "+
			"WHERE machine_id = $1 AND service_id = $2 AND expire > EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::INTEGER",
		machineId, serviceId)
	if err != nil {
		return err
	}

	return nil
}
