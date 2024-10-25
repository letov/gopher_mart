package db

import (
	"context"
	"embed"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopher_mart/internal/infrastructure/config"
	"time"
)

type DB struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

var (
	ErrNoDBConnection = errors.New("DB connection failed")
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) SetPool(pool *pgxpool.Pool) {
	db.pool = pool
}

func (db *DB) Ping(ctx context.Context) (err error) {
	pool := db.GetPool()
	if pool == nil {
		return ErrNoDBConnection
	}

	return pool.Ping(ctx)
}

func (db *DB) makeMigrations() {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	sqlDB := stdlib.OpenDBFromPool(db.GetPool())
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		panic(err)
	}
	if err := goose.Version(sqlDB, "migrations"); err != nil {
		db.log.Fatal(err)
	}
}

func NewDB(lc fx.Lifecycle, log *zap.SugaredLogger, c *config.Config) *DB {
	db := &DB{
		log: log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(c.DBConnection) == 0 {
		log.Info("Empty DB config")
		return db
	}

	log.Info("Init DB pool")

	poolConfig, err := pgxpool.ParseConfig(c.DBConnection)
	if err != nil {
		log.Warn("Failed to parce config: ", err)
		return db
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Warn("Failed to create db pool: ", err)
		return db
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Warn("Failed to ping db: ", err)
		return db
	}

	db.SetPool(pool)
	db.makeMigrations()

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if db.pool == nil {
				return nil
			}

			log.Info("Close DB pool")
			db.pool.Close()
			return nil
		},
	})

	return db
}
