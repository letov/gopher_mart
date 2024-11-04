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

type OrderDBRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func (r *OrderDBRepo) Get(ctx context.Context, orderId string) (domain.Order, error) {
	query := `SELECT * FROM public.orders WHERE order_id = @order_id`
	args := pgx.NamedArgs{
		"order_id": orderId,
	}

	var res domain.Order
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&res.ID,
		&res.OrderID,
		&res.UserID,
		&res.Status,
		&res.Accrual,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	return res, err
}

func (r *OrderDBRepo) GetByUserId(ctx context.Context, userId int64) ([]domain.Order, error) {
	query := `SELECT * FROM public.orders WHERE user_id = @user_id ORDER BY created_at DESC`
	args := pgx.NamedArgs{
		"user_id": userId,
	}

	var res []domain.Order
	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var r domain.Order
		err = rows.Scan(
			&r.ID,
			&r.OrderID,
			&r.UserID,
			&r.Status,
			&r.Accrual,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func (r *OrderDBRepo) Save(ctx context.Context, ra in.RequestAccrual) error {
	query := `INSERT INTO public.orders (order_id, user_id) VALUES (@order_id, @user_id)`
	args := pgx.NamedArgs{
		"order_id": ra.OrderID,
		"user_id":  ra.UserID,
	}
	_, err := r.pool.Exec(ctx, query, args)
	return err
}

func (r *OrderDBRepo) UpdateOrder(ctx context.Context, uo in.UpdateOrder) error {
	query := `UPDATE public.orders SET status = @status, accrual = @accrual WHERE order_id=@order_id`
	args := pgx.NamedArgs{
		"status":   uo.Status,
		"accrual":  uo.Accrual,
		"order_id": uo.OrderID,
	}
	_, err := r.pool.Exec(ctx, query, args)
	return err
}

func NewOrderDBRepo(db *db.DB, log *zap.SugaredLogger) *OrderDBRepo {
	return &OrderDBRepo{
		pool: db.GetPool(),
		log:  log,
	}
}
