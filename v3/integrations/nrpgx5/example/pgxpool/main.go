package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/facily-tech/go-agent/v3/integrations/nrpgx5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	cfg, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432")
	if err != nil {
		panic(err)
	}

	cfg.ConnConfig.Tracer = nrpgx5.NewTracer()
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("PostgreSQL App"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		panic(err)
	}
	//
	// N.B.: We do not recommend using app.WaitForConnection in production code.
	//
	app.WaitForConnection(5 * time.Second)
	txn := app.StartTransaction("postgresQuery")

	ctx := newrelic.NewContext(context.Background(), txn)
	row := db.QueryRow(ctx, "SELECT count(*) FROM pg_catalog.pg_tables")
	count := 0
	err = row.Scan(&count)
	if err != nil {
		log.Println(err)
	}

	txn.End()
	app.Shutdown(5 * time.Second)

	fmt.Println("number of entries in pg_catalog.pg_tables", count)
}
