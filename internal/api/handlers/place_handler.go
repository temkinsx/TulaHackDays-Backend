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

type PlaceHandler struct {
	placeSvc domain.PlaceService
}

func NewPlaceHandler(placeSvc domain.PlaceService) *PlaceHandler {
	return &PlaceHandler{placeSvc: placeSvc}
}

func (h *PlaceHandler) CreatePlace(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var req dto.CreatePlaceRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	input := &domain.CreatePlaceInput{
		Name:        req.Name,
		Description: req.Description,
		Type:        domain.PlaceType(req.Type),
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Phone:       req.Phone,
		Website:     req.Website,
		Images:      req.Images,
	}

	place, err := h.placeSvc.Create(c.Request.Context(), userID, input)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	placeDTO := dto.PlaceEntityToDTO(place)
	c.JSON(http.StatusCreated, placeDTO)
}

func (h *PlaceHandler) GetPlace(c *gin.Context) {
	idStr := c.Param("id")
	placeID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid place ID"})
		return
	}

	place, err := h.placeSvc.GetByID(c.Request.Context(), placeID)
	if err != nil {
		if errors.Is(err, domain.ErrPlaceNotFound) {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	placeDTO := dto.PlaceEntityToDTO(place)
	c.JSON(http.StatusOK, placeDTO)
}

func (h *PlaceHandler) UpdatePlace(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	placeID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid place ID"})
		return
	}

	var req dto.UpdatePlaceRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	input := &domain.UpdatePlaceInput{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Phone:       req.Phone,
		Website:     req.Website,
		Images:      req.Images,
	}

	place, err := h.placeSvc.Update(c.Request.Context(), placeID, userID, input)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	placeDTO := dto.PlaceEntityToDTO(place)
	c.JSON(http.StatusOK, placeDTO)
}

func (h *PlaceHandler) DeletePlace(c *gin.Context) {
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	placeID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid place ID"})
		return
	}

	err = h.placeSvc.Delete(c.Request.Context(), placeID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, common.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{Message: "Place deleted successfully"})
}

func (h *PlaceHandler) SearchPlaces(c *gin.Context) {
	var searchReq dto.SearchPlacesRequestDTO

	searchReq.Query = c.Query("query")
	searchReq.Type = c.Query("type")

	if latStr := c.Query("lat"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			searchReq.Lat = lat
		}
	}

	if lngStr := c.Query("lng"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			searchReq.Lng = lng
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			searchReq.Radius = radius
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			searchReq.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			searchReq.Offset = offset
		}
	}

	var placeType *domain.PlaceType
	if searchReq.Type != "" {
		pt := domain.PlaceType(searchReq.Type)
		placeType = &pt
	}
	params := &domain.SearchPlacesParams{
		Query:     searchReq.Query,
		Type:      placeType,
		Latitude:  &searchReq.Lat,
		Longitude: &searchReq.Lng,
		RadiusKm:  searchReq.Radius,
		Limit:     searchReq.Limit,
		Offset:    searchReq.Offset,
	}

	response, err := h.placeSvc.Search(c.Request.Context(), params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	placeDTOs := dto.PlacesEntitiesToDTOs(response.Places)

	c.JSON(http.StatusOK, dto.PlacesResponseDTO{
		Places: placeDTOs,
		Total:  response.Total,
		Limit:  response.Limit,
		Offset: response.Offset,
	})
}

func (h *PlaceHandler) GetNearbyPlaces(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Latitude and longitude are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid latitude"})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid longitude"})
		return
	}

	radius := 5.0
	if radiusStr := c.Query("radius"); radiusStr != "" {
		if r, err := strconv.ParseFloat(radiusStr, 64); err == nil && r > 0 {
			radius = r
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	params := &domain.NearbyPlacesParams{
		Latitude:  lat,
		Longitude: lng,
		RadiusKm:  radius,
		Limit:     limit,
	}

	places, err := h.placeSvc.GetNearby(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}

	placeDTOs := dto.PlacesEntitiesToDTOs(places)
	c.JSON(http.StatusOK, placeDTOs)
}
