//go:build integration
// +build integration

package repository_test

import (
	"context"
	"log"
	"os"
	"test-task/internal/database"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	var ctx = context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate("../../../migrations", dsn)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	db.Close()
	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}
