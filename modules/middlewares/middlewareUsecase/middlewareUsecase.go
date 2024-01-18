package middlewareUsecase

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/middlewares"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/middlewares/middlewareRepository"
)

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecase struct {
	middlewareRepository middlewareRepository.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewareRepository middlewareRepository.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewareRepository: middlewareRepository,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecase) FindRole() ([]*middlewares.Role, error) {
	role, err := u.middlewareRepository.FindRole()
	if err != nil {
		return nil, err
	}
	return role, nil
}
