package dto

import (
	"time"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type PlaceDTO struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	Address       string    `json:"address"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Phone         string    `json:"phone"`
	Website       string    `json:"website"`
	Images        []string  `json:"images"`
	AverageRating float64   `json:"average_rating"`
	ReviewCount   int       `json:"review_count"`
	IsActive      bool      `json:"is_active"`
	CreatedByID   uuid.UUID `json:"created_by_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreatePlaceRequestDTO struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Type        string   `json:"type" binding:"required,oneof=healthy_food alcohol tobacco health"`
	Address     string   `json:"address" binding:"required"`
	Latitude    float64  `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude   float64  `json:"longitude" binding:"required,gte=-180,lte=180"`
	Phone       string   `json:"phone"`
	Website     string   `json:"website"`
	Images      []string `json:"images"`
}

type UpdatePlaceRequestDTO struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Address     *string   `json:"address"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Phone       *string   `json:"phone"`
	Website     *string   `json:"website"`
	Images      *[]string `json:"images"`
}

type SearchPlacesRequestDTO struct {
	Query  string  `form:"query"`
	Type   string  `form:"type"`
	Lat    float64 `form:"lat"`
	Lng    float64 `form:"lng"`
	Radius float64 `form:"radius" default:"5"`
	Limit  int     `form:"limit" default:"20"`
	Offset int     `form:"offset" default:"0"`
}

type PlacesResponseDTO struct {
	Places []PlaceDTO `json:"places"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// Convert entity to DTO
func PlaceEntityToDTO(place *domain.Place) PlaceDTO {
	return PlaceDTO{
		ID:            place.ID,
		Name:          place.Name,
		Description:   place.Description,
		Type:          string(place.Type),
		Address:       place.Address,
		Latitude:      place.Latitude,
		Longitude:     place.Longitude,
		Phone:         place.Phone,
		Website:       place.Website,
		Images:        place.Images,
		AverageRating: place.AverageRating,
		ReviewCount:   place.ReviewCount,
		IsActive:      place.IsActive,
		CreatedByID:   place.CreatedByID,
		CreatedAt:     place.CreatedAt,
	}
}

func PlacesEntitiesToDTOs(places []*domain.Place) []PlaceDTO {
	dtos := make([]PlaceDTO, len(places))
	for i, place := range places {
		dtos[i] = PlaceEntityToDTO(place)
	}
	return dtos
}
