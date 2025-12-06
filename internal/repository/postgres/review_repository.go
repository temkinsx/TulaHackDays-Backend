package postgres

import (
	"TulaHackDays-Backend/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type reviewRepository struct {
	q Querier
}

func NewReviewRepository(q Querier) domain.ReviewRepository {
	return &reviewRepository{q: q}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) error {
	const q = `
		INSERT INTO reviews (place_id, user_id, title, content, images, rating, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	err := r.q.QueryRow(ctx, q, review.PlaceID, review.UserID, review.Title, review.Content,
		review.Images, review.Rating, review.IsActive).Scan(&review.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *reviewRepository) GetByID(ctx context.Context, id uint) (*domain.Review, error) {
	const q = `
		SELECT id, place_id, user_id, title, content, images, rating, is_active, created_at, updated_at
		FROM reviews
		WHERE id = $1;
	`

	var review domain.Review
	err := r.q.QueryRow(ctx, q, id).Scan(
		&review.ID,
		&review.PlaceID,
		&review.UserID,
		&review.Title,
		&review.Content,
		&review.Images,
		&review.Rating,
		&review.IsActive,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReviewNotFound
		}
		return nil, err
	}

	return &review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *domain.Review) error {
	const q = `
		UPDATE reviews
		SET title = $2, content = $3, images = $4, rating = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1;
	`

	cmd, err := r.q.Exec(ctx, q, review.ID, review.Title, review.Content, review.Images,
		review.Rating, review.IsActive)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepository) Delete(ctx context.Context, id uint) error {
	const q = `
		DELETE FROM reviews
		WHERE id = $1;
	`

	cmd, err := r.q.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepository) GetByPlaceID(ctx context.Context, placeID uuid.UUID, limit, offset int) ([]*domain.Review, int, error) {
	baseQuery := `
		SELECT id, place_id, user_id, title, content, images, rating, is_active, created_at, updated_at
		FROM reviews
		WHERE place_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	countQuery := `
		SELECT COUNT(*)
		FROM reviews
		WHERE place_id = $1 AND is_active = true;
	`

	var total int
	err := r.q.QueryRow(ctx, countQuery, placeID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.Query(ctx, baseQuery, placeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.Review, error) {
		review := new(domain.Review)
		if err := row.Scan(
			&review.ID,
			&review.PlaceID,
			&review.UserID,
			&review.Title,
			&review.Content,
			&review.Images,
			&review.Rating,
			&review.IsActive,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return review, nil
	})
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Review, int, error) {
	baseQuery := `
		SELECT id, place_id, user_id, title, content, images, rating, is_active, created_at, updated_at
		FROM reviews
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	countQuery := `
		SELECT COUNT(*)
		FROM reviews
		WHERE user_id = $1 AND is_active = true;
	`

	var total int
	err := r.q.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.Query(ctx, baseQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.Review, error) {
		review := new(domain.Review)
		if err := row.Scan(
			&review.ID,
			&review.PlaceID,
			&review.UserID,
			&review.Title,
			&review.Content,
			&review.Images,
			&review.Rating,
			&review.IsActive,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return review, nil
	})
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *reviewRepository) ExistsByUserAndPlace(ctx context.Context, userID, placeID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM reviews
			WHERE user_id = $1 AND place_id = $2 AND is_active = true
		);
	`

	var exists bool
	err := r.q.QueryRow(ctx, q, userID, placeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *reviewRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	const q = `
		INSERT INTO comments (review_id, user_id, content, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at;
	`

	err := r.q.QueryRow(ctx, q, comment.ReviewID, comment.UserID, comment.Content, comment.IsActive).Scan(
		&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *reviewRepository) GetCommentsByReviewID(ctx context.Context, reviewID uint, limit, offset int) ([]*domain.Comment, int, error) {
	baseQuery := `
		SELECT id, review_id, user_id, content, is_active, created_at
		FROM comments
		WHERE review_id = $1 AND is_active = true
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3;
	`

	countQuery := `
		SELECT COUNT(*)
		FROM comments
		WHERE review_id = $1 AND is_active = true;
	`

	var total int
	err := r.q.QueryRow(ctx, countQuery, reviewID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.Query(ctx, baseQuery, reviewID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	comments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.Comment, error) {
		comment := new(domain.Comment)
		if err := row.Scan(
			&comment.ID,
			&comment.ReviewID,
			&comment.UserID,
			&comment.Content,
			&comment.IsActive,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		return comment, nil
	})
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *reviewRepository) CreatePlaceRating(ctx context.Context, rating *domain.PlaceRating) error {
	const q = `
		INSERT INTO place_ratings (place_id, user_id, cleanliness, service, quality, value, overall)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (place_id, user_id) DO UPDATE SET
			cleanliness = EXCLUDED.cleanliness,
			service = EXCLUDED.service,
			quality = EXCLUDED.quality,
			value = EXCLUDED.value,
			overall = EXCLUDED.overall,
			created_at = NOW()
		RETURNING id, created_at;
	`

	err := r.q.QueryRow(ctx, q, rating.PlaceID, rating.UserID, rating.Cleanliness,
		rating.Service, rating.Quality, rating.Value, rating.Overall).Scan(
		&rating.ID, &rating.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *reviewRepository) GetPlaceRatings(ctx context.Context, placeID uuid.UUID) ([]*domain.PlaceRating, error) {
	const q = `
		SELECT id, place_id, user_id, cleanliness, service, quality, value, overall, created_at
		FROM place_ratings
		WHERE place_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.q.Query(ctx, q, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.PlaceRating, error) {
		rating := new(domain.PlaceRating)
		if err := row.Scan(
			&rating.ID,
			&rating.PlaceID,
			&rating.UserID,
			&rating.Cleanliness,
			&rating.Service,
			&rating.Quality,
			&rating.Value,
			&rating.Overall,
			&rating.CreatedAt,
		); err != nil {
			return nil, err
		}
		return rating, nil
	})
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (r *reviewRepository) UpdatePlaceRatingStats(ctx context.Context, placeID uuid.UUID) error {
	const q = `
		UPDATE places
		SET average_rating = (
			SELECT COALESCE(AVG(overall), 0)
			FROM place_ratings
			WHERE place_id = $1
		),
		review_count = (
			SELECT COUNT(*)
			FROM place_ratings
			WHERE place_id = $1
		),
		updated_at = NOW()
		WHERE id = $1;
	`

	_, err := r.q.Exec(ctx, q, placeID)
	return err
}
