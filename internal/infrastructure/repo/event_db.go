package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/db"
	"time"
)

type EventDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *EventDBRepo) Save(ctx context.Context, e domain.Event) error {
	return nil
}

func (r *EventDBRepo) HasEvent(ctx context.Context, rootID string, action domain.Action) bool {
	return false
}

func (r *EventDBRepo) HasEventWithDuration(ctx context.Context, rootID string, action domain.Action, duration time.Duration) bool {
	return false
}

func NewEventDBRepo(db *db.DB, log *zap.SugaredLogger) *EventDBRepo {
	return &EventDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
