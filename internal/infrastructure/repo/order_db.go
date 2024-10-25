package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/infrastructure/db"
	"gopher_mart/internal/infrastructure/dto/response"
)

type OrderDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *OrderDBRepo) Save(ctx context.Context, su args.CalcAccrual) error {
	return nil
}

func (r *OrderDBRepo) UpdateOrder(ctx context.Context, oa response.OrderAccrual) error {
	return nil
}

func NewOrderDBRepo(db *db.DB, log *zap.SugaredLogger) *OrderDBRepo {
	return &OrderDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
