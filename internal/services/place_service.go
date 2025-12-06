package services

import (
	"context"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type placeService struct {
	placeRepo domain.PlaceRepository
	userSvc   domain.UserService
}

func NewPlaceService(placeRepo domain.PlaceRepository, userSvc domain.UserService) domain.PlaceService {
	return &placeService{
		placeRepo: placeRepo,
		userSvc:   userSvc,
	}
}

func (s *placeService) CreatePlace(ctx context.Context, userID uuid.UUID, input *domain.CreatePlaceInput) (*domain.Place, error) {
	panic("not implemented")
}

func (s *placeService) GetPlace(ctx context.Context, placeID uuid.UUID) (*domain.Place, error) {
	panic("not implemented")
}

func (s *placeService) UpdatePlace(ctx context.Context, placeID, userID uuid.UUID, input *domain.UpdatePlaceInput) (*domain.Place, error) {
	panic("not implemented")
}

func (s *placeService) DeletePlace(ctx context.Context, placeID, userID uuid.UUID) error {
	panic("not implemented")
}

func (s *placeService) SearchPlaces(ctx context.Context, params *domain.SearchPlacesParams) (*domain.PlacesResult, error) {
	panic("not implemented")
}

func (s *placeService) GetNearbyPlaces(ctx context.Context, params *domain.NearbyPlacesParams) ([]*domain.Place, error) {
	panic("not implemented")
}
