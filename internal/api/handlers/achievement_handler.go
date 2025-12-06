package handlers

import (
	"TulaHackDays-Backend/internal/api/common"
	"TulaHackDays-Backend/internal/api/dto"
	"TulaHackDays-Backend/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AchievementHandler struct {
	achService domain.AchievementService
}

func NewAchievementHandler(achievementService domain.AchievementService) *AchievementHandler {
	return &AchievementHandler{achService: achievementService}
}

func (h *AchievementHandler) ListActive(c *gin.Context) {
	achievements, err := h.achService.ListActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AchievementsEntitiesToDTOs(achievements))
}
