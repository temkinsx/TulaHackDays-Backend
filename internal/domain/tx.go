package domain

import "context"

// TxManager allows executing application logic within a database transaction.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, repos *Repos) error) error
}

// Repos groups repository implementations that operate within a transaction.
type Repos struct {
	User        UserRepository
	Place       PlaceRepository
	Review      ReviewRepository
	Achievement AchievementRepository
}
