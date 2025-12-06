package postgres

import (
	"TulaHackDays-Backend/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type placeRepository struct {
	q Querier
}

func NewPlaceRepository(q Querier) domain.PlaceRepository {
	return &placeRepository{q: q}
}

func (p *placeRepository) Create(ctx context.Context, place *domain.Place) error {
	const q = `
		INSERT INTO places (id, name, description, type, address, latitude, longitude, phone, website, images, average_rating, review_count, is_active, created_by_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
	`

	_, err := p.q.Exec(ctx, q, place.ID, place.Name, place.Description, place.Type, place.Address,
		place.Latitude, place.Longitude, place.Phone, place.Website, place.Images,
		place.AverageRating, place.ReviewCount, place.IsActive, place.CreatedByID)
	if err != nil {
		return err
	}

	return nil
}

func (p *placeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Place, error) {
	const q = `
		SELECT id, name, description, type, address, latitude, longitude, phone, website, images,
		       average_rating, review_count, is_active, created_by_id, created_at, updated_at
		FROM places
		WHERE id = $1;
	`

	var place domain.Place
	err := p.q.QueryRow(ctx, q, id).Scan(
		&place.ID,
		&place.Name,
		&place.Description,
		&place.Type,
		&place.Address,
		&place.Latitude,
		&place.Longitude,
		&place.Phone,
		&place.Website,
		&place.Images,
		&place.AverageRating,
		&place.ReviewCount,
		&place.IsActive,
		&place.CreatedByID,
		&place.CreatedAt,
		&place.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlaceNotFound
		}
		return nil, err
	}

	return &place, nil
}

func (p *placeRepository) Update(ctx context.Context, place *domain.Place) error {
	const q = `
		UPDATE places
		SET name = $2, description = $3, type = $4, address = $5, latitude = $6, longitude = $7,
		    phone = $8, website = $9, images = $10, average_rating = $11, review_count = $12,
		    is_active = $13, updated_at = NOW()
		WHERE id = $1;
	`

	cmd, err := p.q.Exec(ctx, q, place.ID, place.Name, place.Description, place.Type, place.Address,
		place.Latitude, place.Longitude, place.Phone, place.Website, place.Images,
		place.AverageRating, place.ReviewCount, place.IsActive)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrPlaceNotFound
	}

	return nil
}

func (p *placeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `
		DELETE FROM places
		WHERE id = $1;
	`

	cmd, err := p.q.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrPlaceNotFound
	}

	return nil
}

func (p *placeRepository) Search(ctx context.Context, query string, placeType *domain.PlaceType, limit, offset int) ([]*domain.Place, int, error) {
	var args []interface{}
	argCount := 0

	baseQuery := `
		SELECT id, name, description, type, address, latitude, longitude, phone, website, images,
		       average_rating, review_count, is_active, created_by_id, created_at, updated_at
		FROM places
		WHERE is_active = true
	`

	countQuery := `
		SELECT COUNT(*)
		FROM places
		WHERE is_active = true
	`

	if query != "" {
		argCount++
		baseQuery += ` AND (name ILIKE $` + string(rune('0'+argCount)) + ` OR description ILIKE $` + string(rune('0'+argCount)) + `)`
		countQuery += ` AND (name ILIKE $` + string(rune('0'+argCount)) + ` OR description ILIKE $` + string(rune('0'+argCount)) + `)`
		args = append(args, "%"+query+"%")
	}

	if placeType != nil {
		argCount++
		baseQuery += ` AND type = $` + string(rune('0'+argCount))
		countQuery += ` AND type = $` + string(rune('0'+argCount))
		args = append(args, *placeType)
	}

	baseQuery += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argCount+1)) + ` OFFSET $` + string(rune('0'+argCount+2))
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := p.q.QueryRow(ctx, countQuery, args[:argCount]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get places
	rows, err := p.q.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	places, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (*domain.Place, error) {
		place := new(domain.Place)
		if err := r.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&place.Type,
			&place.Address,
			&place.Latitude,
			&place.Longitude,
			&place.Phone,
			&place.Website,
			&place.Images,
			&place.AverageRating,
			&place.ReviewCount,
			&place.IsActive,
			&place.CreatedByID,
			&place.CreatedAt,
			&place.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return place, nil
	})
	if err != nil {
		return nil, 0, err
	}

	return places, total, nil
}

func (p *placeRepository) GetNearby(ctx context.Context, coords *domain.Coordinates, radiusKm float64, limit int) ([]*domain.Place, error) {
	const q = `
		SELECT id, name, description, type, address, latitude, longitude, phone, website, images,
		       average_rating, review_count, is_active, created_by_id, created_at, updated_at, distance
		FROM (
			SELECT id, name, description, type, address, latitude, longitude, phone, website, images,
			       average_rating, review_count, is_active, created_by_id, created_at, updated_at,
			       (6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) AS distance
			FROM places
			WHERE is_active = true
		) AS places_with_distance
		WHERE distance < $3
		ORDER BY distance
		LIMIT $4;
	`

	rows, err := p.q.Query(ctx, q, coords.Latitude, coords.Longitude, radiusKm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	places, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (*domain.Place, error) {
		place := new(domain.Place)
		var distance float64
		if err := r.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&place.Type,
			&place.Address,
			&place.Latitude,
			&place.Longitude,
			&place.Phone,
			&place.Website,
			&place.Images,
			&place.AverageRating,
			&place.ReviewCount,
			&place.IsActive,
			&place.CreatedByID,
			&place.CreatedAt,
			&place.UpdatedAt,
			&distance,
		); err != nil {
			return nil, err
		}
		return place, nil
	})
	if err != nil {
		return nil, err
	}

	return places, nil
}

func (p *placeRepository) GetByTypeStats(ctx context.Context) (map[domain.PlaceType]int, error) {
	const q = `
		SELECT type, COUNT(*) as count
		FROM places
		WHERE is_active = true
		GROUP BY type;
	`

	rows, err := p.q.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[domain.PlaceType]int)
	for rows.Next() {
		var placeType domain.PlaceType
		var count int
		if err := rows.Scan(&placeType, &count); err != nil {
			return nil, err
		}
		stats[placeType] = count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (p *placeRepository) GetTopRated(ctx context.Context, limit int) ([]*domain.Place, error) {
	const q = `
		SELECT id, name, description, type, address, latitude, longitude, phone, website, images,
		       average_rating, review_count, is_active, created_by_id, created_at, updated_at
		FROM places
		WHERE is_active = true AND review_count > 0
		ORDER BY average_rating DESC, review_count DESC
		LIMIT $1;
	`

	rows, err := p.q.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	places, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (*domain.Place, error) {
		place := new(domain.Place)
		if err := r.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&place.Type,
			&place.Address,
			&place.Latitude,
			&place.Longitude,
			&place.Phone,
			&place.Website,
			&place.Images,
			&place.AverageRating,
			&place.ReviewCount,
			&place.IsActive,
			&place.CreatedByID,
			&place.CreatedAt,
			&place.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return place, nil
	})
	if err != nil {
		return nil, err
	}

	return places, nil
}
