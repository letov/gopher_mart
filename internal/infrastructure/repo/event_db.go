package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
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
	query := `INSERT INTO public.events (root_id, action, payload) VALUES (@root_id, @action, @payload)`
	args := pgx.NamedArgs{
		"root_id": e.RootID,
		"action":  e.Action.String(),
		"payload": e.Payload,
	}
	_, err := r.pool.Exec(ctx, query, args)
	return err
}

func (r *EventDBRepo) HasEvent(
	ctx context.Context,
	rootID string,
	action domain.Action,
	durationSec time.Duration,
) (bool, error) {
	var (
		query   string
		args    pgx.NamedArgs
		counter int
	)

	if durationSec.Minutes() == 0 {
		query = `SELECT count(*) FROM public.events WHERE root_id=@root_id AND action=@action`
		args = pgx.NamedArgs{
			"root_id": rootID,
			"action":  action.String(),
		}
	} else {
		query = `SELECT count(*) FROM public.events WHERE root_id=@root_id AND action=@action AND created_at < NOW() - INTERVAL '@duration seconds'`
		args = pgx.NamedArgs{
			"root_id":  rootID,
			"action":   action.String(),
			"duration": durationSec.Minutes(),
		}
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&counter)

	if err != nil {
		return false, err
	}

	return counter > 0, nil
}

func NewEventDBRepo(db *db.DB, log *zap.SugaredLogger) *EventDBRepo {
	return &EventDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
