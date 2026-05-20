package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	connection_string := os.Getenv("CONN")
	if connection_string == "" {
		// No connection string? Try reading other env vars
		postgres_host := os.Getenv("POSTGRES_HOST")
		postgres_user := os.Getenv("POSTGRES_USER")
		postgres_password := os.Getenv("POSTGRES_PASSWORD")
		postgres_db := os.Getenv("POSTGRES_DB")

		if postgres_user == "" || postgres_password == "" || postgres_db == "" || postgres_host == "" {
			slog.Error("Missing necessary env vars! Could not find any of CONN, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_HOST")
			os.Exit(-1)
		}

		connection_string = fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", postgres_user, postgres_password, postgres_host, postgres_db)
	}

	m, err := migrate.New("file://db/migrations", connection_string)
	if err != nil {
		slog.Error("Failed to read DB migrations", "error", err)
		os.Exit(-2)
	}

	slog.Info("Running migrations...")

	// err = m.Force(1)
	err = m.Up()
	// err = m.Drop()

	// ErrNoChange is not an error
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to run DB migrations", "error", err)
		os.Exit(-3)
	}

	slog.Info("📦 Migrations applied successfully!")
}
