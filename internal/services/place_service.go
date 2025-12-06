package services

import (
	"context"

	"TulaHackDays-Backend/internal/domain"

	"github.com/google/uuid"
)

type placeService struct {
	placeRepo domain.PlaceRepository
	txManager domain.TxManager
}

func NewPlaceService(placeRepo domain.PlaceRepository, txManager domain.TxManager) domain.PlaceService {
	return &placeService{
		placeRepo: placeRepo,
		txManager: txManager,
	}
}

func (s *placeService) Create(ctx context.Context, userID uuid.UUID, input *domain.CreatePlaceInput) (*domain.Place, error) {
	place := &domain.Place{
		ID:            uuid.New(),
		Name:          input.Name,
		Description:   input.Description,
		Type:          input.Type,
		Address:       input.Address,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		Phone:         input.Phone,
		Website:       input.Website,
		Images:        input.Images,
		AverageRating: 0,
		ReviewCount:   0,
		IsActive:      true,
		CreatedByID:   userID,
	}

	if err := place.Validate(); err != nil {
		return nil, err
	}

	return place, s.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Place = s.placeRepo
		return repos.Place.Create(ctx, place)
	})
}

func (s *placeService) GetByID(ctx context.Context, placeID uuid.UUID) (*domain.Place, error) {
	return s.placeRepo.GetByID(ctx, placeID)
}

func (s *placeService) Update(ctx context.Context, placeID, userID uuid.UUID, input *domain.UpdatePlaceInput) (*domain.Place, error) {
	var updatedPlace *domain.Place

	err := s.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Place = s.placeRepo

		place, err := repos.Place.GetByID(ctx, placeID)
		if err != nil {
			return err
		}

		if place.CreatedByID != userID {
			return domain.ErrAccessDenied
		}

		if input.Name != nil {
			place.Name = *input.Name
		}
		if input.Description != nil {
			place.Description = *input.Description
		}
		if input.Address != nil {
			place.Address = *input.Address
		}
		if input.Latitude != nil {
			place.Latitude = *input.Latitude
		}
		if input.Longitude != nil {
			place.Longitude = *input.Longitude
		}
		if input.Phone != nil {
			place.Phone = *input.Phone
		}
		if input.Website != nil {
			place.Website = *input.Website
		}
		if input.Images != nil {
			place.Images = *input.Images
		}

		if err := place.Validate(); err != nil {
			return err
		}

		if err := repos.Place.Update(ctx, place); err != nil {
			return err
		}

		updatedPlace = place
		return nil
	})

	return updatedPlace, err
}

func (s *placeService) Delete(ctx context.Context, placeID, userID uuid.UUID) error {
	return s.txManager.WithinTx(ctx, func(ctx context.Context, repos *domain.Repos) error {
		repos.Place = s.placeRepo

		place, err := repos.Place.GetByID(ctx, placeID)
		if err != nil {
			return err
		}

		if place.CreatedByID != userID {
			return domain.ErrAccessDenied
		}

		return repos.Place.Delete(ctx, placeID)
	})
}

func (s *placeService) Search(ctx context.Context, params *domain.SearchPlacesParams) (*domain.PlacesResult, error) {
	var coords *domain.Coordinates
	if params.Latitude != nil && params.Longitude != nil {
		var err error
		coords, err = domain.NewCoordinates(*params.Latitude, *params.Longitude)
		if err != nil {
			return nil, err
		}
	}

	places, total, err := s.placeRepo.Search(ctx, params.Query, params.Type, coords, params.RadiusKm, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	return &domain.PlacesResult{
		Places: places,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}, nil
}

func (s *placeService) GetNearby(ctx context.Context, params *domain.NearbyPlacesParams) ([]*domain.Place, error) {
	coords, err := domain.NewCoordinates(params.Latitude, params.Longitude)
	if err != nil {
		return nil, err
	}

	return s.placeRepo.GetNearby(ctx, coords, params.RadiusKm, params.Limit)
}
