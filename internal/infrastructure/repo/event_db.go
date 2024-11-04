package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/db"
)

type EventDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *EventDBRepo) GetLast(
	ctx context.Context,
	rootID string,
	action domain.Action,
) (domain.Event, error) {
	query := `SELECT * FROM public.events WHERE root_id=@root_id AND action=@action`
	args := pgx.NamedArgs{
		"root_id": rootID,
		"action":  action.String(),
	}

	var res domain.Event
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&res.ID,
		&res.RootID,
		&res.Action,
		&res.Payload,
		&res.CreatedAt,
	)

	return res, err
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

func NewEventDBRepo(db *db.DB, log *zap.SugaredLogger) *EventDBRepo {
	return &EventDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
