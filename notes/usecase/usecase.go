package usecase

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/jackc/pgx/v5"
)

type notesUsecase struct {
	notesRepo domain.NotesRepository
}

func NewNotesRepository(notesRepo domain.NotesRepository) domain.NotesUsecase {
	return &notesUsecase{notesRepo}
}

func (n *notesUsecase) Create(c context.Context, text, machineId, userId string) error {
	return n.notesRepo.Create(c, text, machineId, userId)
}

func (n *notesUsecase) ListByMachineId(c context.Context, machineId string, last int) ([]domain.NotesListJson, error) {
	list, err := n.notesRepo.ListByMachineId(c, machineId, last)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []domain.NotesListJson{}, nil
	}

	var response = make([]domain.NotesListJson, 0, len(list))
	for _, note := range list {
		response = append(response, domain.NotesListJson(note))
	}

	return response, nil
}

func (n *notesUsecase) ChangeNote(c context.Context, id, userId, text string) error {
	i, err := n.notesRepo.FindById(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrNoteNotFound
		}

		return err
	}

	if i.UserId != userId {
		return domain.ErrOnlyCreatorCanChange
	}

	if text == i.Text {
		return nil
	}

	return n.notesRepo.ChangeNote(c, id, text)
}

func (n *notesUsecase) Delete(c context.Context, userId, noteId string) error {
	nt, err := n.notesRepo.FindById(c, noteId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}

		return err
	}

	if nt.UserId != userId {
		return domain.ErrOnlyCreatorCanDelete
	}

	return n.notesRepo.Delete(c, noteId)
}

func (n *notesUsecase) CountByMachineId(c context.Context, machineId string) (int, error) {
	return n.notesRepo.CountByMachineId(c, machineId)
}
