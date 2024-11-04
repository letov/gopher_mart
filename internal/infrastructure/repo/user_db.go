package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/db"
)

type UserDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *UserDBRepo) Get(ctx context.Context, login string) (domain.User, error) {
	query := `SELECT * FROM public.users WHERE login=@login`
	args := pgx.NamedArgs{
		"login": login,
	}

	var res domain.User
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&res.ID,
		&res.Login,
		&res.PasswordHash,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	return res, err
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

func NewUserDBRepo(db *db.DB, log *zap.SugaredLogger) *UserDBRepo {
	return &UserDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
