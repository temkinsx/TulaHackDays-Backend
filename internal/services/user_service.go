package services

import (
	"TulaHackDays-Backend/internal/pkg/auth"
	"context"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type userService struct {
	userRepository domain.UserRepository
	txManager      domain.TxManager
	authService    *auth.AuthService
}

func NewUserService(userRepository domain.UserRepository, txManager domain.TxManager, authService *auth.AuthService) domain.UserService {
	return &userService{
		userRepository: userRepository,
		txManager:      txManager,
		authService:    authService,
	}
}

func (u *userService) Register(ctx context.Context, input *domain.RegisterUserInput) (*domain.LoginResult, error) {
	user := &domain.User{
		Email:     input.Email,
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashed, err := u.authService.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed

	err = u.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		exists, err := repos.User.ExistsByEmailOrUsername(ctx, user.Email, user.Username)
		if err != nil {
			return err
		}

		if exists {
			return domain.ErrUserAlreadyExists
		}

		user, err = repos.User.Create(ctx, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	token, err := u.authService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResult{
		Token: token,
		User:  user,
	}, nil
}

func (u *userService) Login(ctx context.Context, input *domain.LoginUserInput) (*domain.LoginResult, error) {
	user, err := u.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	err = u.authService.CheckPassword(user.Password, input.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := u.authService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &domain.LoginResult{
		Token: token,
		User:  user,
	}, nil
}

func (u *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return u.userRepository.GetByID(ctx, userID)
}

func (u *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, input *domain.UpdateProfileInput) (*domain.User, error) {
	user := &domain.User{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Avatar:    input.Avatar,
	}

	user, err := u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) ChangePassword(ctx context.Context, userID uuid.UUID, input *domain.ChangePasswordInput) error {
	return u.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		user, err := repos.User.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		err = u.authService.CheckPassword(user.Password, input.OldPassword)
		if err != nil {
			return err
		}

		newHash, err := u.authService.HashPassword(input.NewPassword)
		if err != nil {
			return err
		}

		_, err = repos.User.Update(ctx, &domain.User{
			ID:       userID,
			Password: newHash,
		})
		return err
	})
}

func (u *userService) GetLeaderboard(ctx context.Context, limit int) ([]*domain.User, error) {
	return u.userRepository.GetLeaderboard(ctx, limit)
}

func (u *userService) AddPoints(ctx context.Context, userID uuid.UUID, points int) (int, error) {
	return u.userRepository.AddPoints(ctx, userID, points)
}
