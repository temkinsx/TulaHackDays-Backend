package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Additional place errors
var (
	ErrInvalidPlaceName = errors.New("invalid place name")
	ErrInvalidAddress   = errors.New("invalid address")
)

type PlaceType string

const (
	PlaceTypeHealthyFood PlaceType = "healthy_food"
	PlaceTypeAlcohol     PlaceType = "alcohol"
	PlaceTypeTobacco     PlaceType = "tobacco"
	PlaceTypeHealth      PlaceType = "health"
)

type Place struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        PlaceType `json:"type"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Phone       string    `json:"phone"`
	Website     string    `json:"website"`
	Images      []string  `json:"images"`

	// Rating stats
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int     `json:"review_count"`

	// Lifecycle
	IsActive    bool      `json:"is_active"`
	CreatedByID uuid.UUID `json:"created_by_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Place) Validate() error {
	if p.Name == "" {
		return ErrInvalidPlaceName
	}
	if p.Address == "" {
		return ErrInvalidAddress
	}
	if p.Latitude < -90 || p.Latitude > 90 || p.Longitude < -180 || p.Longitude > 180 {
		return ErrInvalidCoordinates
	}
	if !isValidPlaceType(p.Type) {
		return ErrInvalidPlaceType
	}
	return nil
}

func isValidPlaceType(placeType PlaceType) bool {
	validTypes := []PlaceType{
		PlaceTypeHealthyFood,
		PlaceTypeAlcohol,
		PlaceTypeTobacco,
		PlaceTypeHealth,
	}

	for _, validType := range validTypes {
		if placeType == validType {
			return true
		}
	}
	return false
}

type PlaceRepository interface {
	Create(ctx context.Context, place *Place) error
	GetByID(ctx context.Context, id uuid.UUID) (*Place, error)
	Update(ctx context.Context, place *Place) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Search and filtering
	Search(ctx context.Context, query string, placeType *PlaceType, limit, offset int) ([]*Place, int, error)
	GetNearby(ctx context.Context, coords *Coordinates, radiusKm float64, limit int) ([]*Place, error)

	// Statistics
	GetByTypeStats(ctx context.Context) (map[PlaceType]int, error)
	GetTopRated(ctx context.Context, limit int) ([]*Place, error)
}

type CreatePlaceInput struct {
	Name        string
	Description string
	Type        PlaceType
	Address     string
	Latitude    float64
	Longitude   float64
	Phone       string
	Website     string
	Images      []string
}

type UpdatePlaceInput struct {
	Name        *string
	Description *string
	Address     *string
	Latitude    *float64
	Longitude   *float64
	Phone       *string
	Website     *string
	Images      *[]string
}

type SearchPlacesParams struct {
	Query     string
	Type      *PlaceType
	Latitude  *float64
	Longitude *float64
	RadiusKm  float64
	Limit     int
	Offset    int
}

type NearbyPlacesParams struct {
	Latitude  float64
	Longitude float64
	RadiusKm  float64
	Limit     int
}

type PlacesResult struct {
	Places []*Place
	Total  int
	Limit  int
	Offset int
}

type PlaceService interface {
	Create(ctx context.Context, authorID uuid.UUID, input *CreatePlaceInput) (*Place, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Place, error)
	Update(ctx context.Context, id uuid.UUID, authorID uuid.UUID, input *UpdatePlaceInput) (*Place, error)
	Delete(ctx context.Context, id uuid.UUID, authorID uuid.UUID) error
	Search(ctx context.Context, params *SearchPlacesParams) (*PlacesResult, error)
	GetNearby(ctx context.Context, params *NearbyPlacesParams) ([]*Place, error)
}
