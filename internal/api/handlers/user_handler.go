package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"TulaHackDays-Backend/internal/api/common"
	"TulaHackDays-Backend/internal/api/dto"
	"TulaHackDays-Backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService        domain.UserService
	achievementService domain.AchievementService
}

func NewUserHandler(userSvc domain.UserService, achievementService domain.AchievementService) *UserHandler {
	return &UserHandler{userService: userSvc, achievementService: achievementService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	loginResult, err := h.userService.Register(c.Request.Context(), &domain.RegisterUserInput{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.LoginUserResponseDTO{
		Token: loginResult.Token,
		User: dto.UserDTO{
			ID:        loginResult.User.ID,
			Email:     loginResult.User.Email,
			Username:  loginResult.User.Username,
			FirstName: loginResult.User.FirstName,
			LastName:  loginResult.User.LastName,
			Avatar:    loginResult.User.Avatar,
			Points:    loginResult.User.Points,
			Level:     loginResult.User.Level,
			IsActive:  loginResult.User.IsActive,
			CreatedAt: loginResult.User.CreatedAt,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginUserRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := h.userService.Login(c.Request.Context(), &domain.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: err.Error()})
		return
	}

	userDTO := dto.UserDTO{
		ID:        response.User.ID,
		Email:     response.User.Email,
		Username:  response.User.Username,
		FirstName: response.User.FirstName,
		LastName:  response.User.LastName,
		Avatar:    response.User.Avatar,
		Points:    response.User.Points,
		Level:     response.User.Level,
		IsActive:  response.User.IsActive,
		CreatedAt: response.User.CreatedAt,
	}

	c.JSON(http.StatusOK, dto.LoginUserResponseDTO{
		Token: response.Token,
		User:  userDTO,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "user not authenticated"})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})
		return
	}

	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Points:    user.Points,
		Level:     user.Level,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, userDTO)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req dto.UpdateProfileRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &domain.UpdateProfileInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Avatar:    req.Avatar,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Points:    user.Points,
		Level:     user.Level,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, userDTO)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req dto.ChangePasswordRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), userID, &domain.ChangePasswordInput{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{Message: "Password changed successfully"})
}

func (h *UserHandler) GetLeaderboard(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	users, err := h.userService.GetLeaderboard(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = dto.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
			Points:    user.Points,
			Level:     user.Level,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, userDTOs)
}

func (h *UserHandler) GetUserAchievements(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "user not authenticated"})
		return
	}

	achievements, err := h.achievementService.GetUserAchievements(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, achievements)
}
