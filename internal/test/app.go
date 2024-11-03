package test

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gopher_mart/internal/infrastructure/db"
	"gopher_mart/internal/infrastructure/di"
	"gopher_mart/internal/infrastructure/httpclient"
	"testing"
)

func InjectApp() fx.Option {
	cs := di.GetConstructors()

	// TODO: dirty switch to test implementation, chandes needed
	cs[4] = NewHttpClient
	cs[5] = fx.Annotate(func(h *HttpClient) httpclient.Client {
		return h
	}, fx.As(new(httpclient.Client)))

	return fx.Provide(
		cs...,
	)
}

func initTest(t *testing.T, r interface{}) {
	t.Setenv("IS_TEST_ENV", "true")
	app := fxtest.New(t, InjectApp(), fx.Invoke(r))
	defer app.RequireStop()
	app.RequireStart()
}

func flushDB(ctx context.Context, db *db.DB) error {
	pool := db.GetPool()
	query := `SELECT table_name "table" FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version';`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return err
	}

	var queries []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return err
		}
		queries = append(queries, fmt.Sprintf("TRUNCATE %v CASCADE;", table))
	}

	tx, _ := pool.Begin(ctx)
	for _, query := range queries {
		_, err = tx.Exec(ctx, query)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
