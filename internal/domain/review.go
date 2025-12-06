package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Additional review errors
var (
	ErrInvalidReviewContent  = errors.New("invalid review content")
	ErrInvalidCommentContent = errors.New("invalid comment content")
)

type Review struct {
	ID       uint      `json:"id"`
	PlaceID  uuid.UUID `json:"place_id"`
	UserID   uuid.UUID `json:"user_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Images   []string  `json:"images"`
	Rating   float64   `json:"rating"`
	IsActive bool      `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlaceRating struct {
	ID          uint      `json:"id"`
	PlaceID     uuid.UUID `json:"place_id"`
	UserID      uuid.UUID `json:"user_id"`
	Cleanliness float64   `json:"cleanliness"`
	Service     float64   `json:"service"`
	Quality     float64   `json:"quality"`
	Value       float64   `json:"value"`
	Overall     float64   `json:"overall"`
	CreatedAt   time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint      `json:"id"`
	ReviewID  uint      `json:"review_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Review) Validate() error {
	if r.Content == "" {
		return ErrInvalidReviewContent
	}
	if r.Rating < 1 || r.Rating > 5 {
		return ErrInvalidRating
	}
	return nil
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	GetByID(ctx context.Context, id uint) (*Review, error)
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id uint) error

	// By place
	GetByPlaceID(ctx context.Context, placeID uuid.UUID, limit, offset int) ([]*Review, int, error)

	// By user
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, int, error)

	// Check if user already reviewed place
	ExistsByUserAndPlace(ctx context.Context, userID, placeID uuid.UUID) (bool, error)

	// Comments
	CreateComment(ctx context.Context, comment *Comment) error
	GetCommentsByReviewID(ctx context.Context, reviewID uint, limit, offset int) ([]*Comment, int, error)

	// Ratings
	CreatePlaceRating(ctx context.Context, rating *PlaceRating) error
	GetPlaceRatings(ctx context.Context, placeID uuid.UUID) ([]*PlaceRating, error)
	UpdatePlaceRatingStats(ctx context.Context, placeID uuid.UUID) error
}

type CreateReviewInput struct {
	PlaceID     uuid.UUID
	Title       string
	Content     string
	Rating      float64
	Images      []string
	Cleanliness float64
	Service     float64
	Quality     float64
	Value       float64
}

type UpdateReviewInput struct {
	Title       *string
	Content     *string
	Rating      *float64
	Images      *[]string
	Cleanliness *float64
	Service     *float64
	Quality     *float64
	Value       *float64
}

type ListParams struct {
	Limit  int
	Offset int
}

type ReviewsResult struct {
	Reviews []*Review
	Total   int
	Limit   int
	Offset  int
}

type CommentsResult struct {
	Comments []*Comment
	Total    int
	Limit    int
	Offset   int
}

type ReviewService interface {
	Create(ctx context.Context, authorID uuid.UUID, input *CreateReviewInput) (*Review, error)
	GetByID(ctx context.Context, id uint) (*Review, error)
	Update(ctx context.Context, id uint, authorID uuid.UUID, input *UpdateReviewInput) (*Review, error)
	Delete(ctx context.Context, id uint, authorID uuid.UUID) error
	ListByPlace(ctx context.Context, placeID uuid.UUID, params *ListParams) (*ReviewsResult, error)
	ListByUser(ctx context.Context, userID uuid.UUID, params *ListParams) (*ReviewsResult, error)
	AddComment(ctx context.Context, reviewID uint, authorID uuid.UUID, content string) (*Comment, error)
	ListComments(ctx context.Context, reviewID uint, params *ListParams) (*CommentsResult, error)
}
