package usecase

import (
	"context"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/bcrypt"
	"github.com/jackc/pgx/v5"
)

type authUsecase struct {
	UserRepository domain.UserRepository
}

func NewAuthUsecase(r domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{r}
}

func (a *authUsecase) Autenticate(c context.Context, username, password string) (*domain.User, error) {
	userInfo, err := a.UserRepository.FindByUsername(c, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUsernameOrPasswordInvalid
		}

		return nil, err
	}

	if err := bcrypt.Compare([]byte(userInfo.Password), []byte(password)); err != nil {
		return nil, domain.ErrUsernameOrPasswordInvalid
	}

	if !userInfo.Status {
		return nil, domain.ErrUserIsNotActive
	}

	return userInfo, nil
}
