package domain

import "context"

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, repos *Repos) error) error
}

type Repos struct {
	User        UserRepository
	Place       PlaceRepository
	Review      ReviewRepository
	Achievement AchievementRepository
}
