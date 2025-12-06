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

func NewReviewService(reviewRepo domain.ReviewRepository, txManager domain.TxManager) domain.ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		txManager:  txManager,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, userID uuid.UUID, input *domain.CreateReviewInput) (*domain.Review, error) {
	var createdReview *domain.Review

	err := s.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Review = s.reviewRepo

		// Check if user already reviewed this place
		exists, err := repos.Review.ExistsByUserAndPlace(ctx, userID, input.PlaceID)
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrUserAlreadyReviewed
		}

		// Create review
		review := &domain.Review{
			PlaceID:  input.PlaceID,
			UserID:   userID,
			Title:    input.Title,
			Content:  input.Content,
			Images:   input.Images,
			Rating:   input.Rating,
			IsActive: true,
		}

		if err := review.Validate(); err != nil {
			return err
		}

		if err := repos.Review.Create(ctx, review); err != nil {
			return err
		}

		// Create detailed rating if provided
		if input.Cleanliness > 0 && input.Service > 0 && input.Quality > 0 && input.Value > 0 {
			overall := (input.Cleanliness + input.Service + input.Quality + input.Value) / 4
			rating := &domain.PlaceRating{
				PlaceID:     input.PlaceID,
				UserID:      userID,
				Cleanliness: input.Cleanliness,
				Service:     input.Service,
				Quality:     input.Quality,
				Value:       input.Value,
				Overall:     overall,
			}

			if err := repos.Review.CreatePlaceRating(ctx, rating); err != nil {
				return err
			}

			// Update place rating stats
			if err := repos.Review.UpdatePlaceRatingStats(ctx, input.PlaceID); err != nil {
				return err
			}
		}

		createdReview = review
		return nil
	})

	return createdReview, err
}

func (s *reviewService) GetReview(ctx context.Context, reviewID uint) (*domain.Review, error) {
	return s.reviewRepo.GetByID(ctx, reviewID)
}

func (s *reviewService) UpdateReview(ctx context.Context, reviewID uint, userID uuid.UUID, input *domain.UpdateReviewInput) (*domain.Review, error) {
	var updatedReview *domain.Review

	err := s.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Review = s.reviewRepo

		// Get existing review
		review, err := repos.Review.GetByID(ctx, reviewID)
		if err != nil {
			return err
		}

		// Check if user is the author
		if review.UserID != userID {
			return domain.ErrAccessDenied
		}

		// Apply updates
		if input.Title != nil {
			review.Title = *input.Title
		}
		if input.Content != nil {
			review.Content = *input.Content
		}
		if input.Rating != nil {
			review.Rating = *input.Rating
		}
		if input.Images != nil {
			review.Images = *input.Images
		}

		// Validate updated review
		if err := review.Validate(); err != nil {
			return err
		}

		// Update in repository
		if err := repos.Review.Update(ctx, review); err != nil {
			return err
		}

		// Update detailed rating if provided
		if input.Cleanliness != nil && input.Service != nil && input.Quality != nil && input.Value != nil {
			overall := (*input.Cleanliness + *input.Service + *input.Quality + *input.Value) / 4
			rating := &domain.PlaceRating{
				PlaceID:     review.PlaceID,
				UserID:      userID,
				Cleanliness: *input.Cleanliness,
				Service:     *input.Service,
				Quality:     *input.Quality,
				Value:       *input.Value,
				Overall:     overall,
			}

			if err := repos.Review.CreatePlaceRating(ctx, rating); err != nil {
				return err
			}

			// Update place rating stats
			if err := repos.Review.UpdatePlaceRatingStats(ctx, review.PlaceID); err != nil {
				return err
			}
		}

		updatedReview = review
		return nil
	})

	return updatedReview, err
}

func (s *reviewService) DeleteReview(ctx context.Context, reviewID uint, userID uuid.UUID) error {
	return s.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Review = s.reviewRepo

		// Get existing review to check ownership
		review, err := repos.Review.GetByID(ctx, reviewID)
		if err != nil {
			return err
		}

		// Check if user is the author
		if review.UserID != userID {
			return domain.ErrAccessDenied
		}

		return repos.Review.Delete(ctx, reviewID)
	})
}

func (s *reviewService) GetReviewsByPlace(ctx context.Context, placeID uuid.UUID, params *domain.ListParams) (*domain.ReviewsResult, error) {
	reviews, total, err := s.reviewRepo.GetByPlaceID(ctx, placeID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	return &domain.ReviewsResult{
		Reviews: reviews,
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
	}, nil
}

func (s *reviewService) GetReviewsByUser(ctx context.Context, userID uuid.UUID, params *domain.ListParams) (*domain.ReviewsResult, error) {
	reviews, total, err := s.reviewRepo.GetByUserID(ctx, userID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	return &domain.ReviewsResult{
		Reviews: reviews,
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
	}, nil
}

func (s *reviewService) AddComment(ctx context.Context, reviewID uint, userID uuid.UUID, content string) (*domain.Comment, error) {
	if content == "" {
		return nil, domain.ErrInvalidCommentContent
	}

	comment := &domain.Comment{
		ReviewID:  reviewID,
		UserID:    userID,
		Content:   content,
		IsActive:  true,
	}

	err := s.reviewRepo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *reviewService) GetCommentsByReview(ctx context.Context, reviewID uint, params *domain.ListParams) (*domain.CommentsResult, error) {
	comments, total, err := s.reviewRepo.GetCommentsByReviewID(ctx, reviewID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	return &domain.CommentsResult{
		Comments: comments,
		Total:    total,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}, nil
}
