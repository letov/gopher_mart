package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/dto/out"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/db"
)

var (
	ErrUnknownStatus = errors.New("unknown operation status")
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

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var res []domain.Operation
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

func (r *OperationDBRepo) GetBalance(ctx context.Context, userId int64) (out.Balance, error) {
	query := `select status, sum(sum) from public.operations where user_id = @user_id GROUP BY status`
	args := pgx.NamedArgs{
		"user_id": userId,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return out.Balance{}, err
	}

	defer rows.Close()
	var (
		added    int64
		deducted int64
	)

	for rows.Next() {
		var status domain.OperationStatus
		var sum int64
		err = rows.Scan(
			&status,
			&sum,
		)
		if err != nil {
			return out.Balance{}, err
		}
		switch status {
		case domain.AddedStatus:
			added = sum
		case domain.DeductedStatus:
			deducted = sum
		default:
			return out.Balance{}, ErrUnknownStatus
		}
	}

	return out.Balance{
		Current:   added - deducted,
		Withdrawn: deducted,
	}, nil
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
