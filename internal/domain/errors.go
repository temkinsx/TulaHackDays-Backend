package domain

import "errors"

var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Place errors
	ErrPlaceNotFound      = errors.New("place not found")
	ErrInvalidPlaceType   = errors.New("invalid place type")
	ErrInvalidCoordinates = errors.New("invalid coordinates")

	// Review errors
	ErrReviewNotFound      = errors.New("review not found")
	ErrUserAlreadyReviewed = errors.New("user already reviewed this place")
	ErrInvalidRating       = errors.New("invalid rating")

	// Achievement errors
	ErrAchievementNotFound = errors.New("achievement not found")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")

	// Authorization errors
	ErrAccessDenied = errors.New("access denied")

	// General errors
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
)
