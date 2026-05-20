package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
		panic(err)
	}

	slog.Info("🔗 Connected to database", "database", connection.Config().ConnConfig.Database)

	return connection
}

type ExistingURLs struct {
	Urls   map[string]string
	Hashes map[string]struct{}
}

func CheckExistingURLs(ctx *gin.Context, pool *pgxpool.Pool, hashes []string) (*ExistingURLs, error) {

	sql := `
		SELECT 
			code,
			url_hash,
			url 
		FROM redirection 
		WHERE url_hash = ANY($1::varchar[])
	`
	rows, err := pool.Query(ctx, sql, hashes)
	if err != nil {
		return nil, err
	}

	var existingUrls ExistingURLs
	existingUrls.Urls = map[string]string{}
	existingUrls.Hashes = map[string]struct{}{}
	var hash, code, url string

	for rows.Next() {

		err := rows.Scan(&code, &hash, &url)
		if err != nil {
			slog.Error("error while scanning existing urls", "error", err.Error())
			return nil, err
		}

		existingUrls.Urls[code] = url
		existingUrls.Hashes[hash] = struct{}{}

	}

	return &existingUrls, nil
}

func SaveURLs(ctx *gin.Context, pool *pgxpool.Pool, original []string, short []string, hashes []string) error {

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	sql := `
		INSERT INTO redirection (code, url, url_hash, expire_at)
			SELECT
				UNNEST($1::varchar[]),  			-- code
				UNNEST($2::text[]),						-- url
				UNNEST($3::varchar[]),					-- url_hash
				NOW() + INTERVAL '2 days' 		-- expire_at
	`
	_, err = tx.Exec(ctx, sql, short, original, hashes)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	sql = `
		INSERT INTO metadata (code)
		SELECT 
			UNNEST($1::varchar[])					-- code
	`

	_, err = tx.Exec(ctx, sql, short)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
	tx.Conn().Close(ctx)

	return err
}

func GetUrlWithCode(ctx *gin.Context, pool *pgxpool.Pool, shortCode string) *string {

	sql := `
		SELECT
			url
		FROM redirection
		WHERE code = $1
	`
	var url string
	row := pool.QueryRow(ctx, sql, shortCode)
	err := row.Scan(&url)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("getUrlWithCode code not found", "code", shortCode, "error", err)
		} else {
			slog.Error("getUrlWithCode query failed", "code", shortCode, "error", err)
		}
		return nil
	}
	return &url
}

func RegisterClick(ctx *gin.Context, pool *pgxpool.Pool, code string) {

	sql := `
		UPDATE metadata
		SET clicks = clicks + 1
		WHERE code = $1
	`
	_, err := pool.Exec(ctx, sql, code)

	if err != nil {
		slog.Error("RegisterClick error in query", "error", err.Error())
	}
}
