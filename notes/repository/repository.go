package repository

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type notesRepository struct {
	Conn *pgxpool.Pool
}

func NewNotesRepository(Conn *pgxpool.Pool) domain.NotesRepository {
	return &notesRepository{Conn}
}

func (n *notesRepository) Create(c context.Context, text, machineId, userId string) error {
	pool, err := n.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "INSERT INTO notes (text, machine_id, user_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", text, machineId, userId)
	return err
}

func (n *notesRepository) ListByMachineId(c context.Context, machineId string, page int) ([]domain.NotesList, error) {
	pool, err := n.Conn.Acquire(c)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	rows, _ := pool.Query(c, "SELECT n.*, u.name AS user_name FROM notes n "+
		"INNER JOIN users u ON u.id = n.user_id WHERE n.machine_id = $1 "+
		"ORDER BY created_at, updated_at DESC LIMIT 10 OFFSET $2", machineId, 10*page)
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[domain.NotesList])
}

func (n *notesRepository) ChangeNote(c context.Context, id, text string) error {
	pool, err := n.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "UPDATE notes SET text = $1, updated_at = now() WHERE id = $2", text, id)
	return err
}

func (n *notesRepository) Delete(c context.Context, id string) error {
	pool, err := n.Conn.Acquire(c)
	if err != nil {
		return err
	}
	defer pool.Release()

	_, err = pool.Exec(c, "DELETE FROM notes WHERE id = $1", id)
	return err
}

func (n *notesRepository) FindById(c context.Context, id string) (domain.Notes, error) {
	pool, err := n.Conn.Acquire(c)
	if err != nil {
		return domain.Notes{}, err
	}
	defer pool.Release()

	var note domain.Notes
	err = pool.QueryRow(c, "SELECT id, text, machine_id, user_id, created_at, updated_at FROM notes WHERE id = $1 ORDER BY created_at, updated_at DESC", id).
		Scan(&note.ID, &note.Text, &note.MachineId, &note.UserId, &note.CreatedAt, &note.UpdatedAt)
	return note, err
}

func (n *notesRepository) CountByMachineId(c context.Context, machineId string) (int, error) {
	pool, err := n.Conn.Acquire(c)
	if err != nil {
		return 0, err
	}
	defer pool.Release()

	var count int
	err = pool.QueryRow(c, "SELECT COUNT(*) FROM notes n "+
		"INNER JOIN users u ON u.id = n.user_id WHERE n.machine_id = $1", machineId).Scan(&count)
	return count, err
}
