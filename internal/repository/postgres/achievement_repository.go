package postgres

import (
	"TulaHackDays-Backend/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type achievementRepository struct {
	q Querier
}

func NewAchievementRepository(q Querier) domain.AchievementRepository {
	return &achievementRepository{q: q}
}

func (r *achievementRepository) Create(ctx context.Context, achievement *domain.Achievement) error {
	const query = `
		INSERT INTO achievements (name, description, icon, points, type, condition, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at;
	`

	return r.q.QueryRow(ctx, query,
		achievement.Name,
		achievement.Description,
		achievement.Icon,
		achievement.Points,
		achievement.Type,
		achievement.Condition,
		achievement.IsActive,
	).Scan(&achievement.ID, &achievement.CreatedAt)
}

func (r *achievementRepository) GetByID(ctx context.Context, id uint) (*domain.Achievement, error) {
	const query = `
		SELECT id, name, description, icon, points, type, condition, is_active, created_at
		FROM achievements
		WHERE id = $1;
	`

	var achievement domain.Achievement
	err := r.q.QueryRow(ctx, query, id).Scan(
		&achievement.ID,
		&achievement.Name,
		&achievement.Description,
		&achievement.Icon,
		&achievement.Points,
		&achievement.Type,
		&achievement.Condition,
		&achievement.IsActive,
		&achievement.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAchievementNotFound
		}
		return nil, err
	}

	return &achievement, nil
}

func (r *achievementRepository) ListActive(ctx context.Context) ([]*domain.Achievement, error) {
	const query = `
		SELECT id, name, description, icon, points, type, condition, is_active, created_at
		FROM achievements
		WHERE is_active = true
		ORDER BY id;
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.Achievement, error) {
		achievement := new(domain.Achievement)
		if err := row.Scan(
			&achievement.ID,
			&achievement.Name,
			&achievement.Description,
			&achievement.Icon,
			&achievement.Points,
			&achievement.Type,
			&achievement.Condition,
			&achievement.IsActive,
			&achievement.CreatedAt,
		); err != nil {
			return nil, err
		}
		return achievement, nil
	})
	if err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *achievementRepository) Update(ctx context.Context, achievement *domain.Achievement) error {
	const query = `
		UPDATE achievements
		SET name = $2,
		    description = $3,
		    icon = $4,
		    points = $5,
		    type = $6,
		    condition = $7,
		    is_active = $8
		WHERE id = $1;
	`

	cmd, err := r.q.Exec(ctx, query,
		achievement.ID,
		achievement.Name,
		achievement.Description,
		achievement.Icon,
		achievement.Points,
		achievement.Type,
		achievement.Condition,
		achievement.IsActive,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrAchievementNotFound
	}

	return nil
}

func (r *achievementRepository) CreateUserAchievement(ctx context.Context, userAchievement *domain.UserAchievement) error {
	const query = `
		INSERT INTO user_achievements (user_id, achievement_id)
		VALUES ($1, $2)
		RETURNING id, unlocked_at;
	`

	return r.q.QueryRow(ctx, query, userAchievement.UserID, userAchievement.AchievementID).
		Scan(&userAchievement.ID, &userAchievement.UnlockedAt)
}

func (r *achievementRepository) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*domain.UserAchievement, error) {
	const query = `
		SELECT id, user_id, achievement_id, unlocked_at
		FROM user_achievements
		WHERE user_id = $1
		ORDER BY unlocked_at DESC;
	`

	rows, err := r.q.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userAchievements, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.UserAchievement, error) {
		item := new(domain.UserAchievement)
		if err := row.Scan(
			&item.ID,
			&item.UserID,
			&item.AchievementID,
			&item.UnlockedAt,
		); err != nil {
			return nil, err
		}
		return item, nil
	})
	if err != nil {
		return nil, err
	}

	return userAchievements, nil
}

func (r *achievementRepository) HasAchievement(ctx context.Context, userID uuid.UUID, achievementID uint) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM user_achievements
			WHERE user_id = $1 AND achievement_id = $2
		);
	`

	var exists bool
	if err := r.q.QueryRow(ctx, query, userID, achievementID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *achievementRepository) CreateMultiple(ctx context.Context, achievements []*domain.Achievement) error {
	const query = `
		INSERT INTO achievements (name, description, icon, points, type, condition, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (name) DO NOTHING;
	`

	for _, achievement := range achievements {
		if achievement == nil {
			continue
		}
		if _, err := r.q.Exec(ctx, query,
			achievement.Name,
			achievement.Description,
			achievement.Icon,
			achievement.Points,
			achievement.Type,
			achievement.Condition,
			achievement.IsActive,
		); err != nil {
			return err
		}
	}

	return nil
}
