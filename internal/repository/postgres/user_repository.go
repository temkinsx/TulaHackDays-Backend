package postgres

import (
	"TulaHackDays-Backend/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	q Querier
}

func NewUserRepository(q Querier) domain.UserRepository {
	return &userRepository{q: q}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (email, username, password, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := u.q.Exec(ctx, q, user.Email, user.Username, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, username, password, 
		       first_name, last_name, avatar, 
		       points, level, is_active, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	var user domain.User
	err := u.q.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Avatar,
		&user.Points,
		&user.Level,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, err
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, username, password, 
		       first_name, last_name, avatar, 
		       points, level, is_active, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	var user domain.User
	err := u.q.QueryRow(ctx, q, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Avatar,
		&user.Points,
		&user.Level,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, err
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const q = `
		SELECT id, email, username, password, 
		       first_name, last_name, avatar, 
		       points, level, is_active, created_at, updated_at
		FROM users
		WHERE username = $1;
	`

	var user domain.User
	err := u.q.QueryRow(ctx, q, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Avatar,
		&user.Points,
		&user.Level,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user.Password = ""

	return &user, err
}

func (u *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	const q = `
  	UPDATE users
  	SET email      = $2,
  	    username   = $3,
  	    password   = $4,
  	    first_name = $5,
  	    last_name  = $6,
  	    avatar     = $7,
  	    points     = $8,
  	    level      = $9,
  	    is_active  = $10,
  	    updated_at = NOW()
  	WHERE id = $1
  	RETURNING id, email, username, first_name, last_name, avatar, points, level, is_active, created_at, updated_at;
  `

	var usr domain.User
	err := u.q.QueryRow(ctx, q, user.ID, user.Email, user.Username, user.Password, user.FirstName, user.LastName, user.Avatar, user.Points, user.Level, user.IsActive).Scan(
		&usr.ID,
		&usr.Email,
		&usr.Username,
		&usr.Password,
		&usr.FirstName,
		&usr.LastName,
		&usr.Avatar,
		&usr.Points,
		&usr.Level,
		&usr.IsActive,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	usr.Password = ""
	return &usr, err
}

func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `
		DELETE
		FROM users
		WHERE id = $1;
	`

	cmd, err := u.q.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (u *userRepository) GetLeaderboard(ctx context.Context, limit int) ([]*domain.User, error) {
	const q = `
		SELECT id, email, username, password, 
		       first_name, last_name, avatar, 
		       points, level, is_active, created_at, updated_at
		FROM users
		ORDER BY points DESC
		LIMIT $1;
	`

	rows, err := u.q.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (*domain.User, error) {
		usr := new(domain.User)
		if err := r.Scan(
			&usr.ID,
			&usr.Email,
			&usr.Username,
			&usr.Password,
			&usr.FirstName,
			&usr.LastName,
			&usr.Avatar,
			&usr.Points,
			&usr.Level,
			&usr.IsActive,
			&usr.CreatedAt,
			&usr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		usr.Password = ""
		return usr, nil
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepository) ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT *
			FROM users
			WHERE email = $1 
			   OR username = $2
		);
	`

	var exists bool
	err := u.q.QueryRow(ctx, q, email, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
