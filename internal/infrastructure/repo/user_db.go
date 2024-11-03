package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/dto/out"
	"gopher_mart/internal/infrastructure/db"
)

type UserDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *UserDBRepo) Save(ctx context.Context, su in.SaveUser) error {
	query := `INSERT INTO public.users (login, password_hash) VALUES (@login, @password_hash)`
	args := pgx.NamedArgs{
		"login":         su.Login,
		"password_hash": su.PasswordHash,
	}
	_, err := r.pool.Exec(ctx, query, args)
	return err
}

func (r *UserDBRepo) Login(ctx context.Context, l in.Login) (out.Login, error) {
	var (
		query string
		args  pgx.NamedArgs
		id    int64
	)

	query = `SELECT id FROM public.users WHERE login=@login AND password_hash=@password_hash`
	args = pgx.NamedArgs{
		"login":         l.Login,
		"password_hash": l.PasswordHash,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&id)

	if err != nil {
		return out.Login{}, err
	}

	return out.Login{UserID: id}, nil
}

func NewUserDBRepo(db *db.DB, log *zap.SugaredLogger) *UserDBRepo {
	return &UserDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
