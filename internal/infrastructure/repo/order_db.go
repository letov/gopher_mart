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
