package usecase

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/bcrypt"
	"github.com/elizandrodantas/machine-controller-v2/internal/util"
	"github.com/jackc/pgx/v5"
)

type userUsecase struct {
	UserRepository domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) domain.UserUsecase {
	return &userUsecase{r}
}

func (u *userUsecase) Create(c context.Context, userCreate *domain.UserCreate) error {
	if u.ExistUsername(c, userCreate.Username) {
		return domain.ErrUserAlreayExist
	}

	passwordHash, err := bcrypt.Hash([]byte(userCreate.Password))
	if err != nil {
		return err
	}

	return u.UserRepository.Create(c, &domain.User{
		Name:     userCreate.Name,
		Username: userCreate.Username,
		Password: passwordHash,
		Scope:    userCreate.Scope,
	})
}

func (u *userUsecase) FindByUsername(c context.Context, username string) (*domain.User, error) {
	return u.UserRepository.FindByUsername(c, username)
}

func (u *userUsecase) ExistUsername(c context.Context, username string) bool {
	_, err := u.FindByUsername(c, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false
		}
	}

	return true
}

func (u *userUsecase) FindById(c context.Context, id string) (*domain.User, error) {
	return u.UserRepository.FindById(c, id)
}

func (u *userUsecase) UpdatePassword(c context.Context, id, password string) error {
	if _, err := u.FindById(c, id); err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}

		return err
	}

	passwordHash, err := bcrypt.Hash([]byte(password))
	if err != nil {
		return err
	}

	return u.UserRepository.UpdatePassword(c, id, passwordHash)
}

func (u *userUsecase) List(c context.Context, page int) ([]domain.UserJson, error) {
	l, err := u.UserRepository.List(c, page)
	if err != nil {
		return nil, err
	}

	if len(l) == 0 {
		return []domain.UserJson{}, nil
	}

	var response = make([]domain.UserJson, 0, len(l))
	for _, item := range l {
		response = append(response, domain.UserJson{
			ID:        item.ID,
			Name:      item.Name,
			Username:  item.Username,
			Status:    item.Status,
			Scope:     item.Scope,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return response, nil
}

func (u *userUsecase) AddScopeWithUserId(c context.Context, id string, s string) error {
	d, err := u.FindById(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}

		return err
	}

	if util.ArrayFind(d.Scope, s) {
		return domain.ErrScopeAlreadyExist
	}

	return u.UserRepository.AddScope(c, id, s)
}

func (u *userUsecase) RemoveScopeWithUserId(c context.Context, id string, s string) error {
	d, err := u.FindById(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}

		return err
	}

	if !util.ArrayFind(d.Scope, s) {
		return nil
	}

	return u.UserRepository.RemoveScope(c, id, s)
}

func (u *userUsecase) Disable(c context.Context, id string) error {
	d, err := u.FindById(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}

		return err
	}

	if !d.Status {
		return nil
	}

	return u.UserRepository.ChangeStatus(c, id)
}

func (u *userUsecase) Enable(c context.Context, id string) error {
	d, err := u.FindById(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}

		return err
	}

	if d.Status {
		return nil
	}

	return u.UserRepository.ChangeStatus(c, id)
}

func (u *userUsecase) Count(c context.Context) (int, error) {
	return u.UserRepository.Count(c)
}

func (u *userUsecase) DeleteById(c context.Context, id string) error {
	if _, err := u.FindById(c, id); err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}
	}

	return u.UserRepository.Delete(c, id)
}

func (u *userUsecase) DeleteByUsername(c context.Context, username string) error {
	d, err := u.FindByUsername(c, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotExist
		}
	}

	return u.UserRepository.Delete(c, d.ID)
}
