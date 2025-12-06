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

type AchievementType string

const (
	AchievementTypeFirstReview    AchievementType = "first_review"
	AchievementTypeReviewMaster   AchievementType = "review_master"
	AchievementTypePlaceCreator   AchievementType = "place_creator"
	AchievementTypeCommenter      AchievementType = "commenter"
	AchievementTypePhotographer   AchievementType = "photographer"
	AchievementTypeHealthAdvocate AchievementType = "health_advocate"
)

// Predefined achievements
var DefaultAchievements = []*Achievement{
	{
		Name:        "Первый отзыв",
		Description: "Оставьте свой первый отзыв",
		Icon:        "⭐",
		Points:      10,
		Type:        string(AchievementTypeFirstReview),
		Condition:   `{"reviews_count": 1}`,
		IsActive:    true,
	},
	{
		Name:        "Мастер отзывов",
		Description: "Оставьте 50 отзывов",
		Icon:        "🏆",
		Points:      100,
		Type:        string(AchievementTypeReviewMaster),
		Condition:   `{"reviews_count": 50}`,
		IsActive:    true,
	},
	{
		Name:        "Создатель мест",
		Description: "Добавьте 10 новых объектов",
		Icon:        "📍",
		Points:      150,
		Type:        string(AchievementTypePlaceCreator),
		Condition:   `{"places_count": 10}`,
		IsActive:    true,
	},
	{
		Name:        "Комментатор",
		Description: "Оставьте 20 комментариев",
		Icon:        "💬",
		Points:      50,
		Type:        string(AchievementTypeCommenter),
		Condition:   `{"comments_count": 20}`,
		IsActive:    true,
	},
	{
		Name:        "Фотограф",
		Description: "Добавьте 25 фотографий",
		Icon:        "📸",
		Points:      75,
		Type:        string(AchievementTypePhotographer),
		Condition:   `{"photos_count": 25}`,
		IsActive:    true,
	},
	{
		Name:        "Защитник здоровья",
		Description: "Получите 500 очков",
		Icon:        "❤️",
		Points:      200,
		Type:        string(AchievementTypeHealthAdvocate),
		Condition:   `{"points": 500}`,
		IsActive:    true,
	},
}

type AchievementRepository interface {
	Create(ctx context.Context, achievement *Achievement) error
	GetByID(ctx context.Context, id uint) (*Achievement, error)
	GetAllActive(ctx context.Context) ([]*Achievement, error)
	Update(ctx context.Context, achievement *Achievement) error

	// User achievements
	CreateUserAchievement(ctx context.Context, userAchievement *UserAchievement) error
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*UserAchievement, error)
	HasAchievement(ctx context.Context, userID uuid.UUID, achievementID uint) (bool, error)

	// Bulk operations for seeding
	CreateMultiple(ctx context.Context, achievements []*Achievement) error
}

type AchievementService interface {
	GetActive(ctx context.Context) ([]*Achievement, error)
	GrantIfEligible(ctx context.Context, userID uuid.UUID, event AchievementEvent) error
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*UserAchievement, error)
}

type AchievementEvent struct {
	Type   AchievementType
	Points int
	Count  int
}
