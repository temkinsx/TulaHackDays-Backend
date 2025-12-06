package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never return password in JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    []byte    `json:"avatar"`
	Points    int       `json:"points"`
	Level     int       `json:"level"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Username == "" {
		return ErrInvalidUsername
	}
	if u.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	AddPoints(ctx context.Context, userID uuid.UUID, points int) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetLeaderboard(ctx context.Context, limit int) ([]*User, error)
	ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error)
}

type RegisterUserInput struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
}

type LoginUserInput struct {
	Email    string
	Password string
}

type UpdateProfileInput struct {
	FirstName string
	LastName  string
	Avatar    []byte
}

type ChangePasswordInput struct {
	OldPassword string
	NewPassword string
}

type LoginResult struct {
	Token string
	User  *User
}

type UserService interface {
	Register(ctx context.Context, input *RegisterUserInput) (*LoginResult, error)
	Login(ctx context.Context, input *LoginUserInput) (*LoginResult, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input *UpdateProfileInput) (*User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, input *ChangePasswordInput) error
	GetLeaderboard(ctx context.Context, limit int) ([]*User, error)
	GetAchievements(ctx context.Context, userID uuid.UUID) ([]*UserAchievement, error)
	AddPoints(ctx context.Context, userID uuid.UUID, points int) (int, error)
}
