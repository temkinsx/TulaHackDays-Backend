package dto

import (
	"TulaHackDays-Backend/internal/domain"
	"time"
)

type AchievementDTO struct {
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

func AchievementEntityToDTO(a *domain.Achievement) AchievementDTO {
	return AchievementDTO{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Icon:        a.Icon,
		Points:      a.Points,
		Type:        a.Type,
		Condition:   a.Condition,
		IsActive:    a.IsActive,
		CreatedAt:   a.CreatedAt,
	}
}

func AchievementsEntitiesToDTOs(achievements []*domain.Achievement) []AchievementDTO {
	result := make([]AchievementDTO, len(achievements))
	for i, achievement := range achievements {
		result[i] = AchievementEntityToDTO(achievement)
	}
	return result
}
