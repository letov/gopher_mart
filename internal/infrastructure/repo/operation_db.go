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

type OperationDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *OperationDBRepo) GetByUserId(ctx context.Context, userId int64) ([]domain.Operation, error) {
	query := `SELECT * FROM public.operations WHERE user_id = @user_id`
	args := pgx.NamedArgs{
		"user_id": userId,
	}

	var res []domain.Operation
	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var r domain.Operation
		err = rows.Scan(
			&r.ID,
			&r.OrderID,
			&r.UserID,
			&r.Status,
			&r.Sum,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func (r *OperationDBRepo) Save(ctx context.Context, so in.SaveOperation) error {
	query := `INSERT INTO public.operations (order_id, user_id, status, sum) VALUES (@order_id, @user_id, @status, @sum)`
	args := pgx.NamedArgs{
		"order_id": so.OrderID,
		"user_id":  so.UserID,
		"status":   so.Status,
		"sum":      so.Sum,
	}
	_, err := r.pool.Exec(ctx, query, args)
	return err
}

func NewOperationDBRepo(db *db.DB, log *zap.SugaredLogger) *OperationDBRepo {
	return &OperationDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
