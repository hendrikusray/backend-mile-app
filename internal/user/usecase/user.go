package usecase

import (
	"context"
	"mile-app-test/domain"
	"mile-app-test/utils"
	"strings"
	"time"
)

type userUseCase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(r domain.UserRepository, duration time.Duration) domain.UserUseCase {
	return &userUseCase{
		userRepo:       r,
		contextTimeout: duration,
	}
}

func (u *userUseCase) Login(ctx context.Context, in *domain.Login) (domain.Token, error) {
	if strings.TrimSpace(in.Username) == "" || strings.TrimSpace(in.Password) == "" {
		return domain.Token{}, utils.ErrValidation
	}

	usr, err := u.userRepo.GetUser(ctx, in.Username)
	if err != nil {
		return domain.Token{}, err
	}
	if usr == nil || usr.Pin != in.Password {
		return domain.Token{}, utils.ErrValidation
	}

	return domain.Token{AccessToken: "mock-token", RefreshToken: "mock-refresh"}, nil
}
