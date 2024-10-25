package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/dto/result"
	"gopher_mart/internal/infrastructure/db"
)

type UserDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *UserDBRepo) Save(ctx context.Context, su args.SaveUser) error {
	return nil
}

func (r *UserDBRepo) Login(ctx context.Context, l args.Login) (result.Login, error) {
	return result.Login{}, nil
}

func NewUserDBRepo(db *db.DB, log *zap.SugaredLogger) *UserDBRepo {
	return &UserDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
