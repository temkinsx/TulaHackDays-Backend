package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"TulaHackDays-Backend/internal/api/common"
	"TulaHackDays-Backend/internal/api/dto"
	"TulaHackDays-Backend/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewSvc domain.ReviewService
}

func NewReviewHandler(reviewSvc domain.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: reviewSvc}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var req dto.CreateReviewRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	placeID, err := uuid.Parse(req.PlaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid place ID"})
		return
	}

	input := &domain.CreateReviewInput{
		PlaceID:     placeID,
		Title:       req.Title,
		Content:     req.Content,
		Rating:      req.Rating,
		Images:      req.Images,
		Cleanliness: req.Cleanliness,
		Service:     req.Service,
		Quality:     req.Quality,
		Value:       req.Value,
	}

	review, err := h.reviewSvc.Create(c.Request.Context(), userID, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyReviewed) {
			c.JSON(http.StatusConflict, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ReviewEntityToDTO(review))
}

func (h *ReviewHandler) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid review ID"})
		return
	}

	review, err := h.reviewSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrReviewNotFound) {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ReviewEntityToDTO(review))
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid review ID"})
		return
	}

	var req dto.UpdateReviewRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	input := &domain.UpdateReviewInput{
		Title:       req.Title,
		Content:     req.Content,
		Rating:      req.Rating,
		Images:      req.Images,
		Cleanliness: req.Cleanliness,
		Service:     req.Service,
		Quality:     req.Quality,
		Value:       req.Value,
	}

	review, err := h.reviewSvc.Update(c.Request.Context(), uint(id), userID, input)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ReviewEntityToDTO(review))
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid review ID"})
		return
	}

	if err := h.reviewSvc.Delete(c.Request.Context(), uint(id), userID); err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{Message: "Review deleted successfully"})
}

func (h *ReviewHandler) GetReviewsByPlace(c *gin.Context) {
	placeIDStr := c.Param("place_id")
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid place ID"})
		return
	}

	limit, offset := parsePagination(c)

	params := &domain.ListParams{Limit: limit, Offset: offset}
	result, err := h.reviewSvc.ListByPlace(c.Request.Context(), placeID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	dtoReviews := dto.ReviewsEntitiesToDTOs(result.Reviews)
	c.JSON(http.StatusOK, dto.ReviewsResponseDTO{
		Reviews: dtoReviews,
		Total:   result.Total,
		Limit:   result.Limit,
		Offset:  result.Offset,
	})
}

func (h *ReviewHandler) GetReviewsByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	limit, offset := parsePagination(c)

	params := &domain.ListParams{Limit: limit, Offset: offset}
	result, err := h.reviewSvc.ListByUser(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	dtoReviews := dto.ReviewsEntitiesToDTOs(result.Reviews)
	c.JSON(http.StatusOK, dto.ReviewsResponseDTO{
		Reviews: dtoReviews,
		Total:   result.Total,
		Limit:   result.Limit,
		Offset:  result.Offset,
	})
}

func (h *ReviewHandler) AddComment(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid review ID"})
		return
	}

	var req dto.AddCommentRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	comment, err := h.reviewSvc.AddComment(c.Request.Context(), uint(reviewID), userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CommentEntityToDTO(comment))
}

func (h *ReviewHandler) GetCommentsByReview(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid review ID"})
		return
	}

	limit, offset := parsePagination(c)

	params := &domain.ListParams{Limit: limit, Offset: offset}
	result, err := h.reviewSvc.ListComments(c.Request.Context(), uint(reviewID), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	dtoComments := make([]dto.CommentDTO, len(result.Comments))
	for i, comment := range result.Comments {
		dtoComments[i] = dto.CommentEntityToDTO(comment)
	}

	c.JSON(http.StatusOK, dto.CommentsResponseDTO{
		Comments: dtoComments,
		Total:    result.Total,
		Limit:    result.Limit,
		Offset:   result.Offset,
	})
}

func parsePagination(c *gin.Context) (limit, offset int) {
	limit = 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset = 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	return
}
