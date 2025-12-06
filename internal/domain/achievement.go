package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidAchievementData = errors.New("invalid achievement data")
)

type Achievement struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Points      int       `json:"points"`
	Type        string    `json:"type"`
	Condition   string    `json:"condition"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserAchievement struct {
	ID            uint      `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AchievementID uint      `json:"achievement_id"`
	UnlockedAt    time.Time `json:"unlocked_at"`
}

type AchievementRepository interface {
	Create(ctx context.Context, achievement *Achievement) error
	GetByID(ctx context.Context, id uint) (*Achievement, error)
	ListActive(ctx context.Context) ([]*Achievement, error)
	Update(ctx context.Context, achievement *Achievement) error

	// User achievements
	CreateUserAchievement(ctx context.Context, userAchievement *UserAchievement) error
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*UserAchievement, error)
	HasAchievement(ctx context.Context, userID uuid.UUID, achievementID uint) (bool, error)

	// Bulk operations for seeding
	CreateMultiple(ctx context.Context, achievements []*Achievement) error
}

type AchievementService interface {
	ListActive(ctx context.Context) ([]*Achievement, error)
	GrantIfEligible(ctx context.Context, userID uuid.UUID, event AchievementEvent) error
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*UserAchievement, error)
}

// AchievementEvent описывает внешнее событие (например, набор очков или количество отзывов),
// на основе которого сервис решает, пора ли выдавать достижение.
type AchievementEvent struct {
	Type  string
	Stats map[string]int
}
