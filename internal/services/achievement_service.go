package services

import (
	"TulaHackDays-Backend/internal/domain"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type achService struct {
	achRepository domain.AchievementRepository
	txManager     domain.TxManager
}

func NewAchievementService(achievementRepository domain.AchievementRepository, txManager domain.TxManager) domain.AchievementService {
	return &achService{achRepository: achievementRepository, txManager: txManager}
}

func (a *achService) ListActive(ctx context.Context) ([]*domain.Achievement, error) {
	return a.achRepository.ListActive(ctx)
}

func (a *achService) GrantIfEligible(ctx context.Context, userID uuid.UUID, event domain.AchievementEvent) error {
	if event.Type == "" {
		return nil
	}

	return a.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Achievement = a.achRepository

		achievements, err := repos.Achievement.ListActive(ctx)
		if err != nil {
			return err
		}

		for _, achievement := range achievements {
			if achievement == nil || achievement.Type != event.Type {
				continue
			}

			var condition map[string]int
			if achievement.Condition != "" {
				if err := json.Unmarshal([]byte(achievement.Condition), &condition); err != nil {
					return err
				}
			}

			if !meetsAchievementCondition(condition, event.Stats) {
				continue
			}

			hasAchievement, err := repos.Achievement.HasAchievement(ctx, userID, achievement.ID)
			if err != nil {
				return err
			}
			if hasAchievement {
				continue
			}

			if err := repos.Achievement.CreateUserAchievement(ctx, &domain.UserAchievement{
				UserID:        userID,
				AchievementID: achievement.ID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *achService) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*domain.UserAchievement, error) {
	return a.achRepository.GetUserAchievements(ctx, userID)
}

func meetsAchievementCondition(condition map[string]int, stats map[string]int) bool {
	if len(condition) == 0 {
		return true
	}
	if len(stats) == 0 {
		return false
	}

	for key, required := range condition {
		if stats[key] < required {
			return false
		}
	}

	return true
}
