package services

import (
	"context"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type reviewService struct {
	reviewRepo domain.ReviewRepository
	txManager  domain.TxManager
}

func NewReviewService(reviewRepo domain.ReviewRepository, txManager domain.TxManager) domain.PlaceService {
	return &reviewService{
		reviewRepo: reviewRepo,
		txManager:  txManager,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, userID uuid.UUID, input *domain.CreateReviewInput) (*domain.Review, error) {
	panic("not implemented")
}

func (s *reviewService) GetReview(ctx context.Context, reviewID uint) (*domain.Review, error) {
	panic("not implemented")
}

func (s *reviewService) UpdateReview(ctx context.Context, reviewID uint, userID uuid.UUID, input *domain.UpdateReviewInput) (*domain.Review, error) {
	panic("not implemented")
}

func (s *reviewService) DeleteReview(ctx context.Context, reviewID uint, userID uuid.UUID) error {
	panic("not implemented")
}

func (s *reviewService) GetReviewsByPlace(ctx context.Context, placeID uuid.UUID, params *domain.ListParams) (*domain.ReviewsResult, error) {
	panic("not implemented")
}

func (s *reviewService) GetReviewsByUser(ctx context.Context, userID uuid.UUID, params *domain.ListParams) (*domain.ReviewsResult, error) {
	panic("not implemented")
}

func (s *reviewService) AddComment(ctx context.Context, reviewID uint, userID uuid.UUID, content string) (*domain.Comment, error) {
	panic("not implemented")
}

func (s *reviewService) GetCommentsByReview(ctx context.Context, reviewID uint, params *domain.ListParams) (*domain.CommentsResult, error) {
	panic("not implemented")
}
