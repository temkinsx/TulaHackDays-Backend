package dto

import (
	"time"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type ReviewDTO struct {
	ID        uint      `json:"id"`
	PlaceID   uuid.UUID `json:"place_id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Images    []string  `json:"images"`
	Rating    float64   `json:"rating"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateReviewRequestDTO struct {
	PlaceID     string   `json:"place_id" binding:"required"`
	Title       string   `json:"title"`
	Content     string   `json:"content" binding:"required"`
	Rating      float64  `json:"rating" binding:"required,gte=1,lte=5"`
	Images      []string `json:"images"`
	Cleanliness float64  `json:"cleanliness" binding:"required,gte=1,lte=5"`
	Service     float64  `json:"service" binding:"required,gte=1,lte=5"`
	Quality     float64  `json:"quality" binding:"required,gte=1,lte=5"`
	Value       float64  `json:"value" binding:"required,gte=1,lte=5"`
}

type UpdateReviewRequestDTO struct {
	Title       *string   `json:"title"`
	Content     *string   `json:"content"`
	Rating      *float64  `json:"rating"`
	Images      *[]string `json:"images"`
	Cleanliness *float64  `json:"cleanliness"`
	Service     *float64  `json:"service"`
	Quality     *float64  `json:"quality"`
	Value       *float64  `json:"value"`
}

type ReviewsResponseDTO struct {
	Reviews []ReviewDTO `json:"reviews"`
	Total   int         `json:"total"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
}

type CommentDTO struct {
	ID        uint      `json:"id"`
	ReviewID  uint      `json:"review_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentsResponseDTO struct {
	Comments []CommentDTO `json:"comments"`
	Total    int          `json:"total"`
	Limit    int          `json:"limit"`
	Offset   int          `json:"offset"`
}

type AddCommentRequestDTO struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// Convert entity to DTO
func ReviewEntityToDTO(review *domain.Review) ReviewDTO {
	return ReviewDTO{
		ID:        review.ID,
		PlaceID:   review.PlaceID,
		UserID:    review.UserID,
		Title:     review.Title,
		Content:   review.Content,
		Images:    review.Images,
		Rating:    review.Rating,
		IsActive:  review.IsActive,
		CreatedAt: review.CreatedAt,
	}
}

func ReviewsEntitiesToDTOs(reviews []*domain.Review) []ReviewDTO {
	dtos := make([]ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = ReviewEntityToDTO(review)
	}
	return dtos
}

func CommentEntityToDTO(comment *domain.Comment) CommentDTO {
	return CommentDTO{
		ID:        comment.ID,
		ReviewID:  comment.ReviewID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		IsActive:  comment.IsActive,
		CreatedAt: comment.CreatedAt,
	}
}

func CommentsEntitiesToDTOs(comments []*domain.Comment) []CommentDTO {
	dtos := make([]CommentDTO, len(comments))
	for i, comment := range comments {
		dtos[i] = CommentEntityToDTO(comment)
	}
	return dtos
}
