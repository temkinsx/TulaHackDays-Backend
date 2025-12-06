package services

import (
	"context"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type userService struct {
	userRepo  domain.UserRepository
	txManager domain.TxManager
}

func NewUserService(userRepo domain.UserRepository, txManager domain.TxManager) domain.UserService {
	return &userService{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (uc *userService) Register(ctx context.Context, input *domain.RegisterUserInput) (*domain.User, error) {
	panic("not implemented")
}

func (uc *userService) Login(ctx context.Context, input *domain.LoginUserInput) (*domain.LoginResult, error) {
	panic("not implemented")
}

func (uc *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	panic("not implemented")
}

func (uc *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, input *domain.UpdateProfileInput) (*domain.User, error) {
	panic("not implemented")
}

func (uc *userService) ChangePassword(ctx context.Context, userID uuid.UUID, input *domain.ChangePasswordInput) error {
	panic("not implemented")
}

func (uc *userService) GetLeaderboard(ctx context.Context, limit int) ([]*domain.User, error) {
	panic("not implemented")
}

func (uc *userService) GetAchievements(ctx context.Context, userID uuid.UUID) ([]*domain.UserAchievement, error) {
	panic("not implemented")
}

func (uc *userService) AddPoints(ctx context.Context, userID uuid.UUID, points int) error {
	panic("not implemented")
}
