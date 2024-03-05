package repository

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(c *pgxpool.Pool) domain.UserRepository {
	return &userRepository{c}
}

func (u *userRepository) Create(c context.Context, user *domain.User) error {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	if user.Scope == nil {
		user.Scope = make([]string, 0)
	}

	_, err = pool.Exec(c, "INSERT INTO users (name, username, password, scope) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		user.Name, user.Username, user.Password, user.Scope)

	return err
}

func (u *userRepository) FindByUsername(c context.Context, username string) (*domain.User, error) {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	user := domain.User{}
	err = pool.QueryRow(c, "SELECT * FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.Name, &user.Username, &user.Password, &user.Status, &user.Scope, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (u *userRepository) Count(c context.Context) (int, error) {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return 0, err
	}
	defer pool.Release()

	var count int
	err = pool.QueryRow(c, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (u *userRepository) FindById(c context.Context, id string) (*domain.User, error) {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	user := domain.User{}
	err = pool.QueryRow(c, "SELECT * FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Status, &user.Scope, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (u *userRepository) UpdatePassword(c context.Context, id, password string) error {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "UPDATE users SET password = $1 WHERE id = $2", password, id)
	return err
}

func (u *userRepository) List(c context.Context, page int) (users []domain.User, err error) {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return
	}
	defer pool.Release()

	rows, err := pool.Query(c, "SELECT * FROM users ORDER BY created_at DESC LIMIT 10 OFFSET $1", page*10)
	if err != nil {
		return
	}
	defer rows.Close()

	users, err = pgx.CollectRows(rows, pgx.RowToStructByPos[domain.User])
	return
}

func (u *userRepository) AddScope(c context.Context, id, s string) error {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "UPDATE users SET scope = array_append(scope, $1), updated_at = now() WHERE id = $2", s, id)
	return err
}

func (u *userRepository) RemoveScope(c context.Context, id string, s string) error {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "UPDATE users SET scope = array_remove(scope, $1), updated_at = now() WHERE id = $2", s, id)
	return err
}

func (u *userRepository) ChangeStatus(c context.Context, id string) error {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "UPDATE users SET status = NOT status, updated_at = now() WHERE id = $1", id)
	return err
}

func (u *userRepository) Delete(c context.Context, id string) error {
	pool, err := u.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "DELETE FROM users WHERE id = $1", id)
	return err
}
