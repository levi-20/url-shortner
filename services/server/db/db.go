package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDatabaseConnection() *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	connectionString := os.Getenv("CONN")

	if connectionString == "" {

		host := os.Getenv("POSTGRES_HOST")
		user := os.Getenv("POSTGRES_USER")
		port := os.Getenv("POSTGRES_PORT")
		pass := os.Getenv("POSTGRES_PASSWORD")
		database := os.Getenv("POSTGRES_DB")

		if user == "" || pass == "" || database == "" || host == "" {
			slog.Error("Missing necessary env vars! Could not find any of CONN, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_HOST")
			os.Exit(-1)
		}

		connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, database)
	}

	connection, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		slog.Error("Database connection failed")
	}

	slog.Info("🔗 Connected to database", "database", connection.Config().ConnConfig.Database)

	return connection
}
