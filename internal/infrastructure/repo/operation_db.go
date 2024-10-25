package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/infrastructure/db"
)

type OperationDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func NewOperationDBRepo(db *db.DB, log *zap.SugaredLogger) *OperationDBRepo {
	return &OperationDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
