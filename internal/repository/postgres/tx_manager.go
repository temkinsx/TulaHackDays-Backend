package postgres

import (
	"context"

	"TulaHackDays-Backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type txManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) domain.TxManager {
	return &txManager{pool: pool}
}

func (m *txManager) WithinTx(ctx context.Context, fn func(ctx context.Context, repos *domain.Repos) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	repos := &domain.Repos{}

	if err := fn(ctx, repos); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
